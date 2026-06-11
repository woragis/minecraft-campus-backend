package httpserver

import (
	"net/http"

	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	metricssvc "github.com/woragis/minecraft-campus-backend/server/internal/metrics/service"
)

type metricsHandler struct {
	svc *metricssvc.Service
}

func newMetricsHandler(svc *metricssvc.Service) *metricsHandler {
	return &metricsHandler{svc: svc}
}

func (h *metricsHandler) overview(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.Overview(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *metricsHandler) territory(w http.ResponseWriter, r *http.Request) {
	rows, err := h.svc.TerritoryByServer(r.Context())
	if err != nil {
		apperrors.WriteError(w, apperrors.InternalCause(apperrors.CodeMetricsTerritoryV1ServiceFailed, apperrors.MsgMetricsTerritoryV1ServiceFailed, err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"servers": rows})
}
