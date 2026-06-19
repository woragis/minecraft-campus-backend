package httpserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	"github.com/woragis/minecraft-campus-backend/server/internal/webauth"
)

type webHandler struct {
	webAuth *webauth.Service
}

func newWebHandler(webAuth *webauth.Service) *webHandler {
	return &webHandler{webAuth: webAuth}
}

type webSessionBody struct {
	Code string `json:"code"`
}

func (h *webHandler) createSession(w http.ResponseWriter, r *http.Request) {
	var body webSessionBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeWebSessionV1HandlerBodyInvalid, apperrors.MsgWebSessionV1HandlerBodyInvalid, err))
		return
	}
	session, err := h.webAuth.RedeemLinkCode(r.Context(), body.Code)
	if err != nil {
		apperrors.WriteError(w, apperrors.Unauthorized(apperrors.CodeWebSessionV1HandlerCodeInvalid, apperrors.MsgWebSessionV1HandlerCodeInvalid))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"token":         session.Token,
		"expiresAt":     session.ExpiresAt,
		"playerId":      session.PlayerID,
		"username":      session.Username,
		"minecraftUuid": session.MinecraftUUID,
	})
}

func bearerToken(r *http.Request) string {
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return strings.TrimSpace(auth[7:])
	}
	return strings.TrimSpace(r.Header.Get("X-Campus-Session"))
}

func (h *webHandler) requireSession(next func(http.ResponseWriter, *http.Request, *webauth.Session)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r)
		if token == "" {
			apperrors.WriteError(w, apperrors.Unauthorized(apperrors.CodeWebSessionV1HandlerMissing, apperrors.MsgWebSessionV1HandlerMissing))
			return
		}
		session, err := h.webAuth.ResolveToken(r.Context(), token)
		if err != nil {
			apperrors.WriteError(w, apperrors.Unauthorized(apperrors.CodeWebSessionV1HandlerInvalid, apperrors.MsgWebSessionV1HandlerInvalid))
			return
		}
		next(w, r, session)
	}
}

func (h *webHandler) logout(w http.ResponseWriter, r *http.Request) {
	token := bearerToken(r)
	if token != "" {
		h.webAuth.RevokeToken(r.Context(), token)
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "logged_out"})
}
