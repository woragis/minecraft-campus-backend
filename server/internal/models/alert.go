package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	AlertTypeGriefingSpike = "griefing_spike"
	AlertSeverityWarning   = "warning"
	AlertSeverityCritical  = "critical"
)

type Alert struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	AlertType    string     `gorm:"not null;column:alert_type" json:"alertType"`
	Severity     string     `gorm:"not null;default:warning" json:"severity"`
	PlayerID     *uuid.UUID `gorm:"type:uuid;column:player_id" json:"playerId,omitempty"`
	Payload      string     `gorm:"type:jsonb;not null;default:'{}'" json:"payload"`
	Acknowledged bool       `gorm:"not null;default:false" json:"acknowledged"`
	CreatedAt    time.Time  `json:"createdAt"`
}

func (Alert) TableName() string { return "alerts" }
