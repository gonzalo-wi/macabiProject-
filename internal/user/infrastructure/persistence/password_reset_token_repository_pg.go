package userpersistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"macabi-back/internal/shared/database"
	userports "macabi-back/internal/user/application/ports"

	"gorm.io/gorm"
)

type PasswordResetTokenModel struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string `gorm:"type:uuid;not null;index"`
	TokenHash string `gorm:"uniqueIndex;not null;size:64"`
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

func (PasswordResetTokenModel) TableName() string {
	return "password_reset_tokens"
}

type PasswordResetTokenRepositoryPG struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepositoryPG(db *gorm.DB) *PasswordResetTokenRepositoryPG {
	return &PasswordResetTokenRepositoryPG{db: db}
}

func (r *PasswordResetTokenRepositoryPG) dbx(ctx context.Context) *gorm.DB {
	return database.TxFromCtx(ctx, r.db).WithContext(ctx)
}

func (r *PasswordResetTokenRepositoryPG) InvalidateUnusedByUserID(ctx context.Context, userID string) error {
	now := time.Now()
	res := r.dbx(ctx).Model(&PasswordResetTokenModel{}).
		Where("user_id = ? AND used_at IS NULL", userID).
		Update("used_at", now)
	if res.Error != nil {
		return fmt.Errorf("invalidate reset tokens: %w", res.Error)
	}
	return nil
}

func (r *PasswordResetTokenRepositoryPG) Create(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	row := PasswordResetTokenModel{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
	if err := r.dbx(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("create reset token: %w", err)
	}
	return nil
}

func (r *PasswordResetTokenRepositoryPG) FindValidByTokenHash(ctx context.Context, tokenHash string) (*userports.PasswordResetToken, error) {
	var row PasswordResetTokenModel
	err := r.dbx(ctx).
		Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", tokenHash, time.Now()).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("find reset token: %w", err)
	}
	return &userports.PasswordResetToken{ID: row.ID, UserID: row.UserID}, nil
}

func (r *PasswordResetTokenRepositoryPG) MarkUsed(ctx context.Context, id string) error {
	now := time.Now()
	res := r.dbx(ctx).Model(&PasswordResetTokenModel{}).
		Where("id = ? AND used_at IS NULL", id).
		Update("used_at", now)
	if res.Error != nil {
		return fmt.Errorf("mark reset token used: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("reset token already used or missing")
	}
	return nil
}
