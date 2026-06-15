package postgresql

import (
	"context"
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
