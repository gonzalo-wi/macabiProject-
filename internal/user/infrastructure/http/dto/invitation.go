package userdto

import (
	"time"

	userports "macabi-back/internal/user/application/ports"
)

type CreateUserInvitationRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
	Role  string `json:"role"`
}

type PendingInvitationResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type ListPendingInvitationsResponse struct {
	Data []PendingInvitationResponse `json:"data"`
}

func ToPendingInvitationResponse(p userports.PendingUserInvitation) PendingInvitationResponse {
	return PendingInvitationResponse{
		ID:        p.ID,
		Email:     p.Email,
		Name:      p.Name,
		Role:      p.Role.String(),
		ExpiresAt: p.ExpiresAt,
		CreatedAt: p.CreatedAt,
	}
}
