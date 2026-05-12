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

type RemoveProjectMember struct {
	repo projectports.ProjectRepository
}

func NewRemoveProjectMember(repo projectports.ProjectRepository) *RemoveProjectMember {
	return &RemoveProjectMember{repo: repo}
}

func (uc *RemoveProjectMember) Execute(ctx context.Context, projectID, userID string) error {
	return uc.repo.RemoveMember(ctx, projectID, userID)
}

type ListProjectMembers struct {
	repo projectports.ProjectRepository
}

func NewListProjectMembers(repo projectports.ProjectRepository) *ListProjectMembers {
	return &ListProjectMembers{repo: repo}
}

func (uc *ListProjectMembers) Execute(ctx context.Context, projectID string) ([]projectdomain.ProjectMember, error) {
	if _, err := uc.repo.FindByID(ctx, projectID); err != nil {
		return nil, err
	}
	return uc.repo.ListMembers(ctx, projectID)
}
