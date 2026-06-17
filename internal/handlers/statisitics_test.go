package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (m *MockStorage) GetTotalCost(ctx context.Context, from, to time.Time, userID *uuid.UUID, serviceName *string) (uint, error) {
	args := m.Called(ctx, from, to, userID, serviceName)
	return args.Get(0).(uint), args.Error(1)
}

func mustParseMonth(s string) time.Time {
	t, _ := time.Parse("01-2006", s)
	return t
}

func TestGetTotalCost_OK(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	from := mustParseMonth("01-2025")
	to := mustParseMonth("12-2025")
	s.On("GetTotalCost", mock.Anything, from, to, (*uuid.UUID)(nil), (*string)(nil)).
		Return(uint(1200), nil)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/cost?from=01-2025&to=12-2025", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]uint
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, uint(1200), resp["total"])
	s.AssertExpectations(t)
}

func TestGetTotalCost_WithUserID(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	uid := uuid.New()
	from := mustParseMonth("03-2025")
	to := mustParseMonth("06-2025")
	s.On("GetTotalCost", mock.Anything, from, to, &uid, (*string)(nil)).
		Return(uint(800), nil)

	req := httptest.NewRequest(http.MethodGet,
		"/subscriptions/cost?from=03-2025&to=06-2025&user_id="+uid.String(), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]uint
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, uint(800), resp["total"])
	s.AssertExpectations(t)
}

func TestGetTotalCost_WithServiceName(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	name := "Yandex Plus"
	from := mustParseMonth("01-2025")
	to := mustParseMonth("06-2025")
	s.On("GetTotalCost", mock.Anything, from, to, (*uuid.UUID)(nil), &name).
		Return(uint(2400), nil)

	req := httptest.NewRequest(http.MethodGet,
		"/subscriptions/cost?from=01-2025&to=06-2025&service_name=Yandex+Plus", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]uint
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, uint(2400), resp["total"])
}

func TestGetTotalCost_WithAllFilters(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	uid := uuid.New()
	name := "Netflix"
	from := mustParseMonth("01-2025")
	to := mustParseMonth("12-2025")
	s.On("GetTotalCost", mock.Anything, from, to, &uid, &name).
		Return(uint(3600), nil)

	req := httptest.NewRequest(http.MethodGet,
		"/subscriptions/cost?from=01-2025&to=12-2025&user_id="+uid.String()+"&service_name=Netflix", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]uint
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, uint(3600), resp["total"])
}

func TestGetTotalCost_MissingFrom(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/cost?to=12-2025", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetTotalCost_MissingTo(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/cost?from=01-2025", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetTotalCost_InvalidFrom(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/cost?from=2025-01&to=12-2025", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetTotalCost_ToBeforeFrom(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/cost?from=12-2025&to=01-2025", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetTotalCost_InvalidUserID(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/cost?from=01-2025&to=12-2025&user_id=bad", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetTotalCost_StorageError(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	from := mustParseMonth("01-2025")
	to := mustParseMonth("12-2025")
	s.On("GetTotalCost", mock.Anything, from, to, (*uuid.UUID)(nil), (*string)(nil)).
		Return(uint(0), assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/cost?from=01-2025&to=12-2025", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
