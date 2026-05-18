package stockusecases

import (
	"context"

	"macabi-back/internal/shared/pagination"
	stockports "macabi-back/internal/stock/application/ports"
	stockdomain "macabi-back/internal/stock/domain"
	userdomain "macabi-back/internal/user/domain"
)

type ListRequests struct {
	repo     stockports.StockRepository
	projects stockports.ProjectMemberReader
}

func NewListRequests(repo stockports.StockRepository, projects stockports.ProjectMemberReader) *ListRequests {
	return &ListRequests{repo: repo, projects: projects}
}

type ListRequestsInput struct {
	Params    pagination.Params
	ProjectID string
	UserID    string
	UserRole  string
}

func (uc *ListRequests) Execute(ctx context.Context, input ListRequestsInput) (pagination.Result[stockdomain.RequestDetail], error) {
	role := userdomain.Role(input.UserRole)

	if role == userdomain.RoleAdmin {
		var only *string
		return uc.repo.ListRequests(ctx, input.Params, input.ProjectID, only)
	}

	if input.ProjectID == "" {
		return pagination.Result[stockdomain.RequestDetail]{}, stockdomain.ErrMissingRequiredField
	}

	ok, err := uc.projects.IsProjectMember(ctx, input.ProjectID, input.UserID)
	if err != nil {
		return pagination.Result[stockdomain.RequestDetail]{}, err
	}
	if !ok {
		return pagination.Result[stockdomain.RequestDetail]{}, stockdomain.ErrForbidden
	}

	coord, err := uc.projects.IsProjectCoordinator(ctx, input.ProjectID, input.UserID)
	if err != nil {
		return pagination.Result[stockdomain.RequestDetail]{}, err
	}

	var filter *string
	if !coord {
		u := input.UserID
		filter = &u
	}
	return uc.repo.ListRequests(ctx, input.Params, input.ProjectID, filter)
}
