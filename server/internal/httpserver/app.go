package httpserver

import (
	invitesvc "github.com/woragis/minecraft-campus-backend/server/internal/invite/service"
	playersvc "github.com/woragis/minecraft-campus-backend/server/internal/player/service"
	"gorm.io/gorm"
)

type App struct {
	DB            *gorm.DB
	PluginAPIKey  string
	Players       *playersvc.Service
	Invites       *invitesvc.Service
}
