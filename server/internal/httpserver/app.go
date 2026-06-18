package httpserver

import (
	alliancesvc "github.com/woragis/minecraft-campus-backend/server/internal/alliance/service"
	alertssvc "github.com/woragis/minecraft-campus-backend/server/internal/alerts/service"
	auditsvc "github.com/woragis/minecraft-campus-backend/server/internal/audit/service"
	citysvc "github.com/woragis/minecraft-campus-backend/server/internal/city/service"
	claimsvc "github.com/woragis/minecraft-campus-backend/server/internal/claim/service"
	"github.com/woragis/minecraft-campus-backend/server/internal/config"
	guildsvc "github.com/woragis/minecraft-campus-backend/server/internal/guild/service"
	invitesvc "github.com/woragis/minecraft-campus-backend/server/internal/invite/service"
	playersvc "github.com/woragis/minecraft-campus-backend/server/internal/player/service"
	presencesvc "github.com/woragis/minecraft-campus-backend/server/internal/presence"
	metricssvc "github.com/woragis/minecraft-campus-backend/server/internal/metrics/service"
	rollbacksvc "github.com/woragis/minecraft-campus-backend/server/internal/rollback/service"
	trustsvc "github.com/woragis/minecraft-campus-backend/server/internal/trust/service"
	"gorm.io/gorm"
)

type App struct {
	DB           *gorm.DB
	PluginAPIKey string
	Config       config.Config
	Players      *playersvc.Service
	Presence     *presencesvc.Service
	Invites      *invitesvc.Service
	Guilds       *guildsvc.Service
	Trust        *trustsvc.Service
	Cities       *citysvc.Service
	Claims       *claimsvc.Service
	Alliances    *alliancesvc.Service
	Audit        *auditsvc.Service
	Rollback     *rollbacksvc.Service
	Metrics      *metricssvc.Service
	Alerts       *alertssvc.Service
}
