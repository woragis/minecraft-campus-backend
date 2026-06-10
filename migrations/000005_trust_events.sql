CREATE TABLE trust_events (
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    delta INT NOT NULL DEFAULT 0,
    reason TEXT,
    actor_player_id UUID REFERENCES players(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_trust_events_player_id ON trust_events (player_id);
CREATE INDEX idx_trust_events_created_at ON trust_events (created_at);
