package main

import (
	"log/slog"
	"net/http"
	"os"
	"testJob/internal/handlers"
	"testJob/internal/repository/postgresql"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

// @title Subscription Service API
// @version 1.0
// @description API для управления подписками пользователей
// @host localhost:8080
// @BasePath /
func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		logger.Error("DATABASE_URL is not set")
		os.Exit(1)
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("connected to database")

	storage := postgresql.NewStorage(db)
	handler := handlers.NewHandler(storage, logger)

	r := chi.NewRouter()
	r.Mount("/", handler.Router())

	addr := ":" + getEnv("APP_PORT", "8080")
	logger.Info("starting server", "addr", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Error("server error", "err", err)
		os.Exit(1)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
