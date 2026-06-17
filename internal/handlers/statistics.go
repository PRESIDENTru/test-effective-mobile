package handlers

import (
	"net/http"

	"github.com/google/uuid"
)

// GetTotalCost рассчитывает общую стоимость
// @Summary Расчет общей стоимости подписок
// @Router /subscriptions/cost [get]
// @Produce json
// @Param from query string true "Начало периода (MM-YYYY)"
// @Param to query string true "Конец периода (MM-YYYY)"
// @Param user_id query string false "Фильтр по ID пользователя"
// @Param service_name query string false "Фильтр по названию сервиса"
// @Success 200
// @Failure 400
// @Failure 500
func (h *Handler) GetTotalCost(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetTotalCost"
	q := r.URL.Query()

	fromStr := q.Get("from")
	toStr := q.Get("to")
	if fromStr == "" || toStr == "" {
		h.logger.Error(op, "missing from/to query parameter")
		writeError(w, http.StatusBadRequest, "from and to are required (MM-YYYY)")
		return
	}

	from, err := parseMonthYear(fromStr)
	if err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusBadRequest, "invalid from: expected MM-YYYY")
		return
	}
	to, err := parseMonthYear(toStr)
	if err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusBadRequest, "invalid to: expected MM-YYYY")
		return
	}
	if to.Before(from) {
		writeError(w, http.StatusBadRequest, "to must be after from")
		return
	}

	var userID *uuid.UUID
	if raw := q.Get("user_id"); raw != "" {
		uid, err := uuid.Parse(raw)
		if err != nil {
			h.logger.Error(op, err)
			writeError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		userID = &uid
	}

	var serviceName *string
	if raw := q.Get("service_name"); raw != "" {
		serviceName = &raw
	}

	total, err := h.storage.GetTotalCost(r.Context(), from, to, userID, serviceName)
	if err != nil {
		h.logger.Error(op, err)
		writeError(w, http.StatusInternalServerError, "failed to calculate total cost")
		return
	}

	h.logger.Info(op, "total cost calculated")
	writeJSON(w, http.StatusOK, map[string]uint{"total": total})
}
