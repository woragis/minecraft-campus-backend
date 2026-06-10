CREATE TABLE game_servers (
    id UUID PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE server_players (
    server_id UUID NOT NULL REFERENCES game_servers(id) ON DELETE CASCADE,
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    play_time_secs BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (server_id, player_id)
);

CREATE INDEX idx_server_players_player_id ON server_players (player_id);
