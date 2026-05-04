package attendanceports

import (
	"context"

	attendancedomain "macabi-back/internal/attendance/domain"
)

type AttendanceRepository interface {
	Save(ctx context.Context, a *attendancedomain.Attendance) error
	FindByUserAndProject(ctx context.Context, userID, projectID string) (*attendancedomain.Attendance, error)
	CountByProject(ctx context.Context, projectID string) (int, error)
	DeleteByUserAndProject(ctx context.Context, userID, projectID string) error
}