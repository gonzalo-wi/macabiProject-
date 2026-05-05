package attendanceusecases

import (
	"context"

	attendanceports "macabi-back/internal/attendance/application/ports"
)

type GetAttendanceCountInput struct {
	ProjectID string
}

type GetAttendanceCountOutput struct {
	ProjectID string
	Confirmed int
}

type GetAttendanceCount struct {
	repo attendanceports.AttendanceRepository
}

func NewGetAttendanceCount(repo attendanceports.AttendanceRepository) *GetAttendanceCount {
	return &GetAttendanceCount{repo: repo}
}

func (uc *GetAttendanceCount) Execute(ctx context.Context, input GetAttendanceCountInput) (*GetAttendanceCountOutput, error) {
	count, err := uc.repo.CountByProject(ctx, input.ProjectID)
	if err != nil {
		return nil, err
	}
	return &GetAttendanceCountOutput{
		ProjectID: input.ProjectID,
		Confirmed: count,
	}, nil
}