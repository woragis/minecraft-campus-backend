package httpserver

import (
	"net/http"
	"strings"

	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
)

const headerPluginKey = "X-Plugin-Key"

func (h *internalHandler) requirePluginKey(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimSpace(r.Header.Get(headerPluginKey))
		if key == "" {
			apperrors.WriteError(w, apperrors.Unauthorized(apperrors.CodePluginAuthV1HandlerKeyMissing, apperrors.MsgPluginAuthV1HandlerKeyMissing))
			return
		}
		if h.pluginAPIKey == "" || key != h.pluginAPIKey {
			apperrors.WriteError(w, apperrors.Unauthorized(apperrors.CodePluginAuthV1HandlerKeyInvalid, apperrors.MsgPluginAuthV1HandlerKeyInvalid))
			return
		}
		next(w, r)
	}
}
