package mealdomain

import "errors"

var (
	ErrInvalidMealType       = errors.New("tipo de vianda inválido")
	ErrInvalidCategory       = errors.New("categoría inválida")
	ErrEmptyTitle            = errors.New("el título no puede estar vacío")
	ErrInvalidDate           = errors.New("las viandas solo están disponibles los sábados")
	ErrMealNotFound          = errors.New("vianda no encontrada")
	ErrMealSoldOut           = errors.New("vianda agotada")
	ErrBookingNotFound       = errors.New("reserva no encontrada")
	ErrBookingNotOwned       = errors.New("no puedes cancelar una reserva que no es tuya")
	ErrInvalidAvailableCount = errors.New("la cantidad disponible no puede ser negativa")
	ErrBookingDeadlinePassed = errors.New("el plazo para reservar ya cerró (viernes 23:59)")
	ErrTemplateNotFound      = errors.New("template de vianda no encontrado")
)
