package userdto

import (
	"time"

	userdomain "macabi-back/internal/user/domain"
)

type ChangeRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type SetUserStatusRequest struct {
	Active bool `json:"active"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type UserResponse struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Email               string    `json:"email"`
	Role                string    `json:"role"`
	Active              bool      `json:"active"`
	PasswordSet         bool      `json:"password_set"`
	InvitationStatus    string    `json:"invitation_status"`
	PendingInvitationID *string   `json:"pending_invitation_id,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
}

type PaginatedUserResponse struct {
	Data       []UserResponse `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

func ToUserResponse(u *userdomain.User) UserResponse {
	return ToUserResponseWithInvitation(u, "", userdomain.InvitationStatus(u, false))
}

func ToUserResponseWithInvitation(u *userdomain.User, pendingInvitationID, invitationStatus string) UserResponse {
	var pendingID *string
	if pendingInvitationID != "" {
		pendingID = &pendingInvitationID
	}
	return UserResponse{
		ID:                  u.ID,
		Name:                u.Name,
		Email:               u.Email,
		Role:                u.Role.String(),
		Active:              u.Active,
		PasswordSet:         u.PasswordSet,
		InvitationStatus:    invitationStatus,
		PendingInvitationID: pendingID,
		CreatedAt:           u.CreatedAt,
	}
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role"`
}
