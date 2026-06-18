package config

import (
	"os"
	"strconv"
	"strings"
)

type Profile string

const (
	ProfileDev        Profile = "dev"
	ProfileBudget     Profile = "budget"
	ProfileProduction Profile = "production"
)

type Config struct {
	Profile Profile

	WorkerEnabled bool

	BackupEnabled          bool
	BackupDatabaseEnabled  bool
	BackupWorldEnabled     bool
	BackupStorage          string
	BackupLocalPath        string
	BackupDatabaseRetentionDays int
	BackupDailyEnabled     bool
	BackupWeeklyEnabled    bool
	BackupMonthlyEnabled   bool

	MetricsRefreshEnabled       bool
	MetricsRefreshIntervalMin   int
	AlertsEnabled               bool
	AlertsIntervalMin           int
	AlertsDiscordWebhook        string
	AlertsGriefingThresholdHour int

	AuditIngestEnabled    bool
	AuditRetentionDays    int
	AuditPurgeEnabled     bool
	AuditPurgeIntervalHrs int

	RedisEnabled bool
	RedisURL     string
	PresenceTTLSeconds int
}

func Load() Config {
	profile := Profile(strings.ToLower(strings.TrimSpace(envOr("CAMPUSWORLD_PROFILE", string(ProfileBudget)))))
	cfg := defaultsForProfile(profile)
	cfg.Profile = profile

	cfg.WorkerEnabled = envBool("WORKER_ENABLED", cfg.WorkerEnabled)

	cfg.BackupEnabled = envBool("BACKUP_ENABLED", cfg.BackupEnabled)
	cfg.BackupDatabaseEnabled = envBool("BACKUP_DATABASE_ENABLED", cfg.BackupDatabaseEnabled)
	cfg.BackupWorldEnabled = envBool("BACKUP_WORLD_ENABLED", cfg.BackupWorldEnabled)
	if v := strings.TrimSpace(os.Getenv("BACKUP_STORAGE")); v != "" {
		cfg.BackupStorage = strings.ToLower(v)
	}
	cfg.BackupLocalPath = envOr("BACKUP_LOCAL_PATH", cfg.BackupLocalPath)
	cfg.BackupDatabaseRetentionDays = envInt("BACKUP_DATABASE_RETENTION_DAYS", cfg.BackupDatabaseRetentionDays)
	cfg.BackupDailyEnabled = envBool("BACKUP_DAILY_ENABLED", cfg.BackupDailyEnabled)
	cfg.BackupWeeklyEnabled = envBool("BACKUP_WEEKLY_ENABLED", cfg.BackupWeeklyEnabled)
	cfg.BackupMonthlyEnabled = envBool("BACKUP_MONTHLY_ENABLED", cfg.BackupMonthlyEnabled)

	cfg.MetricsRefreshEnabled = envBool("METRICS_REFRESH_ENABLED", cfg.MetricsRefreshEnabled)
	cfg.MetricsRefreshIntervalMin = envInt("METRICS_REFRESH_INTERVAL_MINUTES", cfg.MetricsRefreshIntervalMin)
	cfg.AlertsEnabled = envBool("ALERTS_ENABLED", cfg.AlertsEnabled)
	cfg.AlertsIntervalMin = envInt("ALERTS_INTERVAL_MINUTES", cfg.AlertsIntervalMin)
	cfg.AlertsDiscordWebhook = envOr("ALERTS_DISCORD_WEBHOOK", "")
	cfg.AlertsGriefingThresholdHour = envInt("ALERTS_GRIEFING_THRESHOLD_HOUR", cfg.AlertsGriefingThresholdHour)

	cfg.AuditIngestEnabled = envBool("AUDIT_INGEST_ENABLED", cfg.AuditIngestEnabled)
	cfg.AuditRetentionDays = envInt("AUDIT_RETENTION_DAYS", cfg.AuditRetentionDays)
	cfg.AuditPurgeEnabled = envBool("AUDIT_PURGE_ENABLED", cfg.AuditPurgeEnabled)
	cfg.AuditPurgeIntervalHrs = envInt("AUDIT_PURGE_INTERVAL_HOURS", cfg.AuditPurgeIntervalHrs)

	cfg.RedisEnabled = envBool("REDIS_ENABLED", cfg.RedisEnabled)
	cfg.RedisURL = envOr("REDIS_URL", cfg.RedisURL)
	cfg.PresenceTTLSeconds = envInt("PRESENCE_TTL_SECONDS", cfg.PresenceTTLSeconds)

	return cfg
}

func defaultsForProfile(p Profile) Config {
	switch p {
	case ProfileDev:
		return Config{
			WorkerEnabled:               false,
			BackupEnabled:               false,
			BackupDatabaseEnabled:       false,
			BackupWorldEnabled:          false,
			BackupStorage:               "none",
			BackupLocalPath:             "./backups",
			BackupDatabaseRetentionDays: 3,
			MetricsRefreshEnabled:       false,
			MetricsRefreshIntervalMin:   60,
			AlertsEnabled:               false,
			AlertsIntervalMin:           15,
			AlertsGriefingThresholdHour: 200,
			AuditIngestEnabled:          true,
			AuditRetentionDays:          14,
			AuditPurgeEnabled:           true,
			AuditPurgeIntervalHrs:       24,
			RedisEnabled:                false,
			RedisURL:                    "redis://127.0.0.1:6379/0",
			PresenceTTLSeconds:          120,
		}
	case ProfileProduction:
		return Config{
			WorkerEnabled:               true,
			BackupEnabled:               true,
			BackupDatabaseEnabled:       true,
			BackupWorldEnabled:          false,
			BackupStorage:               "local",
			BackupLocalPath:             "/var/backups/campusworld",
			BackupDatabaseRetentionDays: 30,
			BackupDailyEnabled:          true,
			BackupWeeklyEnabled:         true,
			BackupMonthlyEnabled:        true,
			MetricsRefreshEnabled:       true,
			MetricsRefreshIntervalMin:   60,
			AlertsEnabled:               true,
			AlertsIntervalMin:           15,
			AlertsGriefingThresholdHour: 100,
			AuditIngestEnabled:          true,
			AuditRetentionDays:          90,
			AuditPurgeEnabled:           true,
			AuditPurgeIntervalHrs:       24,
			RedisEnabled:                true,
			RedisURL:                    "redis://redis:6379/0",
			PresenceTTLSeconds:          120,
		}
	default: // budget
		return Config{
			WorkerEnabled:               false,
			BackupEnabled:               false,
			BackupDatabaseEnabled:       false,
			BackupWorldEnabled:          false,
			BackupStorage:               "none",
			BackupLocalPath:             "/var/backups/campusworld",
			BackupDatabaseRetentionDays: 7,
			BackupDailyEnabled:          false,
			BackupWeeklyEnabled:         false,
			BackupMonthlyEnabled:        false,
			MetricsRefreshEnabled:       false,
			MetricsRefreshIntervalMin:   60,
			AlertsEnabled:               false,
			AlertsIntervalMin:           15,
			AlertsGriefingThresholdHour: 150,
			AuditIngestEnabled:          true,
			AuditRetentionDays:          30,
			AuditPurgeEnabled:           true,
			AuditPurgeIntervalHrs:       24,
			RedisEnabled:                false,
			RedisURL:                    "redis://127.0.0.1:6379/0",
			PresenceTTLSeconds:          120,
		}
	}
}

func (c Config) BackupActive() bool {
	return c.BackupEnabled && c.BackupStorage != "none"
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
	if err != nil {
		return def
	}
	return n
}

func envBool(key string, def bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	switch strings.ToLower(v) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return def
	}
}
