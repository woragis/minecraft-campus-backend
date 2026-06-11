package config

import (
	"os"
	"testing"
)

func TestLoadBudgetDefaults(t *testing.T) {
	os.Unsetenv("CAMPUSWORLD_PROFILE")
	os.Unsetenv("BACKUP_ENABLED")
	cfg := Load()
	if cfg.Profile != ProfileBudget {
		t.Fatalf("expected budget profile, got %s", cfg.Profile)
	}
	if cfg.BackupEnabled {
		t.Fatal("budget should disable backups by default")
	}
	if cfg.BackupStorage != "none" {
		t.Fatalf("expected storage none, got %s", cfg.BackupStorage)
	}
}

func TestLoadProductionOverride(t *testing.T) {
	t.Setenv("CAMPUSWORLD_PROFILE", "production")
	t.Setenv("BACKUP_WORLD_ENABLED", "0")
	cfg := Load()
	if !cfg.BackupEnabled {
		t.Fatal("production should enable backups")
	}
	if cfg.BackupWorldEnabled {
		t.Fatal("world backup should stay off when overridden")
	}
}
