package attendancepersistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	attendanceports "macabi-back/internal/attendance/application/ports"
	attendancedomain "macabi-back/internal/attendance/domain"
)

type AttendanceModel struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string `gorm:"type:uuid;not null"`
	ProjectID string `gorm:"type:uuid;not null"`
	CreatedAt time.Time
}

func (AttendanceModel) TableName() string { return "attendances" }

type AttendanceRepositoryPG struct {
	db *gorm.DB
}

func NewAttendanceRepositoryPG(db *gorm.DB) *AttendanceRepositoryPG {
	return &AttendanceRepositoryPG{db: db}
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&AttendanceModel{})
}

func (r *AttendanceRepositoryPG) Save(ctx context.Context, a *attendancedomain.Attendance) error {
	m := AttendanceModel{
		UserID:    a.UserID,
		ProjectID: a.ProjectID,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	a.ID = m.ID
	a.CreatedAt = m.CreatedAt
	return nil
}

func (r *AttendanceRepositoryPG) FindByUserAndProject(ctx context.Context, userID, projectID string) (*attendancedomain.Attendance, error) {
	var m AttendanceModel
	if err := r.db.WithContext(ctx).First(&m, "user_id = ? AND project_id = ?", userID, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, attendancedomain.ErrNotFound
		}
		return nil, err
	}
	return &attendancedomain.Attendance{
		ID:        m.ID,
		UserID:    m.UserID,
		ProjectID: m.ProjectID,
		CreatedAt: m.CreatedAt,
	}, nil
}

func (r *AttendanceRepositoryPG) CountByProject(ctx context.Context, projectID string) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&AttendanceModel{}).Where("project_id = ?", projectID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *AttendanceRepositoryPG) DeleteByUserAndProject(ctx context.Context, userID, projectID string) error {
	return r.db.WithContext(ctx).Delete(&AttendanceModel{}, "user_id = ? AND project_id = ?", userID, projectID).Error
}

var _ attendanceports.AttendanceRepository = (*AttendanceRepositoryPG)(nil)