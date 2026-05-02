package projectdomain

import "errors"

var (
	ErrEmptyName       = errors.New("el nombre del proyecto no puede estar vacío")
	ErrMissingAdmin    = errors.New("el proyecto debe tener un coordinador asignado")
	ErrProjectNotFound = errors.New("proyecto no encontrado")
)
