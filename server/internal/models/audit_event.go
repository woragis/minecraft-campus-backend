package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	AuditEventBlockPlace  = "block_place"
	AuditEventBlockBreak  = "block_break"
	AuditEventPlayerJoin  = "player_join"
	AuditEventPlayerQuit  = "player_quit"
)

type AuditEvent struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	ServerID   uuid.UUID  `gorm:"type:uuid;not null;column:server_id" json:"serverId"`
	World      string     `gorm:"not null" json:"world"`
	PlayerID   uuid.UUID  `gorm:"type:uuid;not null;column:player_id" json:"playerId"`
	EventType  string     `gorm:"not null;column:event_type" json:"eventType"`
	BlockX     *int       `gorm:"column:block_x" json:"blockX,omitempty"`
	BlockY     *int       `gorm:"column:block_y" json:"blockY,omitempty"`
	BlockZ     *int       `gorm:"column:block_z" json:"blockZ,omitempty"`
	BlockType  string     `gorm:"column:block_type" json:"blockType,omitempty"`
	ClaimID    *uuid.UUID `gorm:"type:uuid;column:claim_id" json:"claimId,omitempty"`
	Metadata   string     `gorm:"type:jsonb;not null;default:'{}'" json:"metadata,omitempty"`
	OccurredAt time.Time  `gorm:"not null;column:occurred_at" json:"occurredAt"`
	CreatedAt  time.Time  `json:"createdAt"`
}

func (AuditEvent) TableName() string { return "audit_events" }
