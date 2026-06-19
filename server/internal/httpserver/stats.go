package httpserver

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	statssvc "github.com/woragis/minecraft-campus-backend/server/internal/stats"
)

type statsHandler struct {
	svc *statssvc.Service
}

func newStatsHandler(svc *statssvc.Service) *statsHandler {
	return &statsHandler{svc: svc}
}

func (h *statsHandler) playerStats(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodePlayerGetV1HandlerIDInvalid, apperrors.MsgPlayerGetV1HandlerIDInvalid))
		return
	}
	out, err := h.svc.GetPlayerStats(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
