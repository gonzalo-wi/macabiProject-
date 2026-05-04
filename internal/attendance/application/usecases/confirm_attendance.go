package attendanceusecases

import (
	"context"
	"errors"

	attendanceports "macabi-back/internal/attendance/application/ports"
	attendancedomain "macabi-back/internal/attendance/domain"
)

type ConfirmAttendanceInput struct {
	UserID    string
	ProjectID string
}

type ConfirmAttendance struct {
	repo attendanceports.AttendanceRepository
}

func NewConfirmAttendance(repo attendanceports.AttendanceRepository) *ConfirmAttendance {
	return &ConfirmAttendance{repo: repo}
}

func (uc *ConfirmAttendance) Execute(ctx context.Context, input ConfirmAttendanceInput) (*attendancedomain.Attendance, error) {
	_, err := uc.repo.FindByUserAndProject(ctx, input.UserID, input.ProjectID)
	if err == nil {
		return nil, attendancedomain.ErrAlreadyConfirmed
	}
	if !errors.Is(err, attendancedomain.ErrNotFound) {
		return nil, err
	}

	a := attendancedomain.NewAttendance(input.UserID, input.ProjectID)
	if err := uc.repo.Save(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}