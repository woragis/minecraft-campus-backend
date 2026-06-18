package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) FindIdentityByPlatformExternalID(ctx context.Context, platform, externalID string) (*models.PlayerIdentity, error) {
	platform = strings.TrimSpace(platform)
	externalID = strings.TrimSpace(externalID)
	var identity models.PlayerIdentity
	err := r.db.WithContext(ctx).
		Where("platform = ? AND external_id = ?", platform, externalID).
		First(&identity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("player identity find: %w", err)
	}
	return &identity, nil
}

func (r *Repository) UpsertIdentity(ctx context.Context, identity *models.PlayerIdentity) error {
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "platform"}, {Name: "external_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"username", "updated_at"}),
	}).Create(identity).Error; err != nil {
		return fmt.Errorf("player identity upsert: %w", err)
	}
	return nil
}

func (r *Repository) FindPlayerIDByIdentity(ctx context.Context, platform, externalID string) (uuid.UUID, error) {
	identity, err := r.FindIdentityByPlatformExternalID(ctx, platform, externalID)
	if err != nil {
		return uuid.Nil, err
	}
	return identity.PlayerID, nil
}
