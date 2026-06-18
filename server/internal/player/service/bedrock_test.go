package service

import (
	"testing"

	"github.com/google/uuid"
)

func TestBedrockMinecraftUUID_deterministic(t *testing.T) {
	t.Helper()
	a := BedrockMinecraftUUID("2535420712345678")
	b := BedrockMinecraftUUID("2535420712345678")
	if a != b {
		t.Fatalf("expected same uuid, got %s vs %s", a, b)
	}
	if a == uuid.Nil {
		t.Fatal("expected non-nil uuid")
	}
}

func TestBedrockMinecraftUUID_differsByXUID(t *testing.T) {
	t.Helper()
	a := BedrockMinecraftUUID("111")
	b := BedrockMinecraftUUID("222")
	if a == b {
		t.Fatalf("expected different uuids for different xuids")
	}
}
