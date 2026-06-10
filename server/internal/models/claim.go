package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	ZoneUrban      = "urban"
	ZoneRural      = "rural"
	ZoneIndustrial = "industrial"
	ZoneHistoric   = "historic"
)

type Claim struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	OwnerID    uuid.UUID  `gorm:"type:uuid;not null;column:owner_id" json:"ownerId"`
	CityID     *uuid.UUID `gorm:"type:uuid;column:city_id" json:"cityId,omitempty"`
	GuildID    *uuid.UUID `gorm:"type:uuid;column:guild_id" json:"guildId,omitempty"`
	ServerID   uuid.UUID  `gorm:"type:uuid;not null;column:server_id" json:"serverId"`
	World      string     `gorm:"not null;default:world" json:"world"`
	MinX       int        `gorm:"not null;column:min_x" json:"minX"`
	MinY       int        `gorm:"not null;default:-64;column:min_y" json:"minY"`
	MaxX       int        `gorm:"not null;column:max_x" json:"maxX"`
	MaxY       int        `gorm:"not null;default:320;column:max_y" json:"maxY"`
	MinZ       int        `gorm:"not null;column:min_z" json:"minZ"`
	MaxZ       int        `gorm:"not null;column:max_z" json:"maxZ"`
	ZoneType   string     `gorm:"not null;default:urban;column:zone_type" json:"zoneType"`
	AreaBlocks int        `gorm:"not null;column:area_blocks" json:"areaBlocks"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

func (Claim) TableName() string { return "claims" }
