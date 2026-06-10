package httpserver

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	auditsvc "github.com/woragis/minecraft-campus-backend/server/internal/audit/service"
)

type auditHandler struct {
	svc *auditsvc.Service
}

func newAuditHandler(svc *auditsvc.Service) *auditHandler {
	return &auditHandler{svc: svc}
}

func (h *auditHandler) listByPlayer(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeAuditListV1HandlerIDInvalid, apperrors.MsgAuditListV1HandlerIDInvalid))
		return
	}
	q := r.URL.Query()
	var from, to *time.Time
	if v := q.Get("from"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			from = &t
		}
	}
	if v := q.Get("to"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			to = &t
		}
	}
	events, err := h.svc.ListByPlayer(r.Context(), id, from, to, q.Get("eventType"))
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"events": events})
}
