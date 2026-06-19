package userpersistence

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"macabi-back/internal/shared/database"
	"macabi-back/internal/shared/pagination"
	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"

	"gorm.io/gorm"
)

type UserModel struct {
	ID       string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name     string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null;default:'user'"`
	// No `default` tag on the bool fields below: GORM omits zero-value fields
	// (false) from INSERT when they have a default, which made draft users get
	// password_set=true from the DB default. The DB column still has DEFAULT true.
	Active      bool `gorm:"not null"`
	PasswordSet bool `gorm:"not null"`
	CreatedAt   time.Time
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

func (r *UserRepositoryPG) dbx(ctx context.Context) *gorm.DB {
	return database.TxFromCtx(ctx, r.db).WithContext(ctx)
}

func RunMigrations(_ *gorm.DB) error {
	// Schema is managed with external SQL migrations.
	return nil
}

func (r *UserRepositoryPG) Save(ctx context.Context, user *userdomain.User) error {
	model := toModel(user)
	if err := r.dbx(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	user.ID = model.ID
	user.CreatedAt = model.CreatedAt
	return nil
}

func (r *UserRepositoryPG) FindByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	var model UserModel
	err := r.dbx(ctx).Where("email = ?", email).First(&model).Error
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
	err := r.dbx(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, userdomain.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return toDomain(&model), nil
}

func (r *UserRepositoryPG) Update(ctx context.Context, user *userdomain.User) error {
	err := r.dbx(ctx).Model(&UserModel{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"name":         user.Name,
		"email":        user.Email,
		"role":         string(user.Role),
		"password":     user.Password,
		"active":       user.Active,
		"password_set": user.PasswordSet,
	}).Error
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *UserRepositoryPG) Delete(ctx context.Context, id string) error {
	res := r.dbx(ctx).Where("id = ?", id).Delete(&UserModel{})
	if res.Error != nil {
		return fmt.Errorf("delete user: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return userdomain.ErrUserNotFound
	}
	return nil
}

func applyUserListFilter(q *gorm.DB, filter userports.UserListFilter) *gorm.DB {
	if s := strings.TrimSpace(filter.Query); s != "" {
		like := "%" + strings.ToLower(s) + "%"
		q = q.Where("(LOWER(name) LIKE ? OR LOWER(email) LIKE ?)", like, like)
	}
	return q
}

func (r *UserRepositoryPG) FindAll(ctx context.Context, filter userports.UserListFilter, params pagination.Params) ([]userdomain.User, int64, error) {
	base := applyUserListFilter(r.dbx(ctx).Model(&UserModel{}), filter)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	var models []UserModel
	err := applyUserListFilter(r.dbx(ctx), filter).
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

func (r *UserRepositoryPG) FindAllActiveMembers(ctx context.Context) ([]userports.Member, error) {
	var rows []struct {
		ID    string
		Email string
	}
	if err := r.dbx(ctx).Model(&UserModel{}).
		Where("active = ?", true).
		Select("id", "email").
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("find all active members: %w", err)
	}
	out := make([]userports.Member, len(rows))
	for i, row := range rows {
		out[i] = userports.Member{ID: row.ID, Email: row.Email}
	}
	return out, nil
}

func toModel(u *userdomain.User) *UserModel {
	return &UserModel{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		Password:    u.Password,
		Role:        string(u.Role),
		Active:      u.Active,
		PasswordSet: u.PasswordSet,
	}
}

func toDomain(m *UserModel) *userdomain.User {
	return &userdomain.User{
		ID:          m.ID,
		Name:        m.Name,
		Email:       m.Email,
		Password:    m.Password,
		Role:        userdomain.Role(m.Role),
		Active:      m.Active,
		PasswordSet: m.PasswordSet,
		CreatedAt:   m.CreatedAt,
	}
}
