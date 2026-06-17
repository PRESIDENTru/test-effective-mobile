package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	_ "testJob/docs"
	"testJob/internal/models/db"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Storage interface {
	SaveService(ctx context.Context, service *db.Service) error
	GetServiceByID(ctx context.Context, serviceID uint) (*db.Service, error)
	GetServicesByUserID(ctx context.Context, userID uuid.UUID) (*[]db.Service, error)
	DeleteService(ctx context.Context, serviceID uint) error
	ListServices(ctx context.Context) (*[]db.Service, error)
	UpdateService(ctx context.Context, service *db.Service) error
	GetOrCreateServiceName(ctx context.Context, name string) (uint, error)
	GetTotalCost(ctx context.Context, from, to time.Time, userID *uuid.UUID, serviceName *string) (uint, error)
}

type Handler struct {
	storage Storage
	logger  *slog.Logger
}

func NewHandler(storage Storage, logger *slog.Logger) *Handler {
	return &Handler{
		storage: storage,
		logger:  logger,
	}
}

// Router возвращает маршрутизатор API
// @Router / [get]
func (h *Handler) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/subscriptions", h.CreateSubscription)
	r.Get("/subscriptions", h.ListSubscriptions)
	r.Get("/subscriptions/{id}", h.GetSubscription)
	r.Put("/subscriptions/{id}", h.UpdateSubscription)
	r.Delete("/subscriptions/{id}", h.DeleteSubscription)
	r.Get("/subscriptions/cost", h.GetTotalCost)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	return r
}

func parseMonthYear(s string) (time.Time, error) {
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return time.Time{}, errors.New("expected MM-YYYY")
	}
	return t, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func parseUintParam(r *http.Request, key string) (uint, error) {
	val := chi.URLParam(r, key)
	n, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(n), nil
}
