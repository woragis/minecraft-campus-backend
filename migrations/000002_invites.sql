CREATE TABLE invites (
    id UUID PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    sponsor_id UUID NOT NULL REFERENCES players(id),
    target_username TEXT NOT NULL,
    invited_player_id UUID REFERENCES players(id),
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    accepted_at TIMESTAMPTZ
);

CREATE INDEX idx_invites_sponsor_id ON invites (sponsor_id);
CREATE INDEX idx_invites_code ON invites (code);
CREATE INDEX idx_invites_pending_target_username ON invites (lower(target_username))
    WHERE status = 'pending';
