CREATE TABLE alliances (
    id UUID PRIMARY KEY,
    guild_a_id UUID NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
    guild_b_id UUID NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (guild_a_id <> guild_b_id),
    UNIQUE (guild_a_id, guild_b_id)
);

CREATE INDEX idx_alliances_guild_a ON alliances (guild_a_id);
CREATE INDEX idx_alliances_guild_b ON alliances (guild_b_id);
