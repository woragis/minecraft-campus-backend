package httpserver

import (
	"net/http"
	"strings"

	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	invitesvc "github.com/woragis/minecraft-campus-backend/server/internal/invite/service"
	playersvc "github.com/woragis/minecraft-campus-backend/server/internal/player/service"
)

type inviteHandler struct {
	invites *invitesvc.Service
	players *playersvc.Service
}

func newInviteHandler(invites *invitesvc.Service, players *playersvc.Service) *inviteHandler {
	return &inviteHandler{invites: invites, players: players}
}

func (h *inviteHandler) getByCode(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.PathValue("code"))
	if code == "" {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInviteGetV1HandlerCodeEmpty, apperrors.MsgInviteGetV1HandlerCodeEmpty))
		return
	}
	inv, err := h.invites.GetByCode(r.Context(), code)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, inv)
}
