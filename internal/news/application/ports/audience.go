package newsports

import (
	"context"

	userports "macabi-back/internal/user/application/ports"
)

// ProjectAudienceReader resuelve la audiencia por proyecto de las noticias.
// Lo satisface el repositorio que conoce project_members + users.
type ProjectAudienceReader interface {
	// FindUserProjectIDs devuelve los proyectos a los que pertenece el usuario.
	FindUserProjectIDs(ctx context.Context, userID string) ([]string, error)
	// FindActiveMembersOfProjects devuelve (ID+email) de los miembros activos
	// (madrijim + coordinadores) de los proyectos dados, sin duplicados.
	FindActiveMembersOfProjects(ctx context.Context, projectIDs []string) ([]userports.Member, error)
}
