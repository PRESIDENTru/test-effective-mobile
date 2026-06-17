package postgresql

import (
	"context"
	"fmt"
	"testJob/internal/models/db"

	"github.com/google/uuid"
)

func (s *Storage) SaveService(ctx context.Context, service *db.Service) error {
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

func (s *Storage) GetServicesByUserID(ctx context.Context, userID uuid.UUID) (*[]db.Service, error) {
	const operation = "repository.postgresql.GetServicesByUserID"
	var result []db.Service
	req := `SELECT * FROM services WHERE user_id = $1`
	err := s.db.SelectContext(ctx, &result, req, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return &result, err
}

func (s *Storage) GetServiceByID(ctx context.Context, serviceID uint) (*db.Service, error) {
	const operation = "repository.postgresql.GetServiceByID"
	var result db.Service
	req := `SELECT * FROM services WHERE id = $1`
	err := s.db.GetContext(ctx, &result, req, serviceID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return &result, err
}

func (s *Storage) GetServiceNameByName(ctx context.Context, name string) (*db.ServiceName, error) {
	const op = "repository.postgresql.GetServiceNameByName"
	var result db.ServiceName
	err := s.db.GetContext(ctx, &result, `SELECT * FROM services_name WHERE name=$1`, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &result, nil
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

func (s *Storage) ListServices(ctx context.Context) (*[]db.Service, error) {
	const operation = "repository.postgresql.ListServices"
	var result []db.Service
	req := `SELECT * FROM services`
	err := s.db.SelectContext(ctx, &result, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return &result, err
}

func (s *Storage) UpdateService(ctx context.Context, service *db.Service) error {
	const op = "repository.postgresql.UpdateService"
	req := `UPDATE services SET service_id=$1, price=$2, user_id=$3, start_date=$4, end_date=$5 WHERE id=$6`
	_, err := s.db.ExecContext(ctx, req,
		service.ServiceID, service.Price, service.UserID,
		service.StartDate, service.EndDate, service.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
