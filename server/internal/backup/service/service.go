package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	backuprepo "github.com/woragis/minecraft-campus-backend/server/internal/backup/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/config"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
)

type Service struct {
	cfg  config.Config
	repo *backuprepo.Repository
	dsn  string
}

func New(cfg config.Config, repo *backuprepo.Repository, dsn string) *Service {
	return &Service{cfg: cfg, repo: repo, dsn: dsn}
}

type RunResult struct {
	Skipped bool              `json:"skipped"`
	Reason  string            `json:"reason,omitempty"`
	Snapshot *models.WorldSnapshot `json:"snapshot,omitempty"`
}

func (s *Service) RunDatabaseBackup(ctx context.Context) (*RunResult, error) {
	if !s.cfg.BackupEnabled {
		return &RunResult{Skipped: true, Reason: "BACKUP_ENABLED=0"}, nil
	}
	if !s.cfg.BackupDatabaseEnabled {
		return &RunResult{Skipped: true, Reason: "BACKUP_DATABASE_ENABLED=0"}, nil
	}
	if s.cfg.BackupStorage == "none" {
		return &RunResult{Skipped: true, Reason: "BACKUP_STORAGE=none"}, nil
	}
	if s.cfg.BackupStorage == "s3" || s.cfg.BackupStorage == "b2" {
		return s.recordSkipped(ctx, "cloud storage not implemented; use BACKUP_STORAGE=local")
	}

	if err := os.MkdirAll(s.cfg.BackupLocalPath, 0o750); err != nil {
		return nil, fmt.Errorf("mkdir backup path: %w", err)
	}

	filename := fmt.Sprintf("pg-%s.sql.gz", time.Now().UTC().Format("20060102-150405"))
	fullPath := filepath.Join(s.cfg.BackupLocalPath, filename)

	if err := s.pgDumpGzip(ctx, fullPath); err != nil {
		snap := &models.WorldSnapshot{
			ID:           uuid.New(),
			SnapshotType: snapshotTypeForNow(s.cfg),
			Storage:      models.SnapshotStorageLocal,
			Path:         fullPath,
			Status:       models.SnapshotStatusFailed,
		}
		_ = s.repo.Create(ctx, snap)
		return nil, err
	}

	size, checksum, err := fileMeta(fullPath)
	if err != nil {
		return nil, err
	}

	snap := &models.WorldSnapshot{
		ID:           uuid.New(),
		SnapshotType: snapshotTypeForNow(s.cfg),
		Storage:      models.SnapshotStorageLocal,
		Path:         fullPath,
		SizeBytes:    size,
		Checksum:     checksum,
		Status:       models.SnapshotStatusCompleted,
	}
	if err := s.repo.Create(ctx, snap); err != nil {
		return nil, err
	}

	if s.cfg.BackupWorldEnabled {
		log.Printf("backup: world backup skipped (BACKUP_WORLD_ENABLED requires manual setup; disabled by default in budget)")
	}

	if err := s.purgeOld(ctx); err != nil {
		log.Printf("backup: purge warning: %v", err)
	}

	return &RunResult{Snapshot: snap}, nil
}

func (s *Service) RunWorldBackup(ctx context.Context) (*RunResult, error) {
	if !s.cfg.BackupEnabled || !s.cfg.BackupWorldEnabled {
		return &RunResult{Skipped: true, Reason: "world backup disabled"}, nil
	}
	return s.recordSkipped(ctx, "world backup not implemented; use external tools or enable later")
}

func (s *Service) recordSkipped(ctx context.Context, reason string) (*RunResult, error) {
	snap := &models.WorldSnapshot{
		ID:           uuid.New(),
		SnapshotType: models.SnapshotTypeManual,
		Storage:      models.SnapshotStorageNone,
		Path:         reason,
		Status:       models.SnapshotStatusSkipped,
	}
	_ = s.repo.Create(ctx, snap)
	return &RunResult{Skipped: true, Reason: reason, Snapshot: snap}, nil
}

func (s *Service) purgeOld(ctx context.Context) error {
	if s.cfg.BackupDatabaseRetentionDays <= 0 {
		return nil
	}
	cutoff := time.Now().UTC().Add(-time.Duration(s.cfg.BackupDatabaseRetentionDays) * 24 * time.Hour)
	snaps, err := s.repo.ListRecent(ctx, 500)
	if err != nil {
		return err
	}
	for _, snap := range snaps {
		if snap.CreatedAt.Before(cutoff) && snap.Storage == models.SnapshotStorageLocal && snap.Path != "" {
			_ = os.Remove(snap.Path)
		}
	}
	_, _ = s.repo.DeleteOlderThan(ctx, cutoff)
	return nil
}

func (s *Service) pgDumpGzip(ctx context.Context, dest string) error {
	pgDump, err := exec.LookPath("pg_dump")
	if err != nil {
		return fmt.Errorf("pg_dump not found in PATH: %w", err)
	}
	host, port, user, password, dbname, err := parseDSN(s.dsn)
	if err != nil {
		return err
	}

	tmpSQL := dest + ".tmp"
	args := []string{"-h", host, "-p", port, "-U", user, "-d", dbname, "-f", tmpSQL, "--no-owner", "--no-acl"}
	cmd := exec.CommandContext(ctx, pgDump, args...)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+password)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("pg_dump: %w (%s)", err, strings.TrimSpace(string(out)))
	}

	gzipPath, err := exec.LookPath("gzip")
	if err != nil {
		_ = os.Rename(tmpSQL, dest)
		return nil
	}
	gzipCmd := exec.CommandContext(ctx, gzipPath, "-f", tmpSQL)
	if out, err := gzipCmd.CombinedOutput(); err != nil {
		_ = os.Remove(tmpSQL)
		return fmt.Errorf("gzip: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return os.Rename(tmpSQL+".gz", dest)
}

func parseDSN(dsn string) (host, port, user, password, dbname string, err error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("parse DATABASE_URL: %w", err)
	}
	host = u.Hostname()
	port = u.Port()
	if port == "" {
		port = "5432"
	}
	if u.User != nil {
		user = u.User.Username()
		password, _ = u.User.Password()
	}
	dbname = strings.TrimPrefix(u.Path, "/")
	if dbname == "" {
		return "", "", "", "", "", fmt.Errorf("database name missing in DATABASE_URL")
	}
	return host, port, user, password, dbname, nil
}

func fileMeta(path string) (int64, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, "", err
	}
	defer f.Close()
	h := sha256.New()
	size, err := io.Copy(h, f)
	if err != nil {
		return 0, "", err
	}
	return size, hex.EncodeToString(h.Sum(nil)), nil
}

func snapshotTypeForNow(cfg config.Config) string {
	now := time.Now().UTC()
	if cfg.BackupMonthlyEnabled && now.Day() == 1 {
		return models.SnapshotTypeMonthly
	}
	if cfg.BackupWeeklyEnabled && now.Weekday() == time.Sunday {
		return models.SnapshotTypeWeekly
	}
	if cfg.BackupDailyEnabled {
		return models.SnapshotTypeDaily
	}
	return models.SnapshotTypeManual
}
