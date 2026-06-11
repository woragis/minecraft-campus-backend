package worker

import (
	"context"
	"log"
	"time"

	alertssvc "github.com/woragis/minecraft-campus-backend/server/internal/alerts/service"
	auditsvc "github.com/woragis/minecraft-campus-backend/server/internal/audit/service"
	backupsvc "github.com/woragis/minecraft-campus-backend/server/internal/backup/service"
	"github.com/woragis/minecraft-campus-backend/server/internal/config"
	metricssvc "github.com/woragis/minecraft-campus-backend/server/internal/metrics/service"
)

type Runner struct {
	cfg     config.Config
	backup  *backupsvc.Service
	metrics *metricssvc.Service
	alerts  *alertssvc.Service
	audit   *auditsvc.Service
}

func New(cfg config.Config, backup *backupsvc.Service, metrics *metricssvc.Service, alerts *alertssvc.Service, audit *auditsvc.Service) *Runner {
	return &Runner{cfg: cfg, backup: backup, metrics: metrics, alerts: alerts, audit: audit}
}

func (r *Runner) Run(ctx context.Context) {
	if !r.cfg.WorkerEnabled {
		log.Print("worker: WORKER_ENABLED=0 — exiting (no background jobs)")
		return
	}
	log.Printf("worker: started (profile=%s backup=%v alerts=%v metrics_refresh=%v)",
		r.cfg.Profile, r.cfg.BackupActive(), r.cfg.AlertsEnabled, r.cfg.MetricsRefreshEnabled)

	backupTicker := time.NewTicker(24 * time.Hour)
	alertTicker := time.NewTicker(time.Duration(r.cfg.AlertsIntervalMin) * time.Minute)
	metricsTicker := time.NewTicker(time.Duration(r.cfg.MetricsRefreshIntervalMin) * time.Minute)
	purgeTicker := time.NewTicker(time.Duration(r.cfg.AuditPurgeIntervalHrs) * time.Hour)
	defer backupTicker.Stop()
	defer alertTicker.Stop()
	defer metricsTicker.Stop()
	defer purgeTicker.Stop()

	r.tick(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Print("worker: shutdown")
			return
		case <-backupTicker.C:
			r.runBackup(ctx)
		case <-alertTicker.C:
			r.runAlerts(ctx)
		case <-metricsTicker.C:
			r.runMetrics(ctx)
		case <-purgeTicker.C:
			r.runPurge(ctx)
		}
	}
}

func (r *Runner) tick(ctx context.Context) {
	r.runBackup(ctx)
	r.runAlerts(ctx)
	r.runMetrics(ctx)
	r.runPurge(ctx)
}

func (r *Runner) runBackup(ctx context.Context) {
	if !r.cfg.BackupActive() {
		return
	}
	res, err := r.backup.RunDatabaseBackup(ctx)
	if err != nil {
		log.Printf("worker: backup error: %v", err)
		return
	}
	if res != nil && res.Skipped {
		log.Printf("worker: backup skipped: %s", res.Reason)
	} else if res != nil && res.Snapshot != nil {
		log.Printf("worker: backup ok path=%s size=%d", res.Snapshot.Path, res.Snapshot.SizeBytes)
	}
}

func (r *Runner) runAlerts(ctx context.Context) {
	if !r.cfg.AlertsEnabled {
		return
	}
	n, err := r.alerts.ScanGriefing(ctx)
	if err != nil {
		log.Printf("worker: alerts error: %v", err)
		return
	}
	if n > 0 {
		log.Printf("worker: created %d griefing alerts", n)
	}
}

func (r *Runner) runMetrics(ctx context.Context) {
	if !r.cfg.MetricsRefreshEnabled {
		return
	}
	if err := r.metrics.Refresh(ctx); err != nil {
		log.Printf("worker: metrics refresh: %v", err)
	}
}

func (r *Runner) runPurge(ctx context.Context) {
	if _, err := r.audit.PurgeOld(ctx); err != nil {
		log.Printf("worker: audit purge: %v", err)
	}
}
