package httpserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	invitesvc "github.com/woragis/minecraft-campus-backend/server/internal/invite/service"
	playersvc "github.com/woragis/minecraft-campus-backend/server/internal/player/service"
)

type internalHandler struct {
	pluginAPIKey string
	players      *playersvc.Service
	invites      *invitesvc.Service
}

func newInternalHandler(pluginAPIKey string, players *playersvc.Service, invites *invitesvc.Service) *internalHandler {
	return &internalHandler{
		pluginAPIKey: pluginAPIKey,
		players:      players,
		invites:      invites,
	}
}

func (h *internalHandler) whitelist(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.PathValue("minecraftUuid")
	mcUUID, err := uuid.Parse(uuidStr)
	if err != nil || mcUUID == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeWhitelistGetV1HandlerUUIDInvalid, apperrors.MsgWhitelistGetV1HandlerUUIDInvalid))
		return
	}
	username := strings.TrimSpace(r.URL.Query().Get("username"))
	out, err := h.players.CheckWhitelist(r.Context(), mcUUID, username)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

type upsertPlayerBody struct {
	MinecraftUUID string `json:"minecraftUuid"`
	Username      string `json:"username"`
	ServerSlug    string `json:"serverSlug"`
}

func (h *internalHandler) upsertPlayer(w http.ResponseWriter, r *http.Request) {
	var body upsertPlayerBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodePlayerUpsertV1HandlerBodyInvalid, apperrors.MsgPlayerUpsertV1HandlerBodyInvalid, err))
		return
	}
	mcUUID, err := uuid.Parse(strings.TrimSpace(body.MinecraftUUID))
	if err != nil || mcUUID == uuid.Nil || strings.TrimSpace(body.Username) == "" || strings.TrimSpace(body.ServerSlug) == "" {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodePlayerUpsertV1HandlerBodyInvalid, apperrors.MsgPlayerUpsertV1HandlerBodyInvalid))
		return
	}
	player, err := h.players.UpsertFromPlugin(r.Context(), mcUUID, body.Username, body.ServerSlug)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, player)
}

type createInviteBody struct {
	SponsorUUID    string `json:"sponsorUuid"`
	TargetUsername string `json:"targetUsername"`
}

func (h *internalHandler) createInvite(w http.ResponseWriter, r *http.Request) {
	var body createInviteBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeInvitePostInternalV1HandlerBodyInvalid, apperrors.MsgInvitePostInternalV1HandlerBodyInvalid, err))
		return
	}
	sponsorUUID, err := uuid.Parse(strings.TrimSpace(body.SponsorUUID))
	if err != nil || sponsorUUID == uuid.Nil || strings.TrimSpace(body.TargetUsername) == "" {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInvitePostInternalV1HandlerBodyInvalid, apperrors.MsgInvitePostInternalV1HandlerBodyInvalid))
		return
	}
	inv, err := h.invites.CreateForSponsor(r.Context(), sponsorUUID, body.TargetUsername)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, inv)
}
