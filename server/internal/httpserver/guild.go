package httpserver

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	guildsvc "github.com/woragis/minecraft-campus-backend/server/internal/guild/service"
)

type guildHandler struct {
	svc *guildsvc.Service
}

func newGuildHandler(svc *guildsvc.Service) *guildHandler {
	return &guildHandler{svc: svc}
}

func (h *guildHandler) list(w http.ResponseWriter, r *http.Request) {
	guilds, err := h.svc.List(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"guilds": guilds})
}

func (h *guildHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeGuildGetV1HandlerIDInvalid, apperrors.MsgGuildGetV1HandlerIDInvalid))
		return
	}
	guild, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, guild)
}

func (h *guildHandler) getBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	guild, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, guild)
}

func (h *guildHandler) listMembers(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeGuildGetV1HandlerIDInvalid, apperrors.MsgGuildGetV1HandlerIDInvalid))
		return
	}
	members, err := h.svc.ListMembers(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"members": members})
}
