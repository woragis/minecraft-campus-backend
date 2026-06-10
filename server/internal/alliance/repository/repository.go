package repository

import (
	"context"
	"errors"
	"fmt"

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

func (r *Repository) Create(ctx context.Context, alliance *models.Alliance) error {
	if err := r.db.WithContext(ctx).Create(alliance).Error; err != nil {
		return fmt.Errorf("alliance create: %w", err)
	}
	return nil
}

func (r *Repository) Exists(ctx context.Context, guildA, guildB uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Alliance{}).
		Where("(guild_a_id = ? AND guild_b_id = ?) OR (guild_a_id = ? AND guild_b_id = ?)", guildA, guildB, guildB, guildA).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) ListByGuild(ctx context.Context, guildID uuid.UUID) ([]models.Alliance, error) {
	var alliances []models.Alliance
	err := r.db.WithContext(ctx).
		Where("guild_a_id = ? OR guild_b_id = ?", guildID, guildID).
		Order("created_at DESC").
		Find(&alliances).Error
	if err != nil {
		return nil, fmt.Errorf("alliance list: %w", err)
	}
	return alliances, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Alliance{})
	if res.Error != nil {
		return fmt.Errorf("alliance delete: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.Alliance, error) {
	var a models.Alliance
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("alliance find: %w", err)
	}
	return &a, nil
}
