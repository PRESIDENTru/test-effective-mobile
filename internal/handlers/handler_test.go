package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testJob/internal/models/db"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) SaveService(ctx context.Context, service *db.Service) error {
	return m.Called(ctx, service).Error(0)
}
func (m *MockStorage) GetServiceByID(ctx context.Context, id uint) (*db.Service, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.Service), args.Error(1)
}
func (m *MockStorage) GetServicesByUserID(ctx context.Context, userID uuid.UUID) (*[]db.Service, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]db.Service), args.Error(1)
}
func (m *MockStorage) DeleteService(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockStorage) ListServices(ctx context.Context) (*[]db.Service, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]db.Service), args.Error(1)
}
func (m *MockStorage) UpdateService(ctx context.Context, service *db.Service) error {
	return m.Called(ctx, service).Error(0)
}
func (m *MockStorage) GetOrCreateServiceName(ctx context.Context, name string) (uint, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(uint), args.Error(1)
}
func (m *MockStorage) SaveServiceName(ctx context.Context, name string) error {
	return m.Called(ctx, name).Error(0)
}
func (m *MockStorage) GetServiceNameByName(ctx context.Context, name string) (*db.ServiceName, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.ServiceName), args.Error(1)
}

func newHandler(s *MockStorage) (*Handler, *chi.Mux) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	h := NewHandler(s, logger)
	return h, h.Router()
}

func mustJSON(v any) *bytes.Buffer {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}
