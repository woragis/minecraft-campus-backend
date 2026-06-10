package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	playerrepo "github.com/woragis/minecraft-campus-backend/server/internal/player/repository"
	"gorm.io/gorm"
)

type FounderConfig struct {
	MinecraftUUID uuid.UUID
	Username      string
}

func ParseFounderFromEnv() (FounderConfig, bool) {
	rawUUID := strings.TrimSpace(os.Getenv("BOOTSTRAP_MINECRAFT_UUID"))
	if rawUUID == "" {
		return FounderConfig{}, false
	}
	mcUUID, err := uuid.Parse(rawUUID)
	if err != nil || mcUUID == uuid.Nil {
		return FounderConfig{}, false
	}
	username := strings.TrimSpace(os.Getenv("BOOTSTRAP_USERNAME"))
	if username == "" {
		username = "Fundador"
	}
	return FounderConfig{MinecraftUUID: mcUUID, Username: username}, true
}

// EnsureFounder creates the bootstrap player when configured and missing.
// Returns true when a new founder row was inserted.
func EnsureFounder(ctx context.Context, repo *playerrepo.Repository, cfg FounderConfig) (bool, error) {
	_, err := repo.FindByMinecraftUUID(ctx, cfg.MinecraftUUID)
	if err == nil {
		return false, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("bootstrap founder lookup: %w", err)
	}

	player := &models.Player{
		ID:            uuid.New(),
		MinecraftUUID: cfg.MinecraftUUID,
		Username:      cfg.Username,
		Status:        models.PlayerStatusActive,
		TrustScore:    100,
		SponsorScore:  100,
	}
	if err := repo.Create(ctx, player); err != nil {
		return false, fmt.Errorf("bootstrap founder create: %w", err)
	}
	return true, nil
}
