package httpserver

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	presencesvc "github.com/woragis/minecraft-campus-backend/server/internal/presence"
)

type presenceHandler struct {
	svc *presencesvc.Service
}

func newPresenceHandler(svc *presencesvc.Service) *presenceHandler {
	return &presenceHandler{svc: svc}
}

func (h *presenceHandler) overview(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.Overview(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *presenceHandler) server(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.Server(r.Context(), r.PathValue("slug"))
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *presenceHandler) guild(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeGuildGetV1HandlerIDInvalid, apperrors.MsgGuildGetV1HandlerIDInvalid))
		return
	}
	out, err := h.svc.Guild(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
