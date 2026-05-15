package stockdomain

import "errors"

var (
	ErrEmptyResourceName       = errors.New("el nombre del recurso no puede estar vacío")
	ErrInvalidStock            = errors.New("el stock no puede ser negativo")
	ErrInvalidResourceType     = errors.New("tipo de recurso inválido, debe ser 'returnable' o 'consumable'")
	ErrInvalidQuantity         = errors.New("la cantidad debe ser mayor a cero")
	ErrMissingRequiredField    = errors.New("faltan campos obligatorios")
	ErrResourceNotFound        = errors.New("recurso no encontrado")
	ErrRequestNotFound         = errors.New("solicitud no encontrada")
	ErrNotificationNotFound    = errors.New("notificación no encontrada")
	ErrInsufficientStock       = errors.New("stock disponible insuficiente")
	ErrInvalidStatusTransition = errors.New("transición de estado inválida para la solicitud")
	ErrForbidden               = errors.New("no tiene permisos para realizar esta acción")
	ErrResourceInUse           = errors.New("el recurso tiene solicitudes activas y no puede eliminarse")
	ErrNotReturnable           = errors.New("el recurso es consumible y no puede marcarse como devuelto")
)
