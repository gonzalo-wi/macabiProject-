package deps

import (
	projectusecases "macabi-back/internal/project/application/usecases"
	projecthttp "macabi-back/internal/project/infrastructure/http"
	projectpersistence "macabi-back/internal/project/infrastructure/persistence"

	"gorm.io/gorm"
)

func buildProjectDeps(db *gorm.DB) *projecthttp.ProjectHandler {
	projectRepo := projectpersistence.NewProjectRepositoryPG(db)

	return projecthttp.NewProjectHandler(
		projectusecases.NewCreateProject(projectRepo),
		projectusecases.NewListProjects(projectRepo),
		projectusecases.NewGetProject(projectRepo),
		projectusecases.NewUpdateProject(projectRepo),
		projectusecases.NewDeleteProject(projectRepo),
		projectusecases.NewAddProjectMember(projectRepo),
		projectusecases.NewRemoveProjectMember(projectRepo),
		projectusecases.NewListProjectMembers(projectRepo),
	)
}
