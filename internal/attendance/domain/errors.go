package attendancedomain

import "errors"

var (
	ErrAlreadyConfirmed = errors.New("ya confirmaste tu asistencia a este evento")
	ErrNotFound         = errors.New("asistencia no encontrada")
)