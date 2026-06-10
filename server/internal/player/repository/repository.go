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

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.Player, error) {
	var p models.Player
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("player find by id: %w", err)
	}
	return &p, nil
}

func (r *Repository) FindByMinecraftUUID(ctx context.Context, minecraftUUID uuid.UUID) (*models.Player, error) {
	var p models.Player
	err := r.db.WithContext(ctx).Where("minecraft_uuid = ?", minecraftUUID).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("player find by minecraft uuid: %w", err)
	}
	return &p, nil
}

func (r *Repository) Create(ctx context.Context, p *models.Player) error {
	if err := r.db.WithContext(ctx).Create(p).Error; err != nil {
		return fmt.Errorf("player create: %w", err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, p *models.Player) error {
	if err := r.db.WithContext(ctx).Save(p).Error; err != nil {
		return fmt.Errorf("player update: %w", err)
	}
	return nil
}

func (r *Repository) ListInvitedBy(ctx context.Context, sponsorID uuid.UUID) ([]models.Player, error) {
	var players []models.Player
	err := r.db.WithContext(ctx).
		Where("invited_by_id = ?", sponsorID).
		Order("created_at ASC").
		Find(&players).Error
	if err != nil {
		return nil, fmt.Errorf("player list invited by: %w", err)
	}
	return players, nil
}

func (r *Repository) FindByUsernameInsensitive(ctx context.Context, username string) (*models.Player, error) {
	var p models.Player
	err := r.db.WithContext(ctx).
		Where("lower(username) = ?", strings.ToLower(strings.TrimSpace(username))).
		First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("player find by username: %w", err)
	}
	return &p, nil
}
