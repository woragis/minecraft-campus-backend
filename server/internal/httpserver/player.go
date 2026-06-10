package httpserver

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	guildsvc "github.com/woragis/minecraft-campus-backend/server/internal/guild/service"
	playersvc "github.com/woragis/minecraft-campus-backend/server/internal/player/service"
)

type playerHandler struct {
	svc      *playersvc.Service
	guildSvc *guildsvc.Service
}

type playerProfile struct {
	*models.Player
	Guild *models.Guild `json:"guild,omitempty"`
}

func newPlayerHandler(svc *playersvc.Service, guildSvc *guildsvc.Service) *playerHandler {
	return &playerHandler{svc: svc, guildSvc: guildSvc}
}

func (h *playerHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodePlayerGetV1HandlerIDInvalid, apperrors.MsgPlayerGetV1HandlerIDInvalid))
		return
	}
	player, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, h.profile(r, player))
}

func (h *playerHandler) getByMinecraftUUID(w http.ResponseWriter, r *http.Request) {
	mcUUID, err := uuid.Parse(r.PathValue("minecraftUuid"))
	if err != nil || mcUUID == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodePlayerGetMinecraftV1HandlerUUIDInvalid, apperrors.MsgPlayerGetMinecraftV1HandlerUUIDInvalid))
		return
	}
	player, err := h.svc.GetByMinecraftUUID(r.Context(), mcUUID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, h.profile(r, player))
}

func (h *playerHandler) listInvitees(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodePlayerInvitesListV1HandlerIDInvalid, apperrors.MsgPlayerInvitesListV1HandlerIDInvalid))
		return
	}
	players, err := h.svc.ListDirectInvitees(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"invitees": players})
}

func (h *playerHandler) profile(r *http.Request, player *models.Player) playerProfile {
	out := playerProfile{Player: player}
	if h.guildSvc != nil {
		guild, err := h.guildSvc.PlayerGuild(r.Context(), player.ID)
		if err == nil && guild != nil {
			out.Guild = guild
		}
	}
	return out
}
