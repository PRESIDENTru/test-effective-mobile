package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (s *Storage) SaveServiceName(ctx context.Context, name string) error {
	const op = "repository.postgresql.SaveServiceName"
	req := `INSERT INTO services_name (name) VALUES ($1)`
	_, err := s.db.ExecContext(ctx, req, name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return err
}

func (s *Storage) GetOrCreateServiceName(ctx context.Context, name string) (uint, error) {
	const op = "repository.postgresql.GetOrCreateServiceName"

	sn, err := s.GetServiceNameByName(ctx, name)
	if err == nil {
		return sn.ID, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id uint
	err = s.db.QueryRowContext(ctx,
		`INSERT INTO services_name (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name RETURNING id`,
		name,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
