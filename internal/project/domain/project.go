package projectdomain

import (
	"strings"
	"time"
)

type Project struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
}

func NewProject(name, description string) (*Project, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyName
	}
	return &Project{
		Name:        name,
		Description: description,
	}, nil
}

type ProjectMemberRole string

const (
	MemberRoleCoordinator ProjectMemberRole = "coordinator"
	MemberRoleMadrij      ProjectMemberRole = "madrij"
)

type ProjectMember struct {
	ID        string
	ProjectID string
	UserID    string
	Role      ProjectMemberRole
	CreatedAt time.Time
}
