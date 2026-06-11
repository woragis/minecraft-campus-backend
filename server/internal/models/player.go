package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	PlayerStatusActive    = "active"
	PlayerStatusProbation = "probation"
	PlayerStatusBanned    = "banned"
)

type Player struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	MinecraftUUID  uuid.UUID  `gorm:"type:uuid;not null;column:minecraft_uuid" json:"minecraftUuid"`
	Username       string     `gorm:"not null" json:"username"`
	Status         string     `gorm:"not null;default:probation" json:"status"`
	InvitedByID    *uuid.UUID `gorm:"type:uuid;column:invited_by_id" json:"invitedById,omitempty"`
	TrustScore     int        `gorm:"not null;default:100;column:trust_score" json:"trustScore"`
	SponsorScore   int        `gorm:"not null;default:100;column:sponsor_score" json:"sponsorScore"`
	ProbationUntil *time.Time `gorm:"column:probation_until" json:"probationUntil,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

func (Player) TableName() string { return "players" }
