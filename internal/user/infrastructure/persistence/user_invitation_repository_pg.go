package userpersistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"macabi-back/internal/shared/database"
	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"

	"gorm.io/gorm"
)

type UserInvitationModel struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string `gorm:"index;not null"`
	Name      string `gorm:"not null"`
	Role      string `gorm:"not null;default:'user'"`
	TokenHash string `gorm:"not null"`
	ExpiresAt *time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

func (UserInvitationModel) TableName() string {
	return "user_invitations"
}

type UserInvitationRepositoryPG struct {
	db *gorm.DB
}

func NewUserInvitationRepositoryPG(db *gorm.DB) *UserInvitationRepositoryPG {
	return &UserInvitationRepositoryPG{db: db}
}

func (r *UserInvitationRepositoryPG) dbx(ctx context.Context) *gorm.DB {
	return database.TxFromCtx(ctx, r.db).WithContext(ctx)
}

func (r *UserInvitationRepositoryPG) InvalidatePendingByEmail(ctx context.Context, email string) error {
	now := time.Now()
	res := r.dbx(ctx).Model(&UserInvitationModel{}).
		Where("email = ? AND used_at IS NULL", email).
		Update("used_at", now)
	if res.Error != nil {
		return fmt.Errorf("invalidate invitations: %w", res.Error)
	}
	return nil
}

func (r *UserInvitationRepositoryPG) Create(ctx context.Context, email, name string, role userdomain.Role, tokenHash string, expiresAt time.Time) error {
	exp := expiresAt
	row := UserInvitationModel{
		Email:     email,
		Name:      name,
		Role:      string(role),
		TokenHash: tokenHash,
		ExpiresAt: &exp,
	}
	if err := r.dbx(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("create invitation: %w", err)
	}
	return nil
}

func (r *UserInvitationRepositoryPG) FindPendingByID(ctx context.Context, id string) (*userports.PendingUserInvitation, error) {
	var row UserInvitationModel
	err := r.dbx(ctx).Where("id = ? AND used_at IS NULL", id).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("find pending invitation: %w", err)
	}
	return modelToPending(&row), nil
}

func (r *UserInvitationRepositoryPG) FindPendingByEmail(ctx context.Context, email string) (*userports.PendingUserInvitation, error) {
	var row UserInvitationModel
	err := r.dbx(ctx).Where("email = ? AND used_at IS NULL", email).Order("created_at DESC").First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("find pending invitation by email: %w", err)
	}
	return modelToPending(&row), nil
}

func (r *UserInvitationRepositoryPG) ListPending(ctx context.Context) ([]userports.PendingUserInvitation, error) {
	var rows []UserInvitationModel
	err := r.dbx(ctx).Where("used_at IS NULL").Order("created_at DESC").Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list pending invitations: %w", err)
	}
	out := make([]userports.PendingUserInvitation, 0, len(rows))
	for i := range rows {
		out = append(out, *modelToPending(&rows[i]))
	}
	return out, nil
}

func modelToPending(row *UserInvitationModel) *userports.PendingUserInvitation {
	p := &userports.PendingUserInvitation{
		ID:        row.ID,
		Email:     row.Email,
		Name:      row.Name,
		Role:      userdomain.Role(row.Role),
		CreatedAt: row.CreatedAt,
	}
	if row.ExpiresAt != nil {
		p.ExpiresAt = *row.ExpiresAt
	}
	return p
}

func (r *UserInvitationRepositoryPG) FindValidByTokenHash(ctx context.Context, tokenHash string) (*userports.UserInvitation, error) {
	var row UserInvitationModel
	err := r.dbx(ctx).
		Where("token_hash = ? AND used_at IS NULL", tokenHash).
		Where("expires_at IS NULL OR expires_at > ?", time.Now()).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("find invitation: %w", err)
	}
	role := userdomain.Role(row.Role)
	return &userports.UserInvitation{
		ID:    row.ID,
		Email: row.Email,
		Name:  row.Name,
		Role:  role,
	}, nil
}

func (r *UserInvitationRepositoryPG) MarkUsed(ctx context.Context, id string) error {
	now := time.Now()
	res := r.dbx(ctx).Model(&UserInvitationModel{}).
		Where("id = ? AND used_at IS NULL", id).
		Update("used_at", now)
	if res.Error != nil {
		return fmt.Errorf("mark invitation used: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("invitation already used or missing")
	}
	return nil
}
