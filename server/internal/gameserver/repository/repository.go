package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetOrCreateBySlug(ctx context.Context, slug, name string) (*models.GameServer, error) {
	slug = strings.TrimSpace(strings.ToLower(slug))
	if slug == "" {
		slug = "vanilla"
	}
	if strings.TrimSpace(name) == "" {
		name = slug
	}

	var gs models.GameServer
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&gs).Error
	if err == nil {
		return &gs, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("game server find by slug: %w", err)
	}

	gs = models.GameServer{
		ID:   uuid.New(),
		Slug: slug,
		Name: name,
	}
	if err := r.db.WithContext(ctx).Create(&gs).Error; err != nil {
		return nil, fmt.Errorf("game server create: %w", err)
	}
	return &gs, nil
}

func (r *Repository) TouchPlayerPresence(ctx context.Context, serverID, playerID uuid.UUID, seenAt time.Time) error {
	sp := models.ServerPlayer{
		ServerID:   serverID,
		PlayerID:   playerID,
		LastSeenAt: seenAt,
	}
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "server_id"}, {Name: "player_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"last_seen_at"}),
		}).
		Create(&sp).Error
	if err != nil {
		return fmt.Errorf("server player touch: %w", err)
	}
	return nil
}

func (r *Repository) AddPlayerStats(ctx context.Context, serverID, playerID uuid.UUID, sessionSeconds, mobKills int64, seenAt time.Time) error {
	if sessionSeconds < 0 {
		sessionSeconds = 0
	}
	if mobKills < 0 {
		mobKills = 0
	}
	if sessionSeconds == 0 && mobKills == 0 {
		return r.TouchPlayerPresence(ctx, serverID, playerID, seenAt)
	}
	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO server_players (server_id, player_id, last_seen_at, play_time_secs, mob_kills)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (server_id, player_id) DO UPDATE SET
			last_seen_at = EXCLUDED.last_seen_at,
			play_time_secs = server_players.play_time_secs + EXCLUDED.play_time_secs,
			mob_kills = server_players.mob_kills + EXCLUDED.mob_kills
	`, serverID, playerID, seenAt, sessionSeconds, mobKills).Error
	if err != nil {
		return fmt.Errorf("server player add stats: %w", err)
	}
	return nil
}

type PlayerServerStat struct {
	ServerSlug   string `gorm:"column:server_slug"`
	PlayTimeSecs int64  `gorm:"column:play_time_secs"`
	MobKills     int64  `gorm:"column:mob_kills"`
}

func (r *Repository) ListPlayerStats(ctx context.Context, playerID uuid.UUID) ([]PlayerServerStat, error) {
	var rows []PlayerServerStat
	err := r.db.WithContext(ctx).Raw(`
		SELECT gs.slug AS server_slug, sp.play_time_secs, sp.mob_kills
		FROM server_players sp
		JOIN game_servers gs ON gs.id = sp.server_id
		WHERE sp.player_id = ?
		ORDER BY sp.play_time_secs DESC
	`, playerID).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("server player list stats: %w", err)
	}
	return rows, nil
}
