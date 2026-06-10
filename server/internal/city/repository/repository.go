package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

func (r *Repository) Create(ctx context.Context, city *models.City) error {
	if err := r.db.WithContext(ctx).Create(city).Error; err != nil {
		return fmt.Errorf("city create: %w", err)
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.City, error) {
	var city models.City
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&city).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("city find by id: %w", err)
	}
	return &city, nil
}

func (r *Repository) FindBySlug(ctx context.Context, slug string) (*models.City, error) {
	var city models.City
	err := r.db.WithContext(ctx).Where("slug = ?", strings.ToLower(strings.TrimSpace(slug))).First(&city).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("city find by slug: %w", err)
	}
	return &city, nil
}

func (r *Repository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.City{}).Where("slug = ?", slug).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) List(ctx context.Context, limit int) ([]models.City, error) {
	if limit <= 0 {
		limit = 100
	}
	var cities []models.City
	err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&cities).Error
	if err != nil {
		return nil, fmt.Errorf("city list: %w", err)
	}
	return cities, nil
}

func (r *Repository) UpdateArea(ctx context.Context, id uuid.UUID, area int) error {
	res := r.db.WithContext(ctx).Model(&models.City{}).Where("id = ?", id).Update("area_blocks", area)
	if res.Error != nil {
		return fmt.Errorf("city update area: %w", res.Error)
	}
	return nil
}

func (r *Repository) Population(ctx context.Context, cityID uuid.UUID) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Claim{}).
		Where("city_id = ?", cityID).
		Distinct("owner_id").
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("city population: %w", err)
	}
	return int(count), nil
}
