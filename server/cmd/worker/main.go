package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	auditrepo "github.com/woragis/minecraft-campus-backend/server/internal/audit/repository"
	auditsvc "github.com/woragis/minecraft-campus-backend/server/internal/audit/service"
	alertsrepo "github.com/woragis/minecraft-campus-backend/server/internal/alerts/repository"
	alertssvc "github.com/woragis/minecraft-campus-backend/server/internal/alerts/service"
	backuprepo "github.com/woragis/minecraft-campus-backend/server/internal/backup/repository"
	backupsvc "github.com/woragis/minecraft-campus-backend/server/internal/backup/service"
	"github.com/woragis/minecraft-campus-backend/server/internal/config"
	metricssvc "github.com/woragis/minecraft-campus-backend/server/internal/metrics/service"
	"github.com/woragis/minecraft-campus-backend/server/internal/migrate"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	"github.com/woragis/minecraft-campus-backend/server/internal/platform/postgres"
	"github.com/woragis/minecraft-campus-backend/server/internal/worker"
)

func main() {
	cfg := config.Load()
	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if !cfg.WorkerEnabled {
		log.Print("WORKER_ENABLED=0 — worker exits successfully without running jobs")
		return
	}

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
		}
	}

	_ = db.AutoMigrate(
		&models.WorldSnapshot{},
		&models.Alert{},
	)

	backupService := backupsvc.New(cfg, backuprepo.New(db), dsn)
	metricsService := metricssvc.New(db)
	alertsService := alertssvc.New(cfg, alertsrepo.New(db))
	auditService := auditsvc.New(cfg, auditrepo.New(db), nil, nil)

	runner := worker.New(cfg, backupService, metricsService, alertsService, auditService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
	}()

	runner.Run(ctx)
}
