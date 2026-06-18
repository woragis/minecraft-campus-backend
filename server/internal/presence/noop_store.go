package presence

import (
	"context"

	"github.com/google/uuid"
)

type NoopStore struct{}

func NewNoopStore() *NoopStore { return &NoopStore{} }

func (s *NoopStore) Enabled() bool { return false }

func (s *NoopStore) MarkOnline(ctx context.Context, in MarkOnlineInput) error {
	return nil
}

func (s *NoopStore) MarkOffline(ctx context.Context, playerID uuid.UUID, serverSlug string, guildID *uuid.UUID) error {
	return nil
}

func (s *NoopStore) Heartbeat(ctx context.Context, playerID uuid.UUID, serverSlug string) error {
	return nil
}

func (s *NoopStore) ListServerOnline(ctx context.Context, serverSlug string) ([]OnlinePlayer, error) {
	return nil, nil
}

func (s *NoopStore) ListServersOnline(ctx context.Context, serverSlugs []string) ([]ServerPresence, error) {
	out := make([]ServerPresence, 0, len(serverSlugs))
	for _, slug := range serverSlugs {
		out = append(out, ServerPresence{Slug: slug, OnlineCount: 0, Players: []OnlinePlayer{}})
	}
	return out, nil
}

func (s *NoopStore) CountGuildOnline(ctx context.Context, guildID uuid.UUID, memberIDs []uuid.UUID) ([]OnlinePlayer, error) {
	return nil, nil
}
