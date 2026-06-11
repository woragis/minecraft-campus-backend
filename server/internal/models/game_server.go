package models

import (
	"time"

	"github.com/google/uuid"
)

type GameServer struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Slug      string    `gorm:"not null" json:"slug"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

func (GameServer) TableName() string { return "game_servers" }

type ServerPlayer struct {
	ServerID     uuid.UUID `gorm:"type:uuid;primaryKey;column:server_id" json:"serverId"`
	PlayerID     uuid.UUID `gorm:"type:uuid;primaryKey;column:player_id" json:"playerId"`
	LastSeenAt   time.Time `gorm:"not null;column:last_seen_at" json:"lastSeenAt"`
	PlayTimeSecs int64     `gorm:"not null;default:0;column:play_time_secs" json:"playTimeSecs"`
}

func (ServerPlayer) TableName() string { return "server_players" }
