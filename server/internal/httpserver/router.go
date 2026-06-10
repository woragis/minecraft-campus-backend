package httpserver

import "net/http"

func Mount(mux *http.ServeMux, app *App) {
	ph := newPlayerHandler(app.Players)
	ih := newInviteHandler(app.Invites, app.Players)
	in := newInternalHandler(app.PluginAPIKey, app.Players, app.Invites)

	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /ready", handleReady(app.DB))

	mux.HandleFunc("GET /v1/internal/whitelist/{minecraftUuid}", in.requirePluginKey(in.whitelist))
	mux.HandleFunc("POST /v1/internal/players/upsert", in.requirePluginKey(in.upsertPlayer))
	mux.HandleFunc("POST /v1/internal/invites", in.requirePluginKey(in.createInvite))

	mux.HandleFunc("GET /v1/players/minecraft/{minecraftUuid}", ph.getByMinecraftUUID)
	mux.HandleFunc("GET /v1/players/{id}/invites", ph.listInvitees)
	mux.HandleFunc("GET /v1/players/{id}", ph.getByID)
	mux.HandleFunc("GET /v1/invites/{code}", ih.getByCode)
}
