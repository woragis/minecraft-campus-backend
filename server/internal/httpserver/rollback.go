package httpserver

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	rollbacksvc "github.com/woragis/minecraft-campus-backend/server/internal/rollback/service"
)

type rollbackHandler struct {
	svc *rollbacksvc.Service
}

func newRollbackHandler(svc *rollbacksvc.Service) *rollbackHandler {
	return &rollbackHandler{svc: svc}
}

func (h *rollbackHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeRollbackGetV1HandlerIDInvalid, apperrors.MsgRollbackGetV1HandlerIDInvalid))
		return
	}
	rb, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, rb)
}
