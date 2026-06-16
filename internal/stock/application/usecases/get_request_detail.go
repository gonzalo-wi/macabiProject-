package stockusecases

import (
	"context"

	projectports "macabi-back/internal/project/application/ports"
	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
	userdomain "macabi-back/internal/user/domain"
)

type GetRequestDetail struct {
	repo     stockports.StockRepository
	projects projectports.ProjectMembership
}

func NewGetRequestDetail(repo stockports.StockRepository, projects projectports.ProjectMembership) *GetRequestDetail {
	return &GetRequestDetail{repo: repo, projects: projects}
}

type GetRequestDetailInput struct {
	RequestID string
	UserID    string
	UserRole  string
}

func (uc *GetRequestDetail) Execute(ctx context.Context, input GetRequestDetailInput) (*stockdomain.RequestDetail, error) {
	d, err := uc.repo.FindRequestDetailByID(ctx, input.RequestID)
	if err != nil {
		return nil, err
	}
	if userdomain.Role(input.UserRole) == userdomain.RoleAdmin {
		return d, nil
	}
	if d.Request.RequestedByID == input.UserID {
		return d, nil
	}
	ok, err := uc.projects.IsProjectCoordinator(ctx, d.Request.ProjectID, input.UserID)
	if err != nil {
		return nil, err
	}
	if ok {
		return d, nil
	}
	return nil, stockdomain.ErrForbidden
}
