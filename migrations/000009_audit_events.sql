CREATE TABLE audit_events (
    id UUID PRIMARY KEY,
    server_id UUID NOT NULL REFERENCES game_servers(id),
    world TEXT NOT NULL,
    player_id UUID NOT NULL REFERENCES players(id),
    event_type TEXT NOT NULL,
    block_x INT,
    block_y INT,
    block_z INT,
    block_type TEXT,
    claim_id UUID REFERENCES claims(id),
    metadata JSONB NOT NULL DEFAULT '{}',
    occurred_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_events_player_occurred ON audit_events (player_id, occurred_at DESC);
CREATE INDEX idx_audit_events_claim_occurred ON audit_events (claim_id, occurred_at DESC) WHERE claim_id IS NOT NULL;
CREATE INDEX idx_audit_events_server_world_occurred ON audit_events (server_id, world, occurred_at DESC);
