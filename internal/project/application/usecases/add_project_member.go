package projectusecases

import (
	"context"

	projectports "macabi-back/internal/project/application/ports"
	projectdomain "macabi-back/internal/project/domain"
)

type AddProjectMemberInput struct {
	ProjectID string
	UserID    string
	Role      string
}

type AddProjectMember struct {
	repo projectports.ProjectRepository
}

func NewAddProjectMember(repo projectports.ProjectRepository) *AddProjectMember {
	return &AddProjectMember{repo: repo}
}

func (uc *AddProjectMember) Execute(ctx context.Context, in AddProjectMemberInput) (*projectdomain.ProjectMember, error) {
	if _, err := uc.repo.FindByID(ctx, in.ProjectID); err != nil {
		return nil, err
	}
	role := projectdomain.ProjectMemberRole(in.Role)
	if role != projectdomain.MemberRoleCoordinator && role != projectdomain.MemberRoleMadrij {
		return nil, projectdomain.ErrInvalidMemberRole
	}
	m := &projectdomain.ProjectMember{
		ProjectID: in.ProjectID,
		UserID:    in.UserID,
		Role:      role,
	}
	if err := uc.repo.AddMember(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}
