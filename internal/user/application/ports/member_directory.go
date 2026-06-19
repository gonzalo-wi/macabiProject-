package userports

import "context"

// Member es un destinatario mínimo para notificaciones globales (no ligadas a proyecto).
type Member struct {
	ID    string
	Email string
}

// MemberDirectory expone a todos los miembros activos del sistema.
// Lo usa el módulo de noticias para notificar a todos al publicar.
type MemberDirectory interface {
	FindAllActiveMembers(ctx context.Context) ([]Member, error)
}
