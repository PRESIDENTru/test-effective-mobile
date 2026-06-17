package handlers

import (
	"encoding/json"
	"net/http"
	"testJob/internal/models/db"
	"time"

	"github.com/google/uuid"
)

// CreateSubscription создает новую подписку
// @Summary Создание новой подписки
// @Router /subscriptions [post]
// @Accept json
// @Produce json
// @Success 201
// @Failure 400
// @Failure 500
func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CreateSubscription"
	var req struct {
		ServiceName string `json:"service_name"`
		Price       uint   `json:"price"`
		UserID      string `json:"user_id"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var endDate *time.Time
	if req.EndDate != "" {
		ed, err := parseMonthYear(req.EndDate)
		if err != nil {
			h.logger.Error(op, err)
			writeError(w, http.StatusBadRequest, "invalid end_date, expected MM-YYYY")
			return
		}
		endDate = &ed
	}
	serviceNameID, err := h.storage.GetOrCreateServiceName(r.Context(), req.ServiceName)
	if err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusInternalServerError, "failed to resolve service name")
		return
	}

	svc := &db.Service{
		ServiceID: serviceNameID,
		Price:     req.Price,
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}
	if err := h.storage.SaveService(r.Context(), svc); err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusInternalServerError, "failed to save subscription")
		return
	}

	w.WriteHeader(http.StatusCreated)
	h.logger.Info(op, "subscription created")
}

// GetSubscription возвращает подписку по ID
// @Summary Получение подписки по ID
// @Router /subscriptions/{id} [get]
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200
// @Failure 400
// @Failure 404
func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetSubscription"
	id, err := parseUintParam(r, "id")
	if err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	svc, err := h.storage.GetServiceByID(r.Context(), id)
	if err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusNotFound, "subscription not found")
		return
	}

	writeJSON(w, http.StatusOK, svc)
	h.logger.Info(op, "subscription found")
}

// UpdateSubscription обновляет подписку
// @Summary Обновление подписки
// @Router /subscriptions/{id} [put]
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	const op = "handler.UpdateSubscription"
	id, err := parseUintParam(r, "id")
	if err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	svc, err := h.storage.GetServiceByID(r.Context(), id)
	if err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusNotFound, "subscription not found")
		return
	}

	var req struct {
		ServiceName string `json:"service_name"`
		Price       uint   `json:"price"`
		UserID      string `json:"user_id"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.ServiceName != "" {
		nameID, err := h.storage.GetOrCreateServiceName(r.Context(), req.ServiceName)
		if err != nil {
			h.logger.Error(op, err)
			writeError(w, http.StatusInternalServerError, "failed to resolve service name")
			return
		}
		svc.ServiceID = nameID
	}
	if req.Price != 0 {
		svc.Price = req.Price
	}
	if req.UserID != "" {
		uid, err := uuid.Parse(req.UserID)
		if err != nil {
			h.logger.Error(op, err)
			writeError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		svc.UserID = uid
	}
	if req.StartDate != "" {
		sd, err := parseMonthYear(req.StartDate)
		if err != nil {
			h.logger.Error(op, err)
			writeError(w, http.StatusBadRequest, "invalid start_date")
			return
		}
		svc.StartDate = sd
	}
	if req.EndDate != "" {
		ed, err := parseMonthYear(req.EndDate)
		if err != nil {
			h.logger.Error(op, err)
			writeError(w, http.StatusBadRequest, "invalid end_date")
			return
		}
		svc.EndDate = &ed
	}

	if err := h.storage.UpdateService(r.Context(), svc); err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusInternalServerError, "failed to update subscription")
		return
	}

	h.logger.Info(op, "subscription updated")
	writeJSON(w, http.StatusOK, svc)
}

// DeleteSubscription удаляет подписку
// @Summary Удаление подписки
// @Router /subscriptions/{id} [delete]
// @Param id path int true "ID подписки"
// @Success 204
// @Failure 400
// @Failure 500
func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeleteSubscription"
	id, err := parseUintParam(r, "id")
	if err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.storage.DeleteService(r.Context(), id); err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusInternalServerError, "failed to delete subscription")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.logger.Info(op, "subscription deleted")
}

// ListSubscriptions возвращает список подписок
// @Summary Получение списка подписок
// @Router /subscriptions [get]
// @Produce json
// @Param user_id query string false "Фильтр по ID пользователя"
// @Success 200
// @Failure 400
// @Failure 500
func (h *Handler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	const op = "handler.ListSubscriptions"
	q := r.URL.Query()
	userIDStr := q.Get("user_id")

	if userIDStr != "" {
		uid, err := uuid.Parse(userIDStr)
		if err != nil {
			h.logger.Error(op, err)
			writeError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		result, err := h.storage.GetServicesByUserID(r.Context(), uid)
		if err != nil {
			h.logger.Error(op, err)
			writeError(w, http.StatusInternalServerError, "failed to list subscriptions")
			return
		}
		writeJSON(w, http.StatusOK, result)
		h.logger.Info(op, "subscriptions found")
		return
	}

	result, err := h.storage.ListServices(r.Context())
	if err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusInternalServerError, "failed to list subscriptions")
		return
	}
	writeJSON(w, http.StatusOK, result)
	h.logger.Info(op, "subscriptions found")
}
