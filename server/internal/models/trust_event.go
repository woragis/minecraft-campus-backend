package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	TrustEventConfirmedReport  = "confirmed_report"
	TrustEventRollbackApplied  = "rollback_applied"
	TrustEventProbationClean   = "probation_day_clean"
	TrustEventBan              = "ban"
	TrustEventUnban            = "unban"
	TrustEventManualAdjust     = "manual_adjust"
)

type TrustEvent struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	PlayerID      uuid.UUID  `gorm:"type:uuid;not null;column:player_id" json:"playerId"`
	EventType     string     `gorm:"not null;column:event_type" json:"eventType"`
	Delta         int        `gorm:"not null;default:0" json:"delta"`
	Reason        string     `json:"reason,omitempty"`
	ActorPlayerID *uuid.UUID `gorm:"type:uuid;column:actor_player_id" json:"actorPlayerId,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
}

func (TrustEvent) TableName() string { return "trust_events" }
