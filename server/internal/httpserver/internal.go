package httpserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	guildsvc "github.com/woragis/minecraft-campus-backend/server/internal/guild/service"
	invitesvc "github.com/woragis/minecraft-campus-backend/server/internal/invite/service"
	playersvc "github.com/woragis/minecraft-campus-backend/server/internal/player/service"
	trustsvc "github.com/woragis/minecraft-campus-backend/server/internal/trust/service"
)

type internalHandler struct {
	pluginAPIKey string
	players      *playersvc.Service
	invites      *invitesvc.Service
	guilds       *guildsvc.Service
	trust        *trustsvc.Service
}

func newInternalHandler(pluginAPIKey string, players *playersvc.Service, invites *invitesvc.Service, guilds *guildsvc.Service, trust *trustsvc.Service) *internalHandler {
	return &internalHandler{
		pluginAPIKey: pluginAPIKey,
		players:      players,
		invites:      invites,
		guilds:       guilds,
		trust:        trust,
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

type createGuildBody struct {
	LeaderUUID string `json:"leaderUuid"`
	Name       string `json:"name"`
}

func (h *internalHandler) createGuild(w http.ResponseWriter, r *http.Request) {
	var body createGuildBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeGuildPostV1HandlerBodyInvalid, apperrors.MsgGuildPostV1HandlerBodyInvalid, err))
		return
	}
	leaderUUID, err := uuid.Parse(strings.TrimSpace(body.LeaderUUID))
	if err != nil || leaderUUID == uuid.Nil || strings.TrimSpace(body.Name) == "" {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeGuildPostV1HandlerBodyInvalid, apperrors.MsgGuildPostV1HandlerBodyInvalid))
		return
	}
	guild, err := h.guilds.Create(r.Context(), leaderUUID, body.Name)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, guild)
}

type guildPlayerBody struct {
	PlayerUUID string `json:"playerUuid"`
}

func (h *internalHandler) joinGuild(w http.ResponseWriter, r *http.Request) {
	guildID, err := uuid.Parse(r.PathValue("id"))
	if err != nil || guildID == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeGuildGetV1HandlerIDInvalid, apperrors.MsgGuildGetV1HandlerIDInvalid))
		return
	}
	var body guildPlayerBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeGuildPostV1HandlerBodyInvalid, apperrors.MsgGuildPostV1HandlerBodyInvalid, err))
		return
	}
	playerUUID, err := uuid.Parse(strings.TrimSpace(body.PlayerUUID))
	if err != nil || playerUUID == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeGuildPostV1HandlerBodyInvalid, apperrors.MsgGuildPostV1HandlerBodyInvalid))
		return
	}
	if err := h.guilds.Join(r.Context(), guildID, playerUUID); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "joined"})
}

func (h *internalHandler) leaveGuild(w http.ResponseWriter, r *http.Request) {
	guildID, err := uuid.Parse(r.PathValue("id"))
	if err != nil || guildID == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeGuildGetV1HandlerIDInvalid, apperrors.MsgGuildGetV1HandlerIDInvalid))
		return
	}
	var body guildPlayerBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeGuildPostV1HandlerBodyInvalid, apperrors.MsgGuildPostV1HandlerBodyInvalid, err))
		return
	}
	playerUUID, err := uuid.Parse(strings.TrimSpace(body.PlayerUUID))
	if err != nil || playerUUID == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeGuildPostV1HandlerBodyInvalid, apperrors.MsgGuildPostV1HandlerBodyInvalid))
		return
	}
	if err := h.guilds.Leave(r.Context(), guildID, playerUUID); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "left"})
}

type trustEventBody struct {
	PlayerID  string `json:"playerId"`
	EventType string `json:"eventType"`
	Reason    string `json:"reason"`
}

func (h *internalHandler) recordTrustEvent(w http.ResponseWriter, r *http.Request) {
	var body trustEventBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeTrustEventV1HandlerBodyInvalid, apperrors.MsgTrustEventV1HandlerBodyInvalid, err))
		return
	}
	playerID, err := uuid.Parse(strings.TrimSpace(body.PlayerID))
	if err != nil || playerID == uuid.Nil || strings.TrimSpace(body.EventType) == "" {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeTrustEventV1HandlerBodyInvalid, apperrors.MsgTrustEventV1HandlerBodyInvalid))
		return
	}
	event, player, err := h.trust.RecordEvent(r.Context(), playerID, body.EventType, body.Reason, nil)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"event": event, "player": player})
}
