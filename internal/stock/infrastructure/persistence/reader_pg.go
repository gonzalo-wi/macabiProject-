package stockpersistence

import (
	"context"
)

func (r *StockRepositoryPG) FindProjectCoordinators(ctx context.Context, projectID string) ([]string, error) {
	var userIDs []string
	err := r.db.WithContext(ctx).
		Table("project_members").
		Where("project_id = ? AND role = ?", projectID, "coordinator").
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}

func (r *StockRepositoryPG) IsProjectMember(ctx context.Context, projectID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("project_members").
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *StockRepositoryPG) IsProjectCoordinator(ctx context.Context, projectID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("project_members").
		Where("project_id = ? AND user_id = ? AND role = ?", projectID, userID, "coordinator").
		Count(&count).Error
	return count > 0, err
}

func (r *StockRepositoryPG) FindEmailByID(ctx context.Context, userID string) (string, error) {
	var email string
	err := r.db.WithContext(ctx).
		Table("users").
		Select("email").
		Where("id = ?", userID).
		Scan(&email).Error
	return email, err
}

func (r *StockRepositoryPG) FindEmailsByIDs(ctx context.Context, userIDs []string) (map[string]string, error) {
	if len(userIDs) == 0 {
		return map[string]string{}, nil
	}
	var rows []struct {
		ID    string
		Email string
	}
	err := r.db.WithContext(ctx).
		Table("users").
		Select("id, email").
		Where("id IN ?", userIDs).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, len(rows))
	for _, row := range rows {
		result[row.ID] = row.Email
	}
	return result, nil
}
