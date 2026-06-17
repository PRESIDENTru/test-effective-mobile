package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *Storage) GetTotalCost(ctx context.Context, from, to time.Time, userID *uuid.UUID, serviceName *string) (uint, error) {
	const op = "repository.postgresql.GetTotalCost"

	query := `
		SELECT COALESCE(SUM(sv.price), 0)
		FROM services sv
		JOIN services_name sn ON sn.id = sv.service_id
		WHERE sv.start_date <= $2
		  AND (sv.end_date IS NULL OR sv.end_date >= $1)
	`
	args := []any{from, to}
	argN := 3

	if userID != nil {
		query += fmt.Sprintf(" AND sv.user_id = $%d", argN)
		args = append(args, *userID)
		argN++
	}
	if serviceName != nil {
		query += fmt.Sprintf(" AND sn.name = $%d", argN)
		args = append(args, *serviceName)
	}

	var total uint
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return total, nil
}
