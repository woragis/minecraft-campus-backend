package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	InviteStatusPending  = "pending"
	InviteStatusAccepted = "accepted"
	InviteStatusExpired  = "expired"
	InviteStatusRevoked  = "revoked"
)

type Invite struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Code             string     `gorm:"not null" json:"code"`
	SponsorID        uuid.UUID  `gorm:"type:uuid;not null;column:sponsor_id" json:"sponsorId"`
	TargetUsername   string     `gorm:"not null;column:target_username" json:"targetUsername"`
	InvitedPlayerID  *uuid.UUID `gorm:"type:uuid;column:invited_player_id" json:"invitedPlayerId,omitempty"`
	Status           string     `gorm:"not null;default:pending" json:"status"`
	CreatedAt        time.Time  `json:"createdAt"`
	AcceptedAt       *time.Time `gorm:"column:accepted_at" json:"acceptedAt,omitempty"`
}

func (Invite) TableName() string { return "invites" }
