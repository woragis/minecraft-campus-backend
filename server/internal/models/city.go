package models

import (
	"time"

	"github.com/google/uuid"
)

type City struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Slug       string     `gorm:"uniqueIndex;not null" json:"slug"`
	Name       string     `gorm:"not null" json:"name"`
	FounderID  uuid.UUID  `gorm:"type:uuid;not null;column:founder_id" json:"founderId"`
	GuildID    *uuid.UUID `gorm:"type:uuid;column:guild_id" json:"guildId,omitempty"`
	ServerID   uuid.UUID  `gorm:"type:uuid;not null;column:server_id" json:"serverId"`
	World      string     `gorm:"not null;default:world" json:"world"`
	CenterX    int        `gorm:"not null;default:0;column:center_x" json:"centerX"`
	CenterZ    int        `gorm:"not null;default:0;column:center_z" json:"centerZ"`
	AreaBlocks int        `gorm:"not null;default:0;column:area_blocks" json:"areaBlocks"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	Population int        `gorm:"-" json:"population,omitempty"`
}

func (City) TableName() string { return "cities" }
