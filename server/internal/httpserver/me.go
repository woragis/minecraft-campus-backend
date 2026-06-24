package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	affiliationsvc "github.com/woragis/minecraft-campus-backend/server/internal/affiliation/service"
	guildsvc "github.com/woragis/minecraft-campus-backend/server/internal/guild/service"
	invitesvc "github.com/woragis/minecraft-campus-backend/server/internal/invite/service"
	playersvc "github.com/woragis/minecraft-campus-backend/server/internal/player/service"
	"github.com/woragis/minecraft-campus-backend/server/internal/webauth"
)

type meHandler struct {
	players *playersvc.Service
	invites *invitesvc.Service
	guilds  *guildsvc.Service
	web     *webHandler
}

func newMeHandler(players *playersvc.Service, invites *invitesvc.Service, guilds *guildsvc.Service, web *webHandler) *meHandler {
	return &meHandler{players: players, invites: invites, guilds: guilds, web: web}
}

func (h *meHandler) mount(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/me", h.web.requireSession(h.profile))
	mux.HandleFunc("PATCH /v1/me/affiliation", h.web.requireSession(h.patchAffiliation))
	mux.HandleFunc("POST /v1/me/invites", h.web.requireSession(h.createInvite))
	mux.HandleFunc("POST /v1/me/guilds", h.web.requireSession(h.createGuild))
	mux.HandleFunc("POST /v1/me/guilds/{slug}/join", h.web.requireSession(h.joinGuild))
	mux.HandleFunc("POST /v1/me/guilds/{slug}/leave", h.web.requireSession(h.leaveGuild))
}

func (h *meHandler) profile(w http.ResponseWriter, r *http.Request, session *webauth.Session) {
	player, err := h.players.GetByID(r.Context(), session.PlayerID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	out := playerProfile{Player: player}
	guild, err := h.guilds.PlayerGuild(r.Context(), player.ID)
	if err == nil && guild != nil {
		out.Guild = guild
	}
	writeJSON(w, http.StatusOK, out)
}

type meInviteBody struct {
	TargetUsername  string `json:"targetUsername"`
	AffiliationType string `json:"affiliationType"`
}

func (h *meHandler) createInvite(w http.ResponseWriter, r *http.Request, session *webauth.Session) {
	var body meInviteBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeInvitePostInternalV1HandlerBodyInvalid, apperrors.MsgInvitePostInternalV1HandlerBodyInvalid, err))
		return
	}
	inv, err := h.invites.CreateForSponsor(r.Context(), session.MinecraftUUID, body.TargetUsername, body.AffiliationType)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, inv)
}

type meGuildBody struct {
	Name string `json:"name"`
}

func (h *meHandler) createGuild(w http.ResponseWriter, r *http.Request, session *webauth.Session) {
	var body meGuildBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeGuildPostV1HandlerBodyInvalid, apperrors.MsgGuildPostV1HandlerBodyInvalid, err))
		return
	}
	guild, err := h.guilds.Create(r.Context(), session.MinecraftUUID, body.Name)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, guild)
}

func (h *meHandler) joinGuild(w http.ResponseWriter, r *http.Request, session *webauth.Session) {
	if err := h.guilds.JoinBySlug(r.Context(), r.PathValue("slug"), session.MinecraftUUID); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "joined"})
}

func (h *meHandler) leaveGuild(w http.ResponseWriter, r *http.Request, session *webauth.Session) {
	if err := h.guilds.LeaveBySlug(r.Context(), r.PathValue("slug"), session.MinecraftUUID); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "left"})
}

type meAffiliationBody struct {
	AffiliationType string  `json:"affiliationType"`
	UniversitySlug  *string `json:"universitySlug"`
	FacultySlug     *string `json:"facultySlug"`
	CourseSlug      *string `json:"courseSlug"`
}

func (h *meHandler) patchAffiliation(w http.ResponseWriter, r *http.Request, session *webauth.Session) {
	var body meAffiliationBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.InvalidCause(apperrors.CodeAffiliationPatchV1HandlerBodyInvalid, apperrors.MsgAffiliationPatchV1HandlerBodyInvalid, err))
		return
	}
	player, err := h.players.UpdateAffiliation(r.Context(), session.PlayerID, affiliationsvc.AffiliationInput{
		AffiliationType: body.AffiliationType,
		UniversitySlug:  body.UniversitySlug,
		FacultySlug:     body.FacultySlug,
		CourseSlug:      body.CourseSlug,
	})
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, player)
}
