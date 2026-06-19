package stockpersistence

import (
	"context"

	userports "macabi-back/internal/user/application/ports"
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

// FindUserProjectIDs devuelve los proyectos a los que pertenece el usuario.
func (r *StockRepositoryPG) FindUserProjectIDs(ctx context.Context, userID string) ([]string, error) {
	var ids []string
	err := r.db.WithContext(ctx).
		Table("project_members").
		Where("user_id = ?", userID).
		Pluck("project_id", &ids).Error
	return ids, err
}

// FindActiveMembersOfProjects devuelve (ID+email) de los miembros activos
// (madrijim + coordinadores) de los proyectos dados, sin duplicados.
func (r *StockRepositoryPG) FindActiveMembersOfProjects(ctx context.Context, projectIDs []string) ([]userports.Member, error) {
	if len(projectIDs) == 0 {
		return []userports.Member{}, nil
	}
	var rows []struct {
		ID    string
		Email string
	}
	err := r.db.WithContext(ctx).
		Table("project_members pm").
		Select("DISTINCT u.id, u.email").
		Joins("JOIN users u ON u.id = pm.user_id").
		Where("pm.project_id IN ? AND u.active = ?", projectIDs, true).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]userports.Member, len(rows))
	for i, row := range rows {
		out[i] = userports.Member{ID: row.ID, Email: row.Email}
	}
	return out, nil
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
