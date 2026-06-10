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
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, inv *models.Invite) error {
	if err := r.db.WithContext(ctx).Create(inv).Error; err != nil {
		return fmt.Errorf("invite create: %w", err)
	}
	return nil
}

func (r *Repository) FindByCode(ctx context.Context, code string) (*models.Invite, error) {
	var inv models.Invite
	err := r.db.WithContext(ctx).Where("code = ?", strings.TrimSpace(code)).First(&inv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("invite find by code: %w", err)
	}
	return &inv, nil
}

func (r *Repository) FindPendingByTargetUsername(ctx context.Context, username string) (*models.Invite, error) {
	var inv models.Invite
	err := r.db.WithContext(ctx).
		Where("lower(target_username) = ? AND status = ?", strings.ToLower(strings.TrimSpace(username)), models.InviteStatusPending).
		Order("created_at ASC").
		First(&inv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("invite find pending by username: %w", err)
	}
	return &inv, nil
}

func (r *Repository) HasPendingForTargetUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Invite{}).
		Where("lower(target_username) = ? AND status = ?", strings.ToLower(strings.TrimSpace(username)), models.InviteStatusPending).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("invite pending exists: %w", err)
	}
	return count > 0, nil
}

func (r *Repository) Accept(ctx context.Context, inviteID, playerID uuid.UUID, acceptedAt time.Time) error {
	res := r.db.WithContext(ctx).Model(&models.Invite{}).
		Where("id = ? AND status = ?", inviteID, models.InviteStatusPending).
		Updates(map[string]any{
			"status":            models.InviteStatusAccepted,
			"invited_player_id": playerID,
			"accepted_at":       acceptedAt,
		})
	if res.Error != nil {
		return fmt.Errorf("invite accept: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) ListBySponsor(ctx context.Context, sponsorID uuid.UUID) ([]models.Invite, error) {
	var invites []models.Invite
	err := r.db.WithContext(ctx).
		Where("sponsor_id = ?", sponsorID).
		Order("created_at DESC").
		Find(&invites).Error
	if err != nil {
		return nil, fmt.Errorf("invite list by sponsor: %w", err)
	}
	return invites, nil
}
