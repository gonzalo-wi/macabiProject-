package eventpersistence

import (
	"context"
	"time"

	eventdomain "macabi-back/internal/event/domain"
	eventports "macabi-back/internal/event/application/ports"
)

func (r *RepositoryPG) FindEventsWithDeadlineTomorrow(ctx context.Context) ([]eventdomain.EventInstance, error) {
	now := time.Now().UTC()
	startOfTomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	endOfTomorrow := startOfTomorrow.Add(24 * time.Hour)

	var rows []EventInstanceModel
	err := r.db.WithContext(ctx).
		Where("status = ? AND response_deadline_at >= ? AND response_deadline_at < ?",
			string(eventdomain.EventStatusOpen), startOfTomorrow, endOfTomorrow).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]eventdomain.EventInstance, len(rows))
	for i := range rows {
		out[i] = *toDomainInst(&rows[i])
	}
	return out, nil
}

func (r *RepositoryPG) FindUsersWithoutResponse(ctx context.Context, eventID string) ([]eventports.ReminderTarget, error) {
	type row struct {
		UserID string `gorm:"column:user_id"`
		Name   string `gorm:"column:name"`
		Email  string `gorm:"column:email"`
	}

	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT DISTINCT u.id AS user_id, u.name, u.email
		FROM event_instance_projects eip
		JOIN project_members pm ON pm.project_id = eip.project_id
		JOIN users u ON u.id = pm.user_id
		WHERE eip.event_instance_id = ?
		  AND u.active = true
		  AND u.id NOT IN (
		    SELECT user_id FROM event_responses WHERE event_instance_id = ?
		  )
	`, eventID, eventID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]eventports.ReminderTarget, len(rows))
	for i, r := range rows {
		out[i] = eventports.ReminderTarget{UserID: r.UserID, Name: r.Name, Email: r.Email}
	}
	return out, nil
}
