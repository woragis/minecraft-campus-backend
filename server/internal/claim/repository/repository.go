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

func (r *Repository) Create(ctx context.Context, claim *models.Claim) error {
	if err := r.db.WithContext(ctx).Create(claim).Error; err != nil {
		return fmt.Errorf("claim create: %w", err)
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.Claim, error) {
	var claim models.Claim
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&claim).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("claim find by id: %w", err)
	}
	return &claim, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Claim{})
	if res.Error != nil {
		return fmt.Errorf("claim delete: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) TotalAreaByOwner(ctx context.Context, ownerID uuid.UUID) (int, error) {
	type row struct{ Total int }
	var result row
	err := r.db.WithContext(ctx).Model(&models.Claim{}).
		Select("COALESCE(SUM(area_blocks), 0) AS total").
		Where("owner_id = ?", ownerID).
		Scan(&result).Error
	if err != nil {
		return 0, fmt.Errorf("claim total area: %w", err)
	}
	return result.Total, nil
}

func (r *Repository) HasOverlap(ctx context.Context, serverID uuid.UUID, world string, minX, maxX, minZ, maxZ int, excludeID *uuid.UUID) (bool, error) {
	q := r.db.WithContext(ctx).Model(&models.Claim{}).
		Where("server_id = ? AND world = ?", serverID, world).
		Where("NOT (max_x < ? OR min_x > ? OR max_z < ? OR min_z > ?)", minX, maxX, minZ, maxZ)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, fmt.Errorf("claim overlap: %w", err)
	}
	return count > 0, nil
}

func (r *Repository) FindAt(ctx context.Context, serverID uuid.UUID, world string, x, z int) (*models.Claim, error) {
	var claim models.Claim
	err := r.db.WithContext(ctx).
		Where("server_id = ? AND world = ?", serverID, world).
		Where("min_x <= ? AND max_x >= ? AND min_z <= ? AND max_z >= ?", x, x, z, z).
		First(&claim).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("claim find at: %w", err)
	}
	return &claim, nil
}

func (r *Repository) ListByCity(ctx context.Context, cityID uuid.UUID) ([]models.Claim, error) {
	var claims []models.Claim
	err := r.db.WithContext(ctx).Where("city_id = ?", cityID).Order("created_at ASC").Find(&claims).Error
	if err != nil {
		return nil, fmt.Errorf("claims by city: %w", err)
	}
	return claims, nil
}

func (r *Repository) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.Claim, error) {
	var claims []models.Claim
	err := r.db.WithContext(ctx).Where("owner_id = ?", ownerID).Order("created_at ASC").Find(&claims).Error
	if err != nil {
		return nil, fmt.Errorf("claims by owner: %w", err)
	}
	return claims, nil
}

func (r *Repository) SumAreaByCity(ctx context.Context, cityID uuid.UUID) (int, error) {
	type row struct{ Total int }
	var result row
	err := r.db.WithContext(ctx).Model(&models.Claim{}).
		Select("COALESCE(SUM(area_blocks), 0) AS total").
		Where("city_id = ?", cityID).
		Scan(&result).Error
	if err != nil {
		return 0, err
	}
	return result.Total, nil
}
