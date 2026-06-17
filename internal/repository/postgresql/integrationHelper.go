package postgresql

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func integrationDB(t *testing.T) *sqlx.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		dsn = "host=localhost port=5433 user=postgres password=postgres dbname=postgres sslmode=disable TimeZone=UTC"
	}
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	if err = db.Ping(); err != nil {
		t.Fatalf("db.Ping: %v (set TEST_DSN or start postgres)", err)
	}
	cleanupDB(t, db)
	t.Cleanup(func() {
		cleanupDB(t, db)
		db.Close()
	})
	return db
}

func cleanupDB(t *testing.T, db *sqlx.DB) {
	t.Helper()
	_, err := db.Exec(`TRUNCATE TABLE 
        services_name, services
        CASCADE`)
	if err != nil {
		t.Logf("cleanup: %v", err)
	}
}
