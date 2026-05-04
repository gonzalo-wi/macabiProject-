package attendancehttp

import (
	attendanceusecases "macabi-back/internal/attendance/application/usecases"
)

type ConfirmAttendanceRequest struct {
	ProjectID string `json:"project_id" binding:"required"`
}

func (r ConfirmAttendanceRequest) toInput(userID string) attendanceusecases.ConfirmAttendanceInput {
	return attendanceusecases.ConfirmAttendanceInput{
		UserID:    userID,
		ProjectID: r.ProjectID,
	}
}

type AttendanceResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	ProjectID string `json:"project_id"`
}

type AttendanceCountResponse struct {
	ProjectID string `json:"project_id"`
	Confirmed int    `json:"confirmed"`
	Capacity  int    `json:"capacity"`
}