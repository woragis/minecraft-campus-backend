package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	RollbackStatusPending   = "pending"
	RollbackStatusApplying  = "applying"
	RollbackStatusCompleted = "completed"
	RollbackStatusFailed    = "failed"

	RollbackActionRestore = "restore"
	RollbackActionRemove  = "remove"
)

type Rollback struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	TargetPlayerID uuid.UUID  `gorm:"type:uuid;not null;column:target_player_id" json:"targetPlayerId"`
	ActorPlayerID  *uuid.UUID `gorm:"type:uuid;column:actor_player_id" json:"actorPlayerId,omitempty"`
	ServerID       uuid.UUID  `gorm:"type:uuid;not null;column:server_id" json:"serverId"`
	World          string     `gorm:"not null" json:"world"`
	FromAt         time.Time  `gorm:"not null;column:from_at" json:"fromAt"`
	ToAt           time.Time  `gorm:"not null;column:to_at" json:"toAt"`
	Status         string     `gorm:"not null;default:pending" json:"status"`
	ItemCount      int        `gorm:"not null;default:0;column:item_count" json:"itemCount"`
	AppliedCount   int        `gorm:"not null;default:0;column:applied_count" json:"appliedCount"`
	CreatedAt      time.Time  `json:"createdAt"`
	CompletedAt    *time.Time `gorm:"column:completed_at" json:"completedAt,omitempty"`
}

func (Rollback) TableName() string { return "rollbacks" }

type RollbackItem struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	RollbackID   uuid.UUID  `gorm:"type:uuid;not null;column:rollback_id" json:"rollbackId"`
	AuditEventID *uuid.UUID `gorm:"type:uuid;column:audit_event_id" json:"auditEventId,omitempty"`
	Action       string     `gorm:"not null" json:"action"`
	BlockType    string     `gorm:"not null;column:block_type" json:"blockType"`
	BlockX       int        `gorm:"not null;column:block_x" json:"blockX"`
	BlockY       int        `gorm:"not null;column:block_y" json:"blockY"`
	BlockZ       int        `gorm:"not null;column:block_z" json:"blockZ"`
	World        string     `gorm:"not null" json:"world"`
	Applied      bool       `gorm:"not null;default:false" json:"applied"`
	CreatedAt    time.Time  `json:"createdAt"`
}

func (RollbackItem) TableName() string { return "rollback_items" }
