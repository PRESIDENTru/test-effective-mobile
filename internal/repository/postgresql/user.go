package postgresql

import (
	"context"
	"fmt"
	"testJob/internal/models"

	"github.com/google/uuid"
)

func (s *Storage) SaveUser(ctx context.Context, user *models.User) error {
	const operation = "repository.postgresql.SaveUser"
	req := `INSERT INTO users (id) VALUES ($1)`
	_, err := s.db.ExecContext(ctx, req, user.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return err
}

func (s *Storage) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	const operation = "repository.postgresql.GetUserByID"
	var result *models.User
	err := s.db.GetContext(ctx, &result, `SELECT * FROM users WHERE id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return result, err
}

func (s *Storage) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	const operation = "repository.postgresql.DeleteUser"
	_, err := s.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return err
}

func (s *Storage) ListUsers(ctx context.Context) (*[]models.User, error) {
	const operation = "repository.postgresql.ListUsers"
	var result []models.User
	err := s.db.SelectContext(ctx, &result, `SELECT * FROM users`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return &result, err
}
