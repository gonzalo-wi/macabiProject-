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
	Capacity    int
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProject(name, description, adminUserID string, capacity int) (*Project, error) {
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
		Capacity:    capacity,
		Active:      true,
	}, nil
}
