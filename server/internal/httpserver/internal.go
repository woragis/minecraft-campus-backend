package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	alliancesvc "github.com/woragis/minecraft-campus-backend/server/internal/alliance/service"
	auditsvc "github.com/woragis/minecraft-campus-backend/server/internal/audit/service"
	citysvc "github.com/woragis/minecraft-campus-backend/server/internal/city/service"
	claimsvc "github.com/woragis/minecraft-campus-backend/server/internal/claim/service"
	guildsvc "github.com/woragis/minecraft-campus-backend/server/internal/guild/service"
	invitesvc "github.com/woragis/minecraft-campus-backend/server/internal/invite/service"
	playersvc "github.com/woragis/minecraft-campus-backend/server/internal/player/service"
	rollbacksvc "github.com/woragis/minecraft-campus-backend/server/internal/rollback/service"
	trustsvc "github.com/woragis/minecraft-campus-backend/server/internal/trust/service"
)

type internalHandler struct {
	pluginAPIKey string
	players      *playersvc.Service
	invites      *invitesvc.Service
	guilds       *guildsvc.Service
	trust        *trustsvc.Service
	cities       *citysvc.Service
	claims       *claimsvc.Service
	alliances    *alliancesvc.Service
	audit        *auditsvc.Service
	rollback     *rollbacksvc.Service
}

func newInternalHandler(
	pluginAPIKey string,
	players *playersvc.Service,
	invites *invitesvc.Service,
	guilds *guildsvc.Service,
	trust *trustsvc.Service,
	cities *citysvc.Service,
	claims *claimsvc.Service,
	alliances *alliancesvc.Service,
	audit *auditsvc.Service,
	rollback *rollbacksvc.Service,
) *internalHandler {
	return &internalHandler{
		pluginAPIKey: pluginAPIKey,
		players:      players,
		invites:      invites,
		guilds:       guilds,
		trust:        trust,
		cities:       cities,
		claims:       claims,
		alliances:    alliances,
		audit:        audit,
		rollback:     rollback,
	}
}

func (h *internalHandler) bedrockWhitelist(w http.ResponseWriter, r *http.Request) {
	xuid := strings.TrimSpace(r.PathValue("xuid"))
	username := strings.TrimSpace(r.URL.Query().Get("username"))
	out, err := h.players.CheckBedrockWhitelist(r.Context(), xuid, username)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
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

type upsertBedrockPlayerBody struct {
	XUID       string `json:"xuid"`
	Username   string `json:"username"`
	ServerSlug string `json:"serverSlug"`
}

func (h *internalHandler) upsertBedrockPlayer(w http.ResponseWriter, r *http.Request) {
	var body upsertBedrockPlayerBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeBedrockPlayerUpsertV1HandlerBodyInvalid, apperrors.MsgBedrockPlayerUpsertV1HandlerBodyInvalid, err))
		return
	}
	if strings.TrimSpace(body.XUID) == "" || strings.TrimSpace(body.Username) == "" || strings.TrimSpace(body.ServerSlug) == "" {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeBedrockPlayerUpsertV1HandlerBodyInvalid, apperrors.MsgBedrockPlayerUpsertV1HandlerBodyInvalid))
		return
	}
	player, err := h.players.UpsertBedrockFromServer(r.Context(), body.XUID, body.Username, body.ServerSlug)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, player)
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

type createCityBody struct {
	FounderUUID string  `json:"founderUuid"`
	Name        string  `json:"name"`
	ServerSlug  string  `json:"serverSlug"`
	World       string  `json:"world"`
	CenterX     int     `json:"centerX"`
	CenterZ     int     `json:"centerZ"`
	GuildID     *string `json:"guildId"`
}

func (h *internalHandler) createCity(w http.ResponseWriter, r *http.Request) {
	var body createCityBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeCityPostV1HandlerBodyInvalid, apperrors.MsgCityPostV1HandlerBodyInvalid, err))
		return
	}
	founderUUID, err := uuid.Parse(strings.TrimSpace(body.FounderUUID))
	if err != nil || founderUUID == uuid.Nil || strings.TrimSpace(body.Name) == "" || strings.TrimSpace(body.ServerSlug) == "" {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeCityPostV1HandlerBodyInvalid, apperrors.MsgCityPostV1HandlerBodyInvalid))
		return
	}
	var guildID *uuid.UUID
	if body.GuildID != nil && strings.TrimSpace(*body.GuildID) != "" {
		parsed, err := uuid.Parse(strings.TrimSpace(*body.GuildID))
		if err != nil || parsed == uuid.Nil {
			apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeCityPostV1HandlerBodyInvalid, apperrors.MsgCityPostV1HandlerBodyInvalid))
			return
		}
		guildID = &parsed
	}
	city, err := h.cities.Create(r.Context(), citysvc.CreateInput{
		FounderMinecraftUUID: founderUUID,
		Name:                 body.Name,
		ServerSlug:           body.ServerSlug,
		World:                body.World,
		CenterX:              body.CenterX,
		CenterZ:              body.CenterZ,
		GuildID:              guildID,
	})
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, city)
}

type createClaimBody struct {
	OwnerUUID  string  `json:"ownerUuid"`
	ServerSlug string  `json:"serverSlug"`
	World      string  `json:"world"`
	MinX       int     `json:"minX"`
	MaxX       int     `json:"maxX"`
	MinZ       int     `json:"minZ"`
	MaxZ       int     `json:"maxZ"`
	MinY       int     `json:"minY"`
	MaxY       int     `json:"maxY"`
	ZoneType   string  `json:"zoneType"`
	CityID     *string `json:"cityId"`
}

func (h *internalHandler) createClaim(w http.ResponseWriter, r *http.Request) {
	var body createClaimBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeClaimPostV1HandlerBodyInvalid, apperrors.MsgClaimPostV1HandlerBodyInvalid, err))
		return
	}
	ownerUUID, err := uuid.Parse(strings.TrimSpace(body.OwnerUUID))
	if err != nil || ownerUUID == uuid.Nil || strings.TrimSpace(body.ServerSlug) == "" {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeClaimPostV1HandlerBodyInvalid, apperrors.MsgClaimPostV1HandlerBodyInvalid))
		return
	}
	var cityID *uuid.UUID
	if body.CityID != nil && strings.TrimSpace(*body.CityID) != "" {
		parsed, err := uuid.Parse(strings.TrimSpace(*body.CityID))
		if err != nil || parsed == uuid.Nil {
			apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeClaimPostV1HandlerBodyInvalid, apperrors.MsgClaimPostV1HandlerBodyInvalid))
			return
		}
		cityID = &parsed
	}
	claim, err := h.claims.Create(r.Context(), claimsvc.CreateInput{
		OwnerMinecraftUUID: ownerUUID,
		ServerSlug:         body.ServerSlug,
		World:              body.World,
		MinX:               body.MinX,
		MaxX:               body.MaxX,
		MinZ:               body.MinZ,
		MaxZ:               body.MaxZ,
		MinY:               body.MinY,
		MaxY:               body.MaxY,
		ZoneType:           body.ZoneType,
		CityID:             cityID,
	})
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, claim)
}

type deleteClaimBody struct {
	OwnerUUID string `json:"ownerUuid"`
}

func (h *internalHandler) deleteClaim(w http.ResponseWriter, r *http.Request) {
	claimID, err := uuid.Parse(r.PathValue("id"))
	if err != nil || claimID == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeClaimDeleteV1HandlerIDInvalid, apperrors.MsgClaimDeleteV1HandlerIDInvalid))
		return
	}
	var body deleteClaimBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeClaimPostV1HandlerBodyInvalid, apperrors.MsgClaimPostV1HandlerBodyInvalid, err))
		return
	}
	ownerUUID, err := uuid.Parse(strings.TrimSpace(body.OwnerUUID))
	if err != nil || ownerUUID == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeClaimPostV1HandlerBodyInvalid, apperrors.MsgClaimPostV1HandlerBodyInvalid))
		return
	}
	if err := h.claims.Delete(r.Context(), claimID, ownerUUID); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *internalHandler) claimPermission(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	mcUUID, err := uuid.Parse(strings.TrimSpace(q.Get("minecraftUuid")))
	if err != nil || mcUUID == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeClaimPermV1HandlerParamsInvalid, apperrors.MsgClaimPermV1HandlerParamsInvalid))
		return
	}
	serverSlug := strings.TrimSpace(q.Get("serverSlug"))
	world := strings.TrimSpace(q.Get("world"))
	if serverSlug == "" || world == "" {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeClaimPermV1HandlerParamsInvalid, apperrors.MsgClaimPermV1HandlerParamsInvalid))
		return
	}
	x, err := strconv.Atoi(strings.TrimSpace(q.Get("x")))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeClaimPermV1HandlerParamsInvalid, apperrors.MsgClaimPermV1HandlerParamsInvalid))
		return
	}
	z, err := strconv.Atoi(strings.TrimSpace(q.Get("z")))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeClaimPermV1HandlerParamsInvalid, apperrors.MsgClaimPermV1HandlerParamsInvalid))
		return
	}
	out, err := h.claims.CheckPermission(r.Context(), mcUUID, serverSlug, world, x, z)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

type createAllianceBody struct {
	LeaderUUID string `json:"leaderUuid"`
	GuildAID   string `json:"guildAId"`
	GuildBID   string `json:"guildBId"`
}

func (h *internalHandler) createAlliance(w http.ResponseWriter, r *http.Request) {
	var body createAllianceBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeAlliancePostV1HandlerBodyInvalid, apperrors.MsgAlliancePostV1HandlerBodyInvalid, err))
		return
	}
	leaderUUID, err := uuid.Parse(strings.TrimSpace(body.LeaderUUID))
	if err != nil || leaderUUID == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeAlliancePostV1HandlerBodyInvalid, apperrors.MsgAlliancePostV1HandlerBodyInvalid))
		return
	}
	guildAID, err := uuid.Parse(strings.TrimSpace(body.GuildAID))
	if err != nil || guildAID == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeAlliancePostV1HandlerBodyInvalid, apperrors.MsgAlliancePostV1HandlerBodyInvalid))
		return
	}
	guildBID, err := uuid.Parse(strings.TrimSpace(body.GuildBID))
	if err != nil || guildBID == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeAlliancePostV1HandlerBodyInvalid, apperrors.MsgAlliancePostV1HandlerBodyInvalid))
		return
	}
	alliance, err := h.alliances.Create(r.Context(), leaderUUID, guildAID, guildBID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, alliance)
}

type ingestAuditBody struct {
	Events []ingestAuditEventBody `json:"events"`
}

type ingestAuditEventBody struct {
	MinecraftUUID string  `json:"minecraftUuid"`
	ServerSlug    string  `json:"serverSlug"`
	World         string  `json:"world"`
	EventType     string  `json:"eventType"`
	BlockX        *int    `json:"blockX"`
	BlockY        *int    `json:"blockY"`
	BlockZ        *int    `json:"blockZ"`
	BlockType     string  `json:"blockType"`
	ClaimID       *string `json:"claimId"`
	OccurredAt    string  `json:"occurredAt"`
}

func (h *internalHandler) ingestAuditEvents(w http.ResponseWriter, r *http.Request) {
	var body ingestAuditBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeAuditIngestV1HandlerBodyInvalid, apperrors.MsgAuditIngestV1HandlerBodyInvalid, err))
		return
	}
	if len(body.Events) == 0 {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeAuditIngestV1HandlerBodyInvalid, apperrors.MsgAuditIngestV1HandlerBodyInvalid))
		return
	}
	events := make([]auditsvc.IngestEvent, 0, len(body.Events))
	for _, ev := range body.Events {
		mcUUID, err := uuid.Parse(strings.TrimSpace(ev.MinecraftUUID))
		if err != nil || mcUUID == uuid.Nil || strings.TrimSpace(ev.EventType) == "" {
			continue
		}
		var claimID *uuid.UUID
		if ev.ClaimID != nil && strings.TrimSpace(*ev.ClaimID) != "" {
			parsed, err := uuid.Parse(strings.TrimSpace(*ev.ClaimID))
			if err == nil && parsed != uuid.Nil {
				claimID = &parsed
			}
		}
		var occurred time.Time
		if ev.OccurredAt != "" {
			if t, err := time.Parse(time.RFC3339, ev.OccurredAt); err == nil {
				occurred = t
			}
		}
		events = append(events, auditsvc.IngestEvent{
			MinecraftUUID: mcUUID,
			ServerSlug:    ev.ServerSlug,
			World:         ev.World,
			EventType:     ev.EventType,
			BlockX:        ev.BlockX,
			BlockY:        ev.BlockY,
			BlockZ:        ev.BlockZ,
			BlockType:     ev.BlockType,
			ClaimID:       claimID,
			OccurredAt:    occurred,
		})
	}
	out, err := h.audit.IngestBatch(r.Context(), events)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

type createRollbackBody struct {
	TargetUUID    string `json:"targetUuid"`
	ActorUUID     string `json:"actorUuid"`
	ServerSlug    string `json:"serverSlug"`
	World         string `json:"world"`
	WindowMinutes int    `json:"windowMinutes"`
}

func (h *internalHandler) createRollback(w http.ResponseWriter, r *http.Request) {
	var body createRollbackBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeRollbackPostV1HandlerBodyInvalid, apperrors.MsgRollbackPostV1HandlerBodyInvalid, err))
		return
	}
	targetUUID, err := uuid.Parse(strings.TrimSpace(body.TargetUUID))
	if err != nil || targetUUID == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeRollbackPostV1HandlerBodyInvalid, apperrors.MsgRollbackPostV1HandlerBodyInvalid))
		return
	}
	actorUUID, err := uuid.Parse(strings.TrimSpace(body.ActorUUID))
	if err != nil || actorUUID == uuid.Nil || strings.TrimSpace(body.ServerSlug) == "" || body.WindowMinutes <= 0 {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeRollbackPostV1HandlerBodyInvalid, apperrors.MsgRollbackPostV1HandlerBodyInvalid))
		return
	}
	out, err := h.rollback.Create(r.Context(), rollbacksvc.CreateInput{
		TargetMinecraftUUID: targetUUID,
		ActorMinecraftUUID:  actorUUID,
		ServerSlug:          body.ServerSlug,
		World:               body.World,
		WindowMinutes:       body.WindowMinutes,
	})
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (h *internalHandler) listRollbackItems(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeRollbackGetV1HandlerIDInvalid, apperrors.MsgRollbackGetV1HandlerIDInvalid))
		return
	}
	items, err := h.rollback.ListItems(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

type completeRollbackBody struct {
	AppliedCount int `json:"appliedCount"`
}

func (h *internalHandler) completeRollback(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeRollbackGetV1HandlerIDInvalid, apperrors.MsgRollbackGetV1HandlerIDInvalid))
		return
	}
	var body completeRollbackBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeRollbackCompleteV1HandlerBodyInvalid, apperrors.MsgRollbackCompleteV1HandlerBodyInvalid, err))
		return
	}
	rb, err := h.rollback.Complete(r.Context(), id, body.AppliedCount)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, rb)
}
