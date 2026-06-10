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

func (r *Repository) CreateBatch(ctx context.Context, events []models.AuditEvent) error {
	if len(events) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).Create(&events).Error; err != nil {
		return fmt.Errorf("audit create batch: %w", err)
	}
	return nil
}

type ListFilter struct {
	PlayerID uuid.UUID
	From     *time.Time
	To       *time.Time
	EventType string
	Limit    int
}

func (r *Repository) ListByPlayer(ctx context.Context, f ListFilter) ([]models.AuditEvent, error) {
	if f.Limit <= 0 {
		f.Limit = 100
	}
	q := r.db.WithContext(ctx).Where("player_id = ?", f.PlayerID)
	if f.From != nil {
		q = q.Where("occurred_at >= ?", *f.From)
	}
	if f.To != nil {
		q = q.Where("occurred_at <= ?", *f.To)
	}
	if f.EventType != "" {
		q = q.Where("event_type = ?", f.EventType)
	}
	var events []models.AuditEvent
	err := q.Order("occurred_at DESC").Limit(f.Limit).Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("audit list: %w", err)
	}
	return events, nil
}

func (r *Repository) ListForRollback(ctx context.Context, playerID, serverID uuid.UUID, world string, from, to time.Time) ([]models.AuditEvent, error) {
	var events []models.AuditEvent
	err := r.db.WithContext(ctx).
		Where("player_id = ? AND server_id = ? AND world = ?", playerID, serverID, world).
		Where("occurred_at >= ? AND occurred_at <= ?", from, to).
		Where("claim_id IS NOT NULL").
		Where("event_type IN ?", []string{models.AuditEventBlockPlace, models.AuditEventBlockBreak}).
		Order("occurred_at DESC").
		Limit(5000).
		Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("audit list rollback: %w", err)
	}
	return events, nil
}
