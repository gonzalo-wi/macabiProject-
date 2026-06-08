package projectpersistence

import (
	"strings"

	"gorm.io/gorm"

	projectports "macabi-back/internal/project/application/ports"
)

type ProjectRepositoryPG struct {
	db *gorm.DB
}

var _ projectports.ProjectRepository = (*ProjectRepositoryPG)(nil)

func NewProjectRepositoryPG(db *gorm.DB) *ProjectRepositoryPG {
	return &ProjectRepositoryPG{db: db}
}

func RunMigrations(_ *gorm.DB) error {
	return nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "duplicate key")
}
