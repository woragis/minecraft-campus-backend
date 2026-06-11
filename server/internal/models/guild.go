package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	GuildRoleLeader = "leader"
	GuildRoleMember = "member"
)

type Guild struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Slug        string    `gorm:"not null" json:"slug"`
	Name        string    `gorm:"not null" json:"name"`
	LeaderID    uuid.UUID `gorm:"type:uuid;not null;column:leader_id" json:"leaderId"`
	TrustScore  int       `gorm:"not null;default:100;column:trust_score" json:"trustScore"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	MemberCount int       `gorm:"-" json:"memberCount,omitempty"`
}

func (Guild) TableName() string { return "guilds" }

type GuildMember struct {
	GuildID  uuid.UUID `gorm:"type:uuid;primaryKey;column:guild_id" json:"guildId"`
	PlayerID uuid.UUID `gorm:"type:uuid;primaryKey;column:player_id" json:"playerId"`
	Role     string    `gorm:"not null;default:member" json:"role"`
	JoinedAt time.Time `gorm:"not null;column:joined_at" json:"joinedAt"`
}

func (GuildMember) TableName() string { return "guild_members" }
