package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateEvent(ctx context.Context, event *models.TrustEvent) error {
	if err := r.db.WithContext(ctx).Create(event).Error; err != nil {
		return fmt.Errorf("trust event create: %w", err)
	}
	return nil
}

func (r *Repository) ListByPlayer(ctx context.Context, playerID uuid.UUID, limit int) ([]models.TrustEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	var events []models.TrustEvent
	err := r.db.WithContext(ctx).
		Where("player_id = ?", playerID).
		Order("created_at DESC").
		Limit(limit).
		Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("trust events list: %w", err)
	}
	return events, nil
}

func (r *Repository) CountRecentBansAmongInvitees(ctx context.Context, sponsorID uuid.UUID, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Player{}).
		Where("invited_by_id = ? AND status = ? AND updated_at >= ?", sponsorID, models.PlayerStatusBanned, since).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("recent invitee bans: %w", err)
	}
	return count, nil
}
