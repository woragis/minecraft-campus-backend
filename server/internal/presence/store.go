package presence

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OnlinePlayer struct {
	PlayerID   uuid.UUID `json:"playerId"`
	Username   string    `json:"username"`
	ServerSlug string    `json:"serverSlug"`
	Platform   string    `json:"platform,omitempty"`
	GuildID    *uuid.UUID `json:"guildId,omitempty"`
	Since      time.Time `json:"since"`
}

type ServerPresence struct {
	Slug         string         `json:"slug"`
	OnlineCount  int            `json:"onlineCount"`
	Players      []OnlinePlayer `json:"players"`
}

type Overview struct {
	Enabled     bool             `json:"enabled"`
	TotalOnline int              `json:"totalOnline"`
	Servers     []ServerPresence `json:"servers"`
}

type GuildPresence struct {
	GuildID      uuid.UUID      `json:"guildId"`
	OnlineCount  int            `json:"onlineCount"`
	Members      []OnlinePlayer `json:"members"`
}

type MarkOnlineInput struct {
	PlayerID   uuid.UUID
	Username   string
	ServerSlug string
	Platform   string
	GuildID    *uuid.UUID
}

type Store interface {
	Enabled() bool
	MarkOnline(ctx context.Context, in MarkOnlineInput) error
	MarkOffline(ctx context.Context, playerID uuid.UUID, serverSlug string, guildID *uuid.UUID) error
	Heartbeat(ctx context.Context, playerID uuid.UUID, serverSlug string) error
	ListServerOnline(ctx context.Context, serverSlug string) ([]OnlinePlayer, error)
	ListServersOnline(ctx context.Context, serverSlugs []string) ([]ServerPresence, error)
	CountGuildOnline(ctx context.Context, guildID uuid.UUID, memberIDs []uuid.UUID) ([]OnlinePlayer, error)
}
