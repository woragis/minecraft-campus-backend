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

func (r *Repository) Create(ctx context.Context, alert *models.Alert) error {
	if err := r.db.WithContext(ctx).Create(alert).Error; err != nil {
		return fmt.Errorf("alert create: %w", err)
	}
	return nil
}

func (r *Repository) ListUnacknowledged(ctx context.Context, limit int) ([]models.Alert, error) {
	if limit <= 0 {
		limit = 50
	}
	var alerts []models.Alert
	err := r.db.WithContext(ctx).Where("acknowledged = ?", false).Order("created_at DESC").Limit(limit).Find(&alerts).Error
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

func (r *Repository) Acknowledge(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Model(&models.Alert{}).Where("id = ?", id).Update("acknowledged", true)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

type GriefingRow struct {
	PlayerID uuid.UUID
	Count    int64
}

func (r *Repository) GriefingSpikes(ctx context.Context, since time.Time, threshold int64) ([]GriefingRow, error) {
	var rows []GriefingRow
	err := r.db.WithContext(ctx).Model(&models.AuditEvent{}).
		Select("player_id, COUNT(*) AS count").
		Where("event_type = ? AND claim_id IS NOT NULL AND occurred_at >= ?", models.AuditEventBlockBreak, since).
		Group("player_id").
		Having("COUNT(*) >= ?", threshold).
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("griefing spikes: %w", err)
	}
	return rows, nil
}

func (r *Repository) RecentDuplicate(ctx context.Context, alertType string, playerID uuid.UUID, since time.Time) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Alert{}).
		Where("alert_type = ? AND player_id = ? AND created_at >= ?", alertType, playerID, since).
		Count(&count).Error
	return count > 0, err
}
