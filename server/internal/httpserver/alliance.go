package httpserver

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	alliancesvc "github.com/woragis/minecraft-campus-backend/server/internal/alliance/service"
)

type allianceHandler struct {
	svc *alliancesvc.Service
}

func newAllianceHandler(svc *alliancesvc.Service) *allianceHandler {
	return &allianceHandler{svc: svc}
}

func (h *allianceHandler) listByGuild(w http.ResponseWriter, r *http.Request) {
	guildID, err := uuid.Parse(r.PathValue("id"))
	if err != nil || guildID == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeAllianceListV1HandlerIDInvalid, apperrors.MsgAllianceListV1HandlerIDInvalid))
		return
	}
	alliances, err := h.svc.ListByGuild(r.Context(), guildID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"alliances": alliances})
}
