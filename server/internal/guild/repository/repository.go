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

func (r *Repository) Create(ctx context.Context, guild *models.Guild, leaderMember *models.GuildMember) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(guild).Error; err != nil {
			return fmt.Errorf("guild create: %w", err)
		}
		if err := tx.Create(leaderMember).Error; err != nil {
			return fmt.Errorf("guild leader member: %w", err)
		}
		return nil
	})
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.Guild, error) {
	var g models.Guild
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&g).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("guild find by id: %w", err)
	}
	return &g, nil
}

func (r *Repository) FindBySlug(ctx context.Context, slug string) (*models.Guild, error) {
	var g models.Guild
	err := r.db.WithContext(ctx).Where("slug = ?", strings.ToLower(strings.TrimSpace(slug))).First(&g).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("guild find by slug: %w", err)
	}
	return &g, nil
}

func (r *Repository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Guild{}).Where("slug = ?", slug).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("guild slug exists: %w", err)
	}
	return count > 0, nil
}

func (r *Repository) List(ctx context.Context, limit int) ([]models.Guild, error) {
	if limit <= 0 {
		limit = 50
	}
	var guilds []models.Guild
	err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&guilds).Error
	if err != nil {
		return nil, fmt.Errorf("guild list: %w", err)
	}
	return guilds, nil
}

func (r *Repository) Update(ctx context.Context, guild *models.Guild) error {
	if err := r.db.WithContext(ctx).Save(guild).Error; err != nil {
		return fmt.Errorf("guild update: %w", err)
	}
	return nil
}

func (r *Repository) AddMember(ctx context.Context, member *models.GuildMember) error {
	if err := r.db.WithContext(ctx).Create(member).Error; err != nil {
		return fmt.Errorf("guild add member: %w", err)
	}
	return nil
}

func (r *Repository) RemoveMember(ctx context.Context, guildID, playerID uuid.UUID) error {
	res := r.db.WithContext(ctx).
		Where("guild_id = ? AND player_id = ?", guildID, playerID).
		Delete(&models.GuildMember{})
	if res.Error != nil {
		return fmt.Errorf("guild remove member: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) FindMember(ctx context.Context, guildID, playerID uuid.UUID) (*models.GuildMember, error) {
	var m models.GuildMember
	err := r.db.WithContext(ctx).
		Where("guild_id = ? AND player_id = ?", guildID, playerID).
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("guild find member: %w", err)
	}
	return &m, nil
}

func (r *Repository) ListMembers(ctx context.Context, guildID uuid.UUID) ([]models.GuildMember, error) {
	var members []models.GuildMember
	err := r.db.WithContext(ctx).
		Where("guild_id = ?", guildID).
		Order("joined_at ASC").
		Find(&members).Error
	if err != nil {
		return nil, fmt.Errorf("guild list members: %w", err)
	}
	return members, nil
}

func (r *Repository) PlayerGuild(ctx context.Context, playerID uuid.UUID) (*models.Guild, error) {
	var member models.GuildMember
	err := r.db.WithContext(ctx).Where("player_id = ?", playerID).First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("player guild lookup: %w", err)
	}
	return r.FindByID(ctx, member.GuildID)
}

func (r *Repository) MemberCount(ctx context.Context, guildID uuid.UUID) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.GuildMember{}).Where("guild_id = ?", guildID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("guild member count: %w", err)
	}
	return int(count), nil
}

func (r *Repository) AverageMemberTrust(ctx context.Context, guildID uuid.UUID) (int, error) {
	type row struct{ Avg float64 }
	var result row
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(AVG(p.trust_score), 100) AS avg
		FROM guild_members gm
		JOIN players p ON p.id = gm.player_id
		WHERE gm.guild_id = ?
	`, guildID).Scan(&result).Error
	if err != nil {
		return 100, fmt.Errorf("guild avg trust: %w", err)
	}
	return int(result.Avg), nil
}
