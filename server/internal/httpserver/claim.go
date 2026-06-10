package httpserver

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	claimsvc "github.com/woragis/minecraft-campus-backend/server/internal/claim/service"
)

type claimHandler struct {
	svc *claimsvc.Service
}

func newClaimHandler(svc *claimsvc.Service) *claimHandler {
	return &claimHandler{svc: svc}
}

func (h *claimHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeClaimGetV1HandlerIDInvalid, apperrors.MsgClaimGetV1HandlerIDInvalid))
		return
	}
	claim, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, claim)
}

func (h *claimHandler) listByCity(w http.ResponseWriter, r *http.Request) {
	cityID, err := uuid.Parse(r.PathValue("id"))
	if err != nil || cityID == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeClaimListCityV1HandlerIDInvalid, apperrors.MsgClaimListCityV1HandlerIDInvalid))
		return
	}
	claims, err := h.svc.ListByCity(r.Context(), cityID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"claims": claims})
}
