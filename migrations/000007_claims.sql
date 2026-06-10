CREATE TABLE claims (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL REFERENCES players(id),
    city_id UUID REFERENCES cities(id),
    guild_id UUID REFERENCES guilds(id),
    server_id UUID NOT NULL REFERENCES game_servers(id),
    world TEXT NOT NULL DEFAULT 'world',
    min_x INT NOT NULL,
    min_y INT NOT NULL DEFAULT -64,
    max_x INT NOT NULL,
    max_y INT NOT NULL DEFAULT 320,
    min_z INT NOT NULL,
    max_z INT NOT NULL,
    zone_type TEXT NOT NULL DEFAULT 'urban',
    area_blocks INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (min_x <= max_x),
    CHECK (min_z <= max_z),
    CHECK (min_y <= max_y)
);

CREATE INDEX idx_claims_owner_id ON claims (owner_id);
CREATE INDEX idx_claims_city_id ON claims (city_id);
CREATE INDEX idx_claims_server_world ON claims (server_id, world);
CREATE INDEX idx_claims_bounds ON claims (server_id, world, min_x, max_x, min_z, max_z);
