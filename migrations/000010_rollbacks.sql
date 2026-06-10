CREATE TABLE rollbacks (
    id UUID PRIMARY KEY,
    target_player_id UUID NOT NULL REFERENCES players(id),
    actor_player_id UUID REFERENCES players(id),
    server_id UUID NOT NULL REFERENCES game_servers(id),
    world TEXT NOT NULL,
    from_at TIMESTAMPTZ NOT NULL,
    to_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    item_count INT NOT NULL DEFAULT 0,
    applied_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ
);

CREATE INDEX idx_rollbacks_target_created ON rollbacks (target_player_id, created_at DESC);

CREATE TABLE rollback_items (
    id UUID PRIMARY KEY,
    rollback_id UUID NOT NULL REFERENCES rollbacks(id) ON DELETE CASCADE,
    audit_event_id UUID REFERENCES audit_events(id),
    action TEXT NOT NULL,
    block_type TEXT NOT NULL,
    block_x INT NOT NULL,
    block_y INT NOT NULL,
    block_z INT NOT NULL,
    world TEXT NOT NULL,
    applied BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_rollback_items_rollback_id ON rollback_items (rollback_id);
