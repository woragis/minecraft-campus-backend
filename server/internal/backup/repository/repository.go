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

func (r *Repository) Create(ctx context.Context, snap *models.WorldSnapshot) error {
	if err := r.db.WithContext(ctx).Create(snap).Error; err != nil {
		return fmt.Errorf("snapshot create: %w", err)
	}
	return nil
}

func (r *Repository) DeleteOlderThan(ctx context.Context, before time.Time) (int64, error) {
	res := r.db.WithContext(ctx).Where("created_at < ?", before).Delete(&models.WorldSnapshot{})
	if res.Error != nil {
		return 0, fmt.Errorf("snapshot purge: %w", res.Error)
	}
	return res.RowsAffected, nil
}

func (r *Repository) ListRecent(ctx context.Context, limit int) ([]models.WorldSnapshot, error) {
	if limit <= 0 {
		limit = 20
	}
	var snaps []models.WorldSnapshot
	err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&snaps).Error
	if err != nil {
		return nil, err
	}
	return snaps, nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.WorldSnapshot, error) {
	var snap models.WorldSnapshot
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&snap).Error
	if err != nil {
		return nil, err
	}
	return &snap, nil
}
