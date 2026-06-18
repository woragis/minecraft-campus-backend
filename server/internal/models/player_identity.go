package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	PlatformJava   = "java"
	PlatformBedrock = "bedrock"
)

type PlayerIdentity struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PlayerID   uuid.UUID `gorm:"type:uuid;not null;column:player_id" json:"playerId"`
	Platform   string    `gorm:"not null" json:"platform"`
	ExternalID string    `gorm:"not null;column:external_id" json:"externalId"`
	Username   string    `gorm:"not null" json:"username"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (PlayerIdentity) TableName() string { return "player_identities" }
