package bootstrap

import (
	"testing"

	"github.com/google/uuid"
)

func TestParseFounderFromEnvMissing(t *testing.T) {
	t.Setenv("BOOTSTRAP_MINECRAFT_UUID", "")
	_, ok := ParseFounderFromEnv()
	if ok {
		t.Fatal("expected false when env is empty")
	}
}

func TestParseFounderFromEnvValid(t *testing.T) {
	id := "11111111-1111-1111-1111-111111111111"
	t.Setenv("BOOTSTRAP_MINECRAFT_UUID", id)
	t.Setenv("BOOTSTRAP_USERNAME", "Admin")

	cfg, ok := ParseFounderFromEnv()
	if !ok {
		t.Fatal("expected true")
	}
	if cfg.MinecraftUUID != uuid.MustParse(id) {
		t.Fatalf("uuid mismatch: %s", cfg.MinecraftUUID)
	}
	if cfg.Username != "Admin" {
		t.Fatalf("username: %s", cfg.Username)
	}
}

func TestParseFounderFromEnvInvalidUUID(t *testing.T) {
	t.Setenv("BOOTSTRAP_MINECRAFT_UUID", "not-a-uuid")
	_, ok := ParseFounderFromEnv()
	if ok {
		t.Fatal("expected false for invalid uuid")
	}
}
