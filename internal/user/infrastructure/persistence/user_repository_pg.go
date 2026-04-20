package userpersistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"macabi-back/internal/shared/pagination"
	userdomain "macabi-back/internal/user/domain"

	"gorm.io/gorm"
)

type UserModel struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"not null;default:'user'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserModel) TableName() string {
	return "users"
}

type UserRepositoryPG struct {
	db *gorm.DB
}

func NewUserRepositoryPG(db *gorm.DB) *UserRepositoryPG {
	return &UserRepositoryPG{db: db}
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&UserModel{})
}

func (r *UserRepositoryPG) Save(ctx context.Context, user *userdomain.User) error {
	model := toModel(user)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	user.ID = model.ID
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *UserRepositoryPG) FindByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	var model UserModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, userdomain.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return toDomain(&model), nil
}

func (r *UserRepositoryPG) FindByID(ctx context.Context, id string) (*userdomain.User, error) {
	var model UserModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, userdomain.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return toDomain(&model), nil
}

func (r *UserRepositoryPG) Update(ctx context.Context, user *userdomain.User) error {
	err := r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"name":     user.Name,
		"email":    user.Email,
		"role":     string(user.Role),
		"password": user.Password,
	}).Error
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *UserRepositoryPG) FindAll(ctx context.Context, params pagination.Params) ([]userdomain.User, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	var models []UserModel
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Find(&models).Error
	if err != nil {
		return nil, 0, fmt.Errorf("find all users: %w", err)
	}

	users := make([]userdomain.User, len(models))
	for i := range models {
		users[i] = *toDomain(&models[i])
	}
	return users, total, nil
}

func toModel(u *userdomain.User) *UserModel {
	return &UserModel{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Role:     string(u.Role),
	}
}

func toDomain(m *UserModel) *userdomain.User {
	return &userdomain.User{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		Role:      userdomain.Role(m.Role),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
