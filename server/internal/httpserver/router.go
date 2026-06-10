package httpserver

import "net/http"

func Mount(mux *http.ServeMux, app *App) {
	ph := newPlayerHandler(app.Players, app.Guilds)
	ih := newInviteHandler(app.Invites, app.Players)
	gh := newGuildHandler(app.Guilds)
	th := newTrustHandler(app.Trust)
	in := newInternalHandler(app.PluginAPIKey, app.Players, app.Invites, app.Guilds, app.Trust)

	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /ready", handleReady(app.DB))

	mux.HandleFunc("GET /v1/internal/whitelist/{minecraftUuid}", in.requirePluginKey(in.whitelist))
	mux.HandleFunc("POST /v1/internal/players/upsert", in.requirePluginKey(in.upsertPlayer))
	mux.HandleFunc("POST /v1/internal/invites", in.requirePluginKey(in.createInvite))
	mux.HandleFunc("POST /v1/internal/guilds", in.requirePluginKey(in.createGuild))
	mux.HandleFunc("POST /v1/internal/guilds/{id}/join", in.requirePluginKey(in.joinGuild))
	mux.HandleFunc("POST /v1/internal/guilds/{id}/leave", in.requirePluginKey(in.leaveGuild))
	mux.HandleFunc("POST /v1/internal/trust/events", in.requirePluginKey(in.recordTrustEvent))

	mux.HandleFunc("GET /v1/guilds/slug/{slug}", gh.getBySlug)
	mux.HandleFunc("GET /v1/guilds/{id}/members", gh.listMembers)
	mux.HandleFunc("GET /v1/guilds/{id}", gh.getByID)
	mux.HandleFunc("GET /v1/guilds", gh.list)

	mux.HandleFunc("GET /v1/players/minecraft/{minecraftUuid}", ph.getByMinecraftUUID)
	mux.HandleFunc("GET /v1/players/{id}/trust-events", th.listEvents)
	mux.HandleFunc("GET /v1/players/{id}/sponsor-tree", th.sponsorTree)
	mux.HandleFunc("GET /v1/players/{id}/invites", ph.listInvitees)
	mux.HandleFunc("GET /v1/players/{id}", ph.getByID)
	mux.HandleFunc("GET /v1/invites/{code}", ih.getByCode)
}
