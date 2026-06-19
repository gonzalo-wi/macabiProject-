package newsdomain

import "errors"

var (
	ErrNewsNotFound            = errors.New("noticia no encontrada")
	ErrForbidden               = errors.New("no tiene permisos para realizar esta acción")
	ErrMissingRequiredField    = errors.New("faltan campos obligatorios")
	ErrInvalidStatusTransition = errors.New("transición de estado inválida")
	ErrImageStorageUnavailable = errors.New("almacenamiento de imágenes no configurado")
	ErrInvalidImagePath        = errors.New("ruta de imagen inválida")
	ErrInvalidMimeType         = errors.New("tipo de archivo no permitido para la imagen")
	ErrImageTooLarge           = errors.New("la imagen no puede superar 2 MB")
	ErrImageAttachFailed       = errors.New("no se pudo guardar la imagen; la noticia no fue creada")
	ErrNoImageAttached         = errors.New("esta noticia no tiene imagen adjunta")
	ErrNotificationNotFound    = errors.New("notificación no encontrada")
)
