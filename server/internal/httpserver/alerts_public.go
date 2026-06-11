package httpserver

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	alertssvc "github.com/woragis/minecraft-campus-backend/server/internal/alerts/service"
)

type alertsHandler struct {
	svc *alertssvc.Service
}

func newAlertsHandler(svc *alertssvc.Service) *alertsHandler {
	return &alertsHandler{svc: svc}
}

func (h *alertsHandler) list(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.svc.ListUnacknowledged(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"alerts": alerts})
}

func (h *alertsHandler) acknowledge(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeAlertsAckV1HandlerIDInvalid, apperrors.MsgAlertsAckV1HandlerIDInvalid))
		return
	}
	if err := h.svc.Acknowledge(r.Context(), id); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "acknowledged"})
}
