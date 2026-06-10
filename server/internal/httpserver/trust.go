package httpserver

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	trustsvc "github.com/woragis/minecraft-campus-backend/server/internal/trust/service"
)

type trustHandler struct {
	svc *trustsvc.Service
}

func newTrustHandler(svc *trustsvc.Service) *trustHandler {
	return &trustHandler{svc: svc}
}

func (h *trustHandler) listEvents(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodePlayerGetV1HandlerIDInvalid, apperrors.MsgPlayerGetV1HandlerIDInvalid))
		return
	}
	events, err := h.svc.ListEvents(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"events": events})
}

func (h *trustHandler) sponsorTree(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodePlayerGetV1HandlerIDInvalid, apperrors.MsgPlayerGetV1HandlerIDInvalid))
		return
	}
	invitees, err := h.svc.SponsorTree(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"invitees": invitees})
}
