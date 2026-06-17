package db

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `db:"id"`
}

type ServiceName struct {
	ID   uint   `db:"id"`
	Name string `db:"name"`
}

type Service struct {
	ID        uint       `db:"id"`
	ServiceID uint       `db:"service_id"`
	Price     uint       `db:"price"`
	UserID    uuid.UUID  `db:"user_id"`
	StartDate time.Time  `db:"start_date"`
	EndDate   *time.Time `db:"end_date"`
}
