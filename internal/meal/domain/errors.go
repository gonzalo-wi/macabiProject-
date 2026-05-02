package mealdomain

import "errors"

var (
	ErrInvalidMealType       = errors.New("tipo de vianda inválido")
	ErrInvalidCategory       = errors.New("categoría inválida")
	ErrEmptyTitle            = errors.New("el título no puede estar vacío")
	ErrInvalidDate           = errors.New("las viandas solo están disponibles sábados y domingos")
	ErrMealNotFound          = errors.New("vianda no encontrada")
	ErrMealSoldOut           = errors.New("vianda agotada")
	ErrBookingNotFound       = errors.New("reserva no encontrada")
	ErrBookingNotOwned       = errors.New("no puedes cancelar una reserva que no es tuya")
	ErrInvalidAvailableCount = errors.New("la cantidad disponible no puede ser negativa")
	ErrBookingDeadlinePassed = errors.New("el plazo para reservar ya cerró (viernes 23:59)")
	ErrTemplateNotFound      = errors.New("template de vianda no encontrado")
	ErrAlreadyBooked         = errors.New("ya tenés una reserva para este día; no podés reservar almuerzo y cena el mismo día")
	ErrProjectConflict       = errors.New("ya tenés una reserva en otro proyecto para ese día; no podés asistir a dos proyectos el mismo día")
	ErrGarnishRequired       = errors.New("esta vianda requiere que selecciones una guarnición")
	ErrGarnishNotForMeal     = errors.New("la guarnición seleccionada no pertenece a esta vianda")
	ErrGarnishNotFound       = errors.New("guarnición no encontrada")
	ErrEmptyGarnishName      = errors.New("el nombre de la guarnición no puede estar vacío")
)
