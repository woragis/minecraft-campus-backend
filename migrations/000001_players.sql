CREATE TABLE players (
    id UUID PRIMARY KEY,
    minecraft_uuid UUID NOT NULL UNIQUE,
    username TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'probation',
    invited_by_id UUID REFERENCES players(id),
    trust_score INT NOT NULL DEFAULT 100,
    sponsor_score INT NOT NULL DEFAULT 100,
    probation_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_players_username_lower ON players (lower(username));
CREATE INDEX idx_players_status ON players (status);
CREATE INDEX idx_players_invited_by_id ON players (invited_by_id);
