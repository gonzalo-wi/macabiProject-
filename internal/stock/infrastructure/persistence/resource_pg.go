package stockpersistence

import (
	"context"
	"errors"
	"strings"
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

func applyResourceListFilter(q *gorm.DB, filter stockports.ResourceListFilter) *gorm.DB {
	if s := strings.TrimSpace(filter.Query); s != "" {
		like := "%" + strings.ToLower(s) + "%"
		q = q.Where("LOWER(name) LIKE ?", like)
	}
	return q
}

func (r *StockRepositoryPG) FindAllResources(ctx context.Context, filter stockports.ResourceListFilter, params pagination.Params) (pagination.Result[stockdomain.Resource], error) {
	base := applyResourceListFilter(r.db.WithContext(ctx).Model(&ResourceModel{}), filter)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return pagination.Result[stockdomain.Resource]{}, err
	}
	var models []ResourceModel
	if err := applyResourceListFilter(r.db.WithContext(ctx).Model(&ResourceModel{}), filter).
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
