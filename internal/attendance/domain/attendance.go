package attendancedomain

import "time"

type Attendance struct {
	ID        string
	UserID    string
	ProjectID string
	CreatedAt time.Time
}

func NewAttendance(userID, projectID string) *Attendance {
	return &Attendance{
		UserID:    userID,
		ProjectID: projectID,
	}
}