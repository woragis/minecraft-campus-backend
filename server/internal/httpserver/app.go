package httpserver

import (
	alliancesvc "github.com/woragis/minecraft-campus-backend/server/internal/alliance/service"
	citysvc "github.com/woragis/minecraft-campus-backend/server/internal/city/service"
	claimsvc "github.com/woragis/minecraft-campus-backend/server/internal/claim/service"
	guildsvc "github.com/woragis/minecraft-campus-backend/server/internal/guild/service"
	invitesvc "github.com/woragis/minecraft-campus-backend/server/internal/invite/service"
	playersvc "github.com/woragis/minecraft-campus-backend/server/internal/player/service"
	trustsvc "github.com/woragis/minecraft-campus-backend/server/internal/trust/service"
	"gorm.io/gorm"
)

type App struct {
	DB           *gorm.DB
	PluginAPIKey string
	Players      *playersvc.Service
	Invites      *invitesvc.Service
	Guilds       *guildsvc.Service
	Trust        *trustsvc.Service
	Cities       *citysvc.Service
	Claims       *claimsvc.Service
	Alliances    *alliancesvc.Service
}
