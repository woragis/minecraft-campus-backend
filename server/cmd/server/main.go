package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	alliancerepo "github.com/woragis/minecraft-campus-backend/server/internal/alliance/repository"
	alliancesvc "github.com/woragis/minecraft-campus-backend/server/internal/alliance/service"
	alertsrepo "github.com/woragis/minecraft-campus-backend/server/internal/alerts/repository"
	alertssvc "github.com/woragis/minecraft-campus-backend/server/internal/alerts/service"
	auditrepo "github.com/woragis/minecraft-campus-backend/server/internal/audit/repository"
	auditsvc "github.com/woragis/minecraft-campus-backend/server/internal/audit/service"
	"github.com/woragis/minecraft-campus-backend/server/internal/bootstrap"
	"github.com/woragis/minecraft-campus-backend/server/internal/config"
	cityrepo "github.com/woragis/minecraft-campus-backend/server/internal/city/repository"
	citysvc "github.com/woragis/minecraft-campus-backend/server/internal/city/service"
	claimrepo "github.com/woragis/minecraft-campus-backend/server/internal/claim/repository"
	claimsvc "github.com/woragis/minecraft-campus-backend/server/internal/claim/service"
	gameserverrepo "github.com/woragis/minecraft-campus-backend/server/internal/gameserver/repository"
	guildrepo "github.com/woragis/minecraft-campus-backend/server/internal/guild/repository"
	guildsvc "github.com/woragis/minecraft-campus-backend/server/internal/guild/service"
	"github.com/woragis/minecraft-campus-backend/server/internal/httpserver"
	metricssvc "github.com/woragis/minecraft-campus-backend/server/internal/metrics/service"
	inviterepo "github.com/woragis/minecraft-campus-backend/server/internal/invite/repository"
	invitesvc "github.com/woragis/minecraft-campus-backend/server/internal/invite/service"
	rollbackrepo "github.com/woragis/minecraft-campus-backend/server/internal/rollback/repository"
	rollbacksvc "github.com/woragis/minecraft-campus-backend/server/internal/rollback/service"
	trustrepo "github.com/woragis/minecraft-campus-backend/server/internal/trust/repository"
	trustsvc "github.com/woragis/minecraft-campus-backend/server/internal/trust/service"
	"github.com/woragis/minecraft-campus-backend/server/internal/middleware"
	"github.com/woragis/minecraft-campus-backend/server/internal/migrate"
	"github.com/woragis/minecraft-campus-backend/server/internal/presence"
	redisplatform "github.com/woragis/minecraft-campus-backend/server/internal/platform/redis"
	playerrepo "github.com/woragis/minecraft-campus-backend/server/internal/player/repository"
	playersvc "github.com/woragis/minecraft-campus-backend/server/internal/player/service"
	"github.com/woragis/minecraft-campus-backend/server/internal/platform/postgres"
)

func main() {
	appCfg := config.Load()
	log.Printf("campusworld profile=%s backup=%v worker_hint=%v", appCfg.Profile, appCfg.BackupActive(), appCfg.WorkerEnabled)

	addr := envOr("HTTP_ADDR", ":8080")
	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}
	pluginAPIKey := strings.TrimSpace(os.Getenv("PLUGIN_API_KEY"))
	if pluginAPIKey == "" {
		log.Fatal("PLUGIN_API_KEY is required")
	}
	probationDays := envInt("PROBATION_DAYS", 7)

	db, err := postgres.Open(dsn)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	if skip := strings.TrimSpace(os.Getenv("SKIP_SQL_MIGRATIONS")); skip != "1" && !strings.EqualFold(skip, "true") {
		dir := strings.TrimSpace(os.Getenv("MIGRATIONS_DIR"))
		if dir == "" {
			dir = migrate.ResolveDir()
		}
		if dir != "" {
			sqlDB, err := db.DB()
			if err != nil {
				log.Fatalf("sql db: %v", err)
			}
			mctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			err = migrate.Up(mctx, sqlDB, dir)
			cancel()
			if err != nil {
				log.Fatalf("sql migrate: %v", err)
			}
			log.Printf("sql migrations applied from %s", dir)
		} else {
			log.Print("warning: SQL migrations skipped (no migrations/ directory found)")
		}
	}

	// Schema is owned by SQL migrations in /migrations (GORM AutoMigrate conflicts with PG constraint names).

	playerRepository := playerrepo.New(db)
	inviteRepository := inviterepo.New(db)
	gameServerRepository := gameserverrepo.New(db)
	guildRepository := guildrepo.New(db)
	trustRepository := trustrepo.New(db)
	cityRepository := cityrepo.New(db)
	claimRepository := claimrepo.New(db)
	allianceRepository := alliancerepo.New(db)
	auditRepository := auditrepo.New(db)
	rollbackRepository := rollbackrepo.New(db)

	if cfg, ok := bootstrap.ParseFounderFromEnv(); ok {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		created, err := bootstrap.EnsureFounder(ctx, playerRepository, cfg)
		cancel()
		if err != nil {
			log.Fatalf("bootstrap founder: %v", err)
		}
		if created {
			log.Printf("bootstrap founder created: %s (%s)", cfg.Username, cfg.MinecraftUUID)
		}
	}

	trustService := trustsvc.New(trustRepository, playerRepository)
	playerService := playersvc.New(playerRepository, inviteRepository, gameServerRepository, probationDays, trustService)
	inviteService := invitesvc.New(inviteRepository, playerRepository)
	guildService := guildsvc.New(guildRepository, playerRepository)
	cityService := citysvc.New(cityRepository, playerRepository, guildRepository, gameServerRepository)
	claimService := claimsvc.New(claimRepository, playerRepository, cityRepository, guildRepository, gameServerRepository)
	allianceService := alliancesvc.New(allianceRepository, guildRepository, playerRepository)
	auditService := auditsvc.New(appCfg, auditRepository, playerRepository, gameServerRepository)
	rollbackService := rollbacksvc.New(rollbackRepository, auditRepository, playerRepository, gameServerRepository, trustService)
	metricsService := metricssvc.New(db)
	alertsService := alertssvc.New(appCfg, alertsrepo.New(db))

	var presenceStore presence.Store = presence.NewNoopStore()
	if appCfg.RedisEnabled {
		redisClient, err := redisplatform.Open(appCfg.RedisURL)
		if err != nil {
			log.Fatalf("redis: %v", err)
		}
		defer func() { _ = redisClient.Close() }()
		presenceStore = presence.NewRedisStore(redisClient, appCfg.PresenceTTLSeconds)
		log.Printf("redis presence enabled (ttl=%ds)", appCfg.PresenceTTLSeconds)
	}
	presenceService := presence.New(presenceStore, guildRepository, playerRepository)

	app := &httpserver.App{
		DB:           db,
		PluginAPIKey: pluginAPIKey,
		Config:       appCfg,
		Players:      playerService,
		Presence:     presenceService,
		Invites:      inviteService,
		Guilds:       guildService,
		Trust:        trustService,
		Cities:       cityService,
		Claims:       claimService,
		Alliances:    allianceService,
		Audit:        auditService,
		Rollback:     rollbackService,
		Metrics:      metricsService,
		Alerts:       alertsService,
	}

	handler := httpserver.NewHandler(app, middleware.LoadConfigFromEnv())

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("api listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func envOr(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return def
	}
	return n
}
