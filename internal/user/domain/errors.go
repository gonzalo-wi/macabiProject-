package userdomain

import "errors"

var (
	ErrInvalidRole        = errors.New("rol inválido")
	ErrInvalidEmail       = errors.New("formato de email inválido")
	ErrWeakPassword       = errors.New("la contraseña debe tener al menos 8 caracteres")
	ErrEmptyName          = errors.New("el nombre no puede estar vacío")
	ErrUserNotFound       = errors.New("usuario no encontrado")
	ErrEmailAlreadyTaken  = errors.New("el email ya está registrado")
	ErrInvalidCredentials = errors.New("email o contraseña incorrectos")
	ErrUnauthorized       = errors.New("no autorizado")
	ErrForbidden          = errors.New("permisos insuficientes")
)
