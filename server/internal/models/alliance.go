package models

import (
	"time"

	"github.com/google/uuid"
)

const AllianceStatusActive = "active"

type Alliance struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	GuildAID  uuid.UUID `gorm:"type:uuid;not null;column:guild_a_id" json:"guildAId"`
	GuildBID  uuid.UUID `gorm:"type:uuid;not null;column:guild_b_id" json:"guildBId"`
	Status    string    `gorm:"not null;default:active" json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func (Alliance) TableName() string { return "alliances" }
