package stockpersistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"macabi-back/internal/shared/pagination"
	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
)

type ResourceModel struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name           string `gorm:"not null"`
	Type           string `gorm:"not null"`
	TotalStock     int    `gorm:"not null;default:0"`
	AvailableStock int    `gorm:"not null;default:0"`
	CreatedAt      time.Time
}

func (ResourceModel) TableName() string { return "stock_resources" }

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

type NotificationModel struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string `gorm:"type:uuid;not null"`
	RequestID string `gorm:"type:uuid;not null"`
	Message   string `gorm:"not null"`
	ReadAt    *time.Time
	CreatedAt time.Time
}

func (NotificationModel) TableName() string { return "stock_notifications" }

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

type StockRepositoryPG struct {
	db *gorm.DB
}

func NewStockRepositoryPG(db *gorm.DB) *StockRepositoryPG {
	return &StockRepositoryPG{db: db}
}

var _ stockports.StockRepository = (*StockRepositoryPG)(nil)
var _ stockports.ProjectMemberReader = (*StockRepositoryPG)(nil)
var _ stockports.UserEmailReader = (*StockRepositoryPG)(nil)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&ResourceModel{},
		&RequestModel{},
		&NotificationModel{},
		&PushSubscriptionModel{},
	)
}

func (r *StockRepositoryPG) SaveResource(ctx context.Context, res *stockdomain.Resource) error {
	m := ResourceModel{
		Name:           res.Name,
		Type:           string(res.Type),
		TotalStock:     res.TotalStock,
		AvailableStock: res.AvailableStock,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	res.ID = m.ID
	res.CreatedAt = m.CreatedAt
	return nil
}

func (r *StockRepositoryPG) FindResourceByID(ctx context.Context, id string) (*stockdomain.Resource, error) {
	var m ResourceModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, stockdomain.ErrResourceNotFound
		}
		return nil, err
	}
	return toDomainResource(m), nil
}

func (r *StockRepositoryPG) FindAllResources(ctx context.Context, params pagination.Params) (pagination.Result[stockdomain.Resource], error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&ResourceModel{}).Count(&total).Error; err != nil {
		return pagination.Result[stockdomain.Resource]{}, err
	}

	var models []ResourceModel
	if err := r.db.WithContext(ctx).
		Order("name ASC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Find(&models).Error; err != nil {
		return pagination.Result[stockdomain.Resource]{}, err
	}

	items := make([]stockdomain.Resource, len(models))
	for i, m := range models {
		items[i] = *toDomainResource(m)
	}
	return pagination.NewResult(items, total, params), nil
}

func (r *StockRepositoryPG) UpdateResource(ctx context.Context, res *stockdomain.Resource) error {
	return r.db.WithContext(ctx).Model(&ResourceModel{}).
		Where("id = ?", res.ID).
		Updates(map[string]interface{}{
			"name":            res.Name,
			"type":            string(res.Type),
			"total_stock":     res.TotalStock,
			"available_stock": res.AvailableStock,
		}).Error
}

func (r *StockRepositoryPG) DeleteResource(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&ResourceModel{}, "id = ?", id).Error
}

func (r *StockRepositoryPG) HasActiveRequests(ctx context.Context, resourceID string) (bool, error) {
	var count int64
	activeStatuses := []string{
		string(stockdomain.RequestStatusPending),
		string(stockdomain.RequestStatusApproved),
		string(stockdomain.RequestStatusDelivered),
	}
	err := r.db.WithContext(ctx).Model(&RequestModel{}).
		Where("resource_id = ? AND status IN ?", resourceID, activeStatuses).
		Count(&count).Error
	return count > 0, err
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

func (r *StockRepositoryPG) ListRequests(ctx context.Context, params pagination.Params, projectID string) (pagination.Result[stockdomain.RequestDetail], error) {
	q := r.db.WithContext(ctx).
		Table("stock_requests sr").
		Select("sr.*, res.name as resource_name, res.type as resource_type, p.name as project_name, u.name as requester_name").
		Joins("JOIN stock_resources res ON res.id = sr.resource_id").
		Joins("JOIN projects p ON p.id = sr.project_id").
		Joins("JOIN users u ON u.id = sr.requested_by_id")

	if projectID != "" {
		q = q.Where("sr.project_id = ?", projectID)
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

func (r *StockRepositoryPG) SaveNotification(ctx context.Context, n *stockdomain.StockNotification) error {
	m := NotificationModel{
		UserID:    n.UserID,
		RequestID: n.RequestID,
		Message:   n.Message,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	n.ID = m.ID
	n.CreatedAt = m.CreatedAt
	return nil
}

func (r *StockRepositoryPG) ListNotificationsByUser(ctx context.Context, userID string) ([]stockdomain.StockNotification, error) {
	var models []NotificationModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]stockdomain.StockNotification, len(models))
	for i, m := range models {
		out[i] = toDomainNotification(m)
	}
	return out, nil
}

func (r *StockRepositoryPG) MarkNotificationRead(ctx context.Context, id, userID string) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&NotificationModel{}).
		Where("id = ? AND user_id = ? AND read_at IS NULL", id, userID).
		Update("read_at", now)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return stockdomain.ErrNotificationNotFound
	}
	return nil
}

func (r *StockRepositoryPG) UnreadCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&NotificationModel{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

func (r *StockRepositoryPG) FindProjectCoordinators(ctx context.Context, projectID string) ([]string, error) {
	var userIDs []string
	err := r.db.WithContext(ctx).
		Table("project_members").
		Where("project_id = ? AND role = ?", projectID, "coordinator").
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}

func (r *StockRepositoryPG) IsProjectCoordinator(ctx context.Context, projectID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("project_members").
		Where("project_id = ? AND user_id = ? AND role = ?", projectID, userID, "coordinator").
		Count(&count).Error
	return count > 0, err
}

func (r *StockRepositoryPG) FindEmailByID(ctx context.Context, userID string) (string, error) {
	var email string
	err := r.db.WithContext(ctx).
		Table("users").
		Select("email").
		Where("id = ?", userID).
		Scan(&email).Error
	return email, err
}

func (r *StockRepositoryPG) FindEmailsByIDs(ctx context.Context, userIDs []string) (map[string]string, error) {
	if len(userIDs) == 0 {
		return map[string]string{}, nil
	}
	var rows []struct {
		ID    string
		Email string
	}
	err := r.db.WithContext(ctx).
		Table("users").
		Select("id, email").
		Where("id IN ?", userIDs).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, len(rows))
	for _, row := range rows {
		result[row.ID] = row.Email
	}
	return result, nil
}

func toDomainResource(m ResourceModel) *stockdomain.Resource {
	return &stockdomain.Resource{
		ID:             m.ID,
		Name:           m.Name,
		Type:           stockdomain.ResourceType(m.Type),
		TotalStock:     m.TotalStock,
		AvailableStock: m.AvailableStock,
		CreatedAt:      m.CreatedAt,
	}
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

func toDomainNotification(m NotificationModel) stockdomain.StockNotification {
	return stockdomain.StockNotification{
		ID:        m.ID,
		UserID:    m.UserID,
		RequestID: m.RequestID,
		Message:   m.Message,
		ReadAt:    m.ReadAt,
		CreatedAt: m.CreatedAt,
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
