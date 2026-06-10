package repository

import (
	"context"
	"errors"
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

func (r *Repository) Create(ctx context.Context, rollback *models.Rollback, items []models.RollbackItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(rollback).Error; err != nil {
			return fmt.Errorf("rollback create: %w", err)
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return fmt.Errorf("rollback items create: %w", err)
			}
		}
		return nil
	})
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.Rollback, error) {
	var rb models.Rollback
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&rb).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("rollback find: %w", err)
	}
	return &rb, nil
}

func (r *Repository) ListItems(ctx context.Context, rollbackID uuid.UUID) ([]models.RollbackItem, error) {
	var items []models.RollbackItem
	err := r.db.WithContext(ctx).
		Where("rollback_id = ?", rollbackID).
		Order("created_at ASC").
		Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("rollback items list: %w", err)
	}
	return items, nil
}

func (r *Repository) MarkCompleted(ctx context.Context, id uuid.UUID, appliedCount int, status string) error {
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).Model(&models.Rollback{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":        status,
			"applied_count": appliedCount,
			"completed_at":  now,
		})
	if res.Error != nil {
		return fmt.Errorf("rollback complete: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	res := r.db.WithContext(ctx).Model(&models.Rollback{}).Where("id = ?", id).Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
