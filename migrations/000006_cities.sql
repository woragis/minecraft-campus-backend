CREATE TABLE cities (
    id UUID PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    founder_id UUID NOT NULL REFERENCES players(id),
    guild_id UUID REFERENCES guilds(id),
    server_id UUID NOT NULL REFERENCES game_servers(id),
    world TEXT NOT NULL DEFAULT 'world',
    center_x INT NOT NULL DEFAULT 0,
    center_z INT NOT NULL DEFAULT 0,
    area_blocks INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_cities_guild_id ON cities (guild_id);
CREATE INDEX idx_cities_server_id ON cities (server_id);
