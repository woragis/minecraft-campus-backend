package httpserver

import "net/http"

func Mount(mux *http.ServeMux, app *App) {
	ph := newPlayerHandler(app.Players, app.Guilds)
	ih := newInviteHandler(app.Invites, app.Players)
	gh := newGuildHandler(app.Guilds)
	th := newTrustHandler(app.Trust)
	ch := newCityHandler(app.Cities)
	clh := newClaimHandler(app.Claims)
	ah := newAllianceHandler(app.Alliances)
	auh := newAuditHandler(app.Audit)
	rbh := newRollbackHandler(app.Rollback)
	mh := newMetricsHandler(app.Metrics)
	alh := newAlertsHandler(app.Alerts)
	prh := newPresenceHandler(app.Presence)
	sth := newStatsHandler(app.Stats)
	wh := newWebHandler(app.WebAuth)
	mhMe := newMeHandler(app.Players, app.Invites, app.Guilds, wh)
	in := newInternalHandler(app.PluginAPIKey, app.Players, app.Presence, app.Stats, app.Invites, app.Guilds, app.Trust, app.Cities, app.Claims, app.Alliances, app.Audit, app.Rollback, app.WebAuth)

	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /ready", handleReady(app.DB))

	mux.HandleFunc("GET /v1/internal/whitelist/bedrock/{xuid}", in.requirePluginKey(in.bedrockWhitelist))
	mux.HandleFunc("GET /v1/internal/whitelist/{minecraftUuid}", in.requirePluginKey(in.whitelist))
	mux.HandleFunc("POST /v1/internal/players/bedrock/upsert", in.requirePluginKey(in.upsertBedrockPlayer))
	mux.HandleFunc("POST /v1/internal/players/upsert", in.requirePluginKey(in.upsertPlayer))
	mux.HandleFunc("POST /v1/internal/invites", in.requirePluginKey(in.createInvite))
	mux.HandleFunc("POST /v1/internal/guilds", in.requirePluginKey(in.createGuild))
	mux.HandleFunc("POST /v1/internal/guilds/{id}/join", in.requirePluginKey(in.joinGuild))
	mux.HandleFunc("POST /v1/internal/guilds/by-slug/{slug}/join", in.requirePluginKey(in.joinGuildBySlug))
	mux.HandleFunc("POST /v1/internal/guilds/{id}/leave", in.requirePluginKey(in.leaveGuild))
	mux.HandleFunc("POST /v1/internal/guilds/by-slug/{slug}/leave", in.requirePluginKey(in.leaveGuildBySlug))
	mux.HandleFunc("POST /v1/internal/web/link-codes", in.requirePluginKey(in.createWebLinkCode))
	mux.HandleFunc("POST /v1/internal/trust/events", in.requirePluginKey(in.recordTrustEvent))
	mux.HandleFunc("POST /v1/internal/cities", in.requirePluginKey(in.createCity))
	mux.HandleFunc("POST /v1/internal/claims", in.requirePluginKey(in.createClaim))
	mux.HandleFunc("DELETE /v1/internal/claims/{id}", in.requirePluginKey(in.deleteClaim))
	mux.HandleFunc("GET /v1/internal/claims/permission", in.requirePluginKey(in.claimPermission))
	mux.HandleFunc("POST /v1/internal/alliances", in.requirePluginKey(in.createAlliance))
	mux.HandleFunc("POST /v1/internal/audit/events", in.requirePluginKey(in.ingestAuditEvents))
	mux.HandleFunc("POST /v1/internal/rollbacks", in.requirePluginKey(in.createRollback))
	mux.HandleFunc("GET /v1/internal/rollbacks/{id}/items", in.requirePluginKey(in.listRollbackItems))
	mux.HandleFunc("POST /v1/internal/rollbacks/{id}/complete", in.requirePluginKey(in.completeRollback))

	mux.HandleFunc("POST /v1/internal/stats/ingest", in.requirePluginKey(in.statsIngest))
	mux.HandleFunc("GET /v1/internal/players/{id}/hud", in.requirePluginKey(in.playerHUD))

	mux.HandleFunc("POST /v1/internal/presence/online", in.requirePluginKey(in.presenceOnline))
	mux.HandleFunc("POST /v1/internal/presence/offline", in.requirePluginKey(in.presenceOffline))
	mux.HandleFunc("POST /v1/internal/presence/heartbeat", in.requirePluginKey(in.presenceHeartbeat))

	mux.HandleFunc("GET /v1/presence/overview", prh.overview)
	mux.HandleFunc("GET /v1/presence/servers/{slug}", prh.server)
	mux.HandleFunc("GET /v1/presence/guilds/{id}", prh.guild)

	mux.HandleFunc("GET /v1/rollbacks/{id}", rbh.getByID)

	mux.HandleFunc("GET /v1/metrics/overview", mh.overview)
	mux.HandleFunc("GET /v1/metrics/territory", mh.territory)
	mux.HandleFunc("GET /v1/alerts", alh.list)
	mux.HandleFunc("POST /v1/alerts/{id}/acknowledge", alh.acknowledge)

	mux.HandleFunc("GET /v1/lookup/cities/{slug}", ch.getBySlug)
	mux.HandleFunc("GET /v1/cities/{id}/claims", clh.listByCity)
	mux.HandleFunc("GET /v1/cities/{id}", ch.getByID)
	mux.HandleFunc("GET /v1/cities", ch.list)

	mux.HandleFunc("GET /v1/claims/{id}", clh.getByID)

	mux.HandleFunc("GET /v1/guilds/{id}/alliances", ah.listByGuild)
	mux.HandleFunc("GET /v1/lookup/guilds/{slug}", gh.getBySlug)
	mux.HandleFunc("GET /v1/guilds/{id}/members", gh.listMembers)
	mux.HandleFunc("GET /v1/guilds/{id}", gh.getByID)
	mux.HandleFunc("GET /v1/guilds", gh.list)

	mux.HandleFunc("GET /v1/lookup/players/minecraft/{minecraftUuid}", ph.getByMinecraftUUID)
	mux.HandleFunc("GET /v1/players/{id}/stats", sth.playerStats)
	mux.HandleFunc("GET /v1/players/{id}/audit-events", auh.listByPlayer)
	mux.HandleFunc("GET /v1/players/{id}/trust-events", th.listEvents)
	mux.HandleFunc("GET /v1/players/{id}/sponsor-tree", th.sponsorTree)
	mux.HandleFunc("GET /v1/players/{id}/invites", ph.listInvitees)
	mux.HandleFunc("GET /v1/players/{id}", ph.getByID)
	mux.HandleFunc("GET /v1/invites/{code}", ih.getByCode)

	mux.HandleFunc("POST /v1/web/session", wh.createSession)
	mux.HandleFunc("POST /v1/web/logout", wh.logout)
	mhMe.mount(mux)
}
