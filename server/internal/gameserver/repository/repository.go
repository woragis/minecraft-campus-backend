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
