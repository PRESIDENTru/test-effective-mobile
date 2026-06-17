package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testJob/internal/models/db"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSubscription_OK(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	s.On("GetOrCreateServiceName", mock.Anything, "Yandex Plus").Return(uint(1), nil)
	s.On("SaveService", mock.Anything, mock.MatchedBy(func(svc *db.Service) bool {
		return svc.Price == 400 && svc.ServiceID == 1
	})).Return(nil)

	body := mustJSON(map[string]any{
		"service_name": "Yandex Plus",
		"price":        400,
		"user_id":      uuid.New().String(),
		"start_date":   "07-2025",
	})
	req := httptest.NewRequest(http.MethodPost, "/subscriptions", body)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	s.AssertExpectations(t)
}

func TestCreateSubscription_InvalidJSON(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	req := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewBufferString("not-json"))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateSubscription_InvalidUserID(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	body := mustJSON(map[string]any{
		"service_name": "Netflix",
		"price":        300,
		"user_id":      "not-a-uuid",
		"start_date":   "01-2025",
	})
	req := httptest.NewRequest(http.MethodPost, "/subscriptions", body)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateSubscription_InvalidStartDate(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	body := mustJSON(map[string]any{
		"service_name": "Netflix",
		"price":        300,
		"user_id":      uuid.New().String(),
		"start_date":   "2025-01-01",
	})
	req := httptest.NewRequest(http.MethodPost, "/subscriptions", body)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateSubscription_WithEndDate(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	s.On("GetOrCreateServiceName", mock.Anything, "Spotify").Return(uint(2), nil)
	s.On("SaveService", mock.Anything, mock.MatchedBy(func(svc *db.Service) bool {
		return svc.EndDate != nil
	})).Return(nil)

	body := mustJSON(map[string]any{
		"service_name": "Spotify",
		"price":        199,
		"user_id":      uuid.New().String(),
		"start_date":   "01-2025",
		"end_date":     "12-2025",
	})
	req := httptest.NewRequest(http.MethodPost, "/subscriptions", body)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	s.AssertExpectations(t)
}

func TestGetSubscription_OK(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	svc := &db.Service{ID: 1, Price: 400, ServiceID: 1, UserID: uuid.New(), StartDate: time.Now()}
	s.On("GetServiceByID", mock.Anything, uint(1)).Return(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var got db.Service
	json.NewDecoder(rec.Body).Decode(&got)
	assert.Equal(t, uint(1), got.ID)
}

func TestGetSubscription_NotFound(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	s.On("GetServiceByID", mock.Anything, uint(99)).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/99", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteSubscription_OK(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	s.On("DeleteService", mock.Anything, uint(5)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/subscriptions/5", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	s.AssertExpectations(t)
}

func TestDeleteSubscription_StorageError(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	s.On("DeleteService", mock.Anything, uint(5)).Return(assert.AnError)

	req := httptest.NewRequest(http.MethodDelete, "/subscriptions/5", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUpdateSubscription_OK(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	existing := &db.Service{ID: 3, Price: 200, ServiceID: 1, UserID: uuid.New(), StartDate: time.Now()}
	s.On("GetServiceByID", mock.Anything, uint(3)).Return(existing, nil)
	s.On("GetOrCreateServiceName", mock.Anything, "Netflix").Return(uint(2), nil)
	s.On("UpdateService", mock.Anything, mock.MatchedBy(func(svc *db.Service) bool {
		return svc.Price == 500 && svc.ServiceID == 2
	})).Return(nil)

	body := mustJSON(map[string]any{
		"service_name": "Netflix",
		"price":        500,
	})
	req := httptest.NewRequest(http.MethodPut, "/subscriptions/3", body)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	s.AssertExpectations(t)
}

func TestUpdateSubscription_NotFound(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	s.On("GetServiceByID", mock.Anything, uint(99)).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodPut, "/subscriptions/99", mustJSON(map[string]any{"price": 100}))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestListSubscriptions_All(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	svcs := &[]db.Service{
		{ID: 1, Price: 100},
		{ID: 2, Price: 200},
	}
	s.On("ListServices", mock.Anything).Return(svcs, nil)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var got []db.Service
	json.NewDecoder(rec.Body).Decode(&got)
	assert.Len(t, got, 2)
}

func TestListSubscriptions_ByUserID(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	uid := uuid.New()
	svcs := &[]db.Service{{ID: 1, UserID: uid, Price: 300}}
	s.On("GetServicesByUserID", mock.Anything, uid).Return(svcs, nil)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions?user_id="+uid.String(), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	s.AssertExpectations(t)
}

func TestListSubscriptions_InvalidUserID(t *testing.T) {
	s := new(MockStorage)
	_, router := newHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions?user_id=bad-uuid", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
