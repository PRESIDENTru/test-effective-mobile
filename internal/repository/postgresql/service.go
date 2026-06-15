package postgresql

import (
	"context"
	"fmt"
	"testJob/internal/models"

	"github.com/google/uuid"
)

func (s *Storage) SaveService(ctx context.Context, service *models.Service) error {
	const operation = "storage.postgresql.SaveService"
	req := `INSERT INTO services (service_id, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.ExecContext(
		ctx,
		req,
		service.ServiceID,
		service.Price,
		service.UserID,
		service.StartDate,
		service.EndDate,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return err
}

func (s *Storage) GetServicesByUserID(ctx context.Context, userID uuid.UUID) (*[]models.Service, error) {
	const operation = "repository.postgresql.GetServicesByUserID"
	var result []models.Service
	req := `SELECT * FROM services WHERE user_id = $1`
	err := s.db.SelectContext(ctx, &result, req, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return &result, err
}

func (s *Storage) GetServiceByID(ctx context.Context, serviceID uint) (*models.Service, error) {
	const operation = "repository.postgresql.GetServiceByID"
	var result *models.Service
	req := `SELECT * FROM services WHERE id = $1`
	err := s.db.GetContext(ctx, &result, req, serviceID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return result, err
}

func (s *Storage) DeleteService(ctx context.Context, serviceID uint) error {
	const operation = "repository.postgresql.DeleteService"
	req := `DELETE FROM services WHERE id = $1`
	_, err := s.db.ExecContext(ctx, req, serviceID)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return err
}

func (s *Storage) ListServices(ctx context.Context) (*[]models.Service, error) {
	const operation = "repository.postgresql.ListServices"
	var result *[]models.Service
	req := `SELECT * FROM services`
	err := s.db.SelectContext(ctx, &result, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return result, err
}
