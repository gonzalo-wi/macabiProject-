package expensesdomain

import "errors"

var (
	ErrExpenseNotFound           = errors.New("gasto no encontrado")
	ErrNoReceiptAttached         = errors.New("este gasto no tiene comprobante adjunto")
	ErrForbidden                 = errors.New("no tiene permisos para realizar esta acción")
	ErrMissingRequiredField      = errors.New("faltan campos obligatorios")
	ErrInvalidAmount             = errors.New("monto inválido")
	ErrInvalidStatusTransition   = errors.New("transición de estado inválida")
	ErrReceiptStorageUnavailable = errors.New("almacenamiento de comprobantes no configurado")
	ErrInvalidReceiptPath        = errors.New("ruta de comprobante inválida")
	ErrInvalidMimeType           = errors.New("tipo de archivo no permitido para el comprobante")
	ErrReceiptTooLarge           = errors.New("el comprobante no puede superar 2 MB")
	ErrReceiptAttachFailed       = errors.New("no se pudo guardar el comprobante; el gasto no fue creado")
	ErrNotificationNotFound      = errors.New("notificación no encontrada")
)
