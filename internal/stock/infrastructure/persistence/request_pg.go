package stockpersistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"macabi-back/internal/shared/pagination"
	stockdomain "macabi-back/internal/stock/domain"
)

type RequestModel struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID      string    `gorm:"type:uuid;not null"`
	ResourceID     string    `gorm:"type:uuid;not null"`
	RequestedByID  string    `gorm:"type:uuid;not null"`
	Quantity       int       `gorm:"not null"`
	WithdrawalDate time.Time `gorm:"not null"`
	ReturnDate     *time.Time
	Status         string `gorm:"not null;default:'PENDIENTE'"`
	Notes          string
	CreatedAt      time.Time
}

func (RequestModel) TableName() string { return "stock_requests" }

type requestDetailRow struct {
	ID             string     `gorm:"column:id"`
	ProjectID      string     `gorm:"column:project_id"`
	ResourceID     string     `gorm:"column:resource_id"`
	RequestedByID  string     `gorm:"column:requested_by_id"`
	Quantity       int        `gorm:"column:quantity"`
	WithdrawalDate time.Time  `gorm:"column:withdrawal_date"`
	ReturnDate     *time.Time `gorm:"column:return_date"`
	Status         string     `gorm:"column:status"`
	Notes          string     `gorm:"column:notes"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	ResourceName   string     `gorm:"column:resource_name"`
	ResourceType   string     `gorm:"column:resource_type"`
	ProjectName    string     `gorm:"column:project_name"`
	RequesterName  string     `gorm:"column:requester_name"`
}

func (r *StockRepositoryPG) SaveRequest(ctx context.Context, req *stockdomain.ResourceRequest) error {
	m := RequestModel{
		ProjectID:      req.ProjectID,
		ResourceID:     req.ResourceID,
		RequestedByID:  req.RequestedByID,
		Quantity:       req.Quantity,
		WithdrawalDate: req.WithdrawalDate,
		ReturnDate:     req.ReturnDate,
		Status:         string(req.Status),
		Notes:          req.Notes,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	req.ID = m.ID
	req.CreatedAt = m.CreatedAt
	return nil
}

func (r *StockRepositoryPG) FindRequestByID(ctx context.Context, id string) (*stockdomain.ResourceRequest, error) {
	var m RequestModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, stockdomain.ErrRequestNotFound
		}
		return nil, err
	}
	return toDomainRequest(m), nil
}

func (r *StockRepositoryPG) FindRequestDetailByID(ctx context.Context, id string) (*stockdomain.RequestDetail, error) {
	var row requestDetailRow
	err := r.db.WithContext(ctx).
		Table("stock_requests sr").
		Select("sr.*, res.name as resource_name, res.type as resource_type, p.name as project_name, u.name as requester_name").
		Joins("JOIN stock_resources res ON res.id = sr.resource_id").
		Joins("JOIN projects p ON p.id = sr.project_id").
		Joins("JOIN users u ON u.id = sr.requested_by_id").
		Where("sr.id = ?", id).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == "" {
		return nil, stockdomain.ErrRequestNotFound
	}
	return toRequestDetail(row), nil
}

func (r *StockRepositoryPG) ApproveRequest(ctx context.Context, requestID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var req RequestModel
		if err := tx.First(&req, "id = ?", requestID).Error; err != nil {
			return err
		}
		if err := tx.Model(&RequestModel{}).Where("id = ?", requestID).
			Update("status", string(stockdomain.RequestStatusApproved)).Error; err != nil {
			return err
		}
		res := tx.Model(&ResourceModel{}).
			Where("id = ? AND available_stock >= ?", req.ResourceID, req.Quantity).
			UpdateColumn("available_stock", gorm.Expr("available_stock - ?", req.Quantity))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return stockdomain.ErrInsufficientStock
		}
		return nil
	})
}

func (r *StockRepositoryPG) ReturnRequest(ctx context.Context, requestID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var req RequestModel
		if err := tx.First(&req, "id = ?", requestID).Error; err != nil {
			return err
		}
		if err := tx.Model(&RequestModel{}).Where("id = ?", requestID).
			Update("status", string(stockdomain.RequestStatusReturned)).Error; err != nil {
			return err
		}
		return tx.Model(&ResourceModel{}).
			Where("id = ?", req.ResourceID).
			UpdateColumn("available_stock", gorm.Expr("available_stock + ?", req.Quantity)).Error
	})
}

func (r *StockRepositoryPG) UpdateRequestStatus(ctx context.Context, id string, status stockdomain.RequestStatus) error {
	return r.db.WithContext(ctx).Model(&RequestModel{}).
		Where("id = ?", id).
		Update("status", string(status)).Error
}

func (r *StockRepositoryPG) ListRequests(ctx context.Context, params pagination.Params, projectID string, onlyRequestedBy *string) (pagination.Result[stockdomain.RequestDetail], error) {
	q := r.db.WithContext(ctx).
		Table("stock_requests sr").
		Select("sr.*, res.name as resource_name, res.type as resource_type, p.name as project_name, u.name as requester_name").
		Joins("JOIN stock_resources res ON res.id = sr.resource_id").
		Joins("JOIN projects p ON p.id = sr.project_id").
		Joins("JOIN users u ON u.id = sr.requested_by_id")
	if projectID != "" {
		q = q.Where("sr.project_id = ?", projectID)
	}
	if onlyRequestedBy != nil && *onlyRequestedBy != "" {
		q = q.Where("sr.requested_by_id = ?", *onlyRequestedBy)
	}
	var total int64
	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return pagination.Result[stockdomain.RequestDetail]{}, err
	}
	var rows []requestDetailRow
	if err := q.Order("sr.created_at DESC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Scan(&rows).Error; err != nil {
		return pagination.Result[stockdomain.RequestDetail]{}, err
	}
	items := make([]stockdomain.RequestDetail, len(rows))
	for i, row := range rows {
		items[i] = *toRequestDetail(row)
	}
	return pagination.NewResult(items, total, params), nil
}

func (r *StockRepositoryPG) ListMyRequests(ctx context.Context, userID string, params pagination.Params) (pagination.Result[stockdomain.RequestDetail], error) {
	q := r.db.WithContext(ctx).
		Table("stock_requests sr").
		Select("sr.*, res.name as resource_name, res.type as resource_type, p.name as project_name, u.name as requester_name").
		Joins("JOIN stock_resources res ON res.id = sr.resource_id").
		Joins("JOIN projects p ON p.id = sr.project_id").
		Joins("JOIN users u ON u.id = sr.requested_by_id").
		Where("sr.requested_by_id = ?", userID)
	var total int64
	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return pagination.Result[stockdomain.RequestDetail]{}, err
	}
	var rows []requestDetailRow
	if err := q.Order("sr.created_at DESC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Scan(&rows).Error; err != nil {
		return pagination.Result[stockdomain.RequestDetail]{}, err
	}
	items := make([]stockdomain.RequestDetail, len(rows))
	for i, row := range rows {
		items[i] = *toRequestDetail(row)
	}
	return pagination.NewResult(items, total, params), nil
}

func toDomainRequest(m RequestModel) *stockdomain.ResourceRequest {
	return &stockdomain.ResourceRequest{
		ID:             m.ID,
		ProjectID:      m.ProjectID,
		ResourceID:     m.ResourceID,
		RequestedByID:  m.RequestedByID,
		Quantity:       m.Quantity,
		WithdrawalDate: m.WithdrawalDate,
		ReturnDate:     m.ReturnDate,
		Status:         stockdomain.RequestStatus(m.Status),
		Notes:          m.Notes,
		CreatedAt:      m.CreatedAt,
	}
}

func toRequestDetail(row requestDetailRow) *stockdomain.RequestDetail {
	return &stockdomain.RequestDetail{
		Request: stockdomain.ResourceRequest{
			ID:             row.ID,
			ProjectID:      row.ProjectID,
			ResourceID:     row.ResourceID,
			RequestedByID:  row.RequestedByID,
			Quantity:       row.Quantity,
			WithdrawalDate: row.WithdrawalDate,
			ReturnDate:     row.ReturnDate,
			Status:         stockdomain.RequestStatus(row.Status),
			Notes:          row.Notes,
			CreatedAt:      row.CreatedAt,
		},
		ResourceName:  row.ResourceName,
		ResourceType:  stockdomain.ResourceType(row.ResourceType),
		ProjectName:   row.ProjectName,
		RequesterName: row.RequesterName,
	}
}
