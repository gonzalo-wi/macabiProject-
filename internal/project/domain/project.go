package projectdomain

import (
	"strings"
	"time"
)

type Project struct {
	ID          string
	Name        string
	Description string
	AdminUserID string
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProject(name, description, adminUserID string) (*Project, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyName
	}
	if adminUserID == "" {
		return nil, ErrMissingAdmin
	}
	return &Project{
		Name:        name,
		Description: description,
		AdminUserID: adminUserID,
		Active:      true,
	}, nil
}
