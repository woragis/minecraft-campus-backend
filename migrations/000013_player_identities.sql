CREATE TABLE player_identities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    platform TEXT NOT NULL CHECK (platform IN ('java', 'bedrock')),
    external_id TEXT NOT NULL,
    username TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (platform, external_id)
);

CREATE INDEX idx_player_identities_player_id ON player_identities (player_id);

INSERT INTO game_servers (id, slug, name)
VALUES
    ('a1b2c3d4-e5f6-4789-a012-3456789abcde', 'bedrock', 'CampusWorld Bedrock'),
    ('b2c3d4e5-f6a7-4890-b123-456789abcdef', 'crossplay', 'CampusWorld Cross-play')
ON CONFLICT (slug) DO NOTHING;
