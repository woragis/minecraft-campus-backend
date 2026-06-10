package httpserver

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	citysvc "github.com/woragis/minecraft-campus-backend/server/internal/city/service"
)

type cityHandler struct {
	svc *citysvc.Service
}

func newCityHandler(svc *citysvc.Service) *cityHandler {
	return &cityHandler{svc: svc}
}

func (h *cityHandler) list(w http.ResponseWriter, r *http.Request) {
	cities, err := h.svc.List(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"cities": cities})
}

func (h *cityHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeCityGetV1HandlerIDInvalid, apperrors.MsgCityGetV1HandlerIDInvalid))
		return
	}
	city, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, city)
}

func (h *cityHandler) getBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	city, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, city)
}
