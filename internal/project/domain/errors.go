package projectdomain

import "errors"

var (
	ErrEmptyName          = errors.New("el nombre del proyecto no puede estar vacío")
	ErrProjectNotFound    = errors.New("proyecto no encontrado")
	ErrMemberNotFound     = errors.New("miembro del proyecto no encontrado")
	ErrDuplicateMember    = errors.New("el usuario ya es miembro de este proyecto")
	ErrInvalidMemberRole  = errors.New("rol de miembro inválido")
	ErrMissingCoordinator = errors.New("el proyecto debe tener un coordinador")
)
