CREATE TABLE world_snapshots (
    id UUID PRIMARY KEY,
    snapshot_type TEXT NOT NULL,
    storage TEXT NOT NULL,
    path TEXT NOT NULL,
    size_bytes BIGINT NOT NULL DEFAULT 0,
    checksum TEXT,
    status TEXT NOT NULL DEFAULT 'completed',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_world_snapshots_created ON world_snapshots (created_at DESC);
