package eventdomain

import "errors"

var (
	ErrNotFound               = errors.New("recurso no encontrado")
	ErrInvalidEventType       = errors.New("tipo de evento inválido")
	ErrInvalidEventStatus     = errors.New("estado de evento inválido")
	ErrInvalidModuleType      = errors.New("tipo de módulo inválido")
	ErrInvalidOptionGroupType = errors.New("tipo de grupo de opciones inválido")
	ErrOptionCapacityReached  = errors.New("cupo de opción agotado")
	ErrDuplicateResponse      = errors.New("ya existe una respuesta para este evento")
	ErrEmptyTitle             = errors.New("el título no puede estar vacío")
)
