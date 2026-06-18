package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	playersvc "github.com/woragis/minecraft-campus-backend/server/internal/player/service"
)

func TestBedrockInviteAndWhitelistFlow(t *testing.T) {
	db, sqlDB := openTestDB(t)
	defer func() { _ = sqlDB.Close() }()

	sponsorUUID := uuid.New()
	sponsor := &models.Player{
		ID:            uuid.New(),
		MinecraftUUID: sponsorUUID,
		Username:      "BedrockSponsor",
		Status:        models.PlayerStatusActive,
		TrustScore:    100,
		SponsorScore:  100,
	}
	if err := db.Create(sponsor).Error; err != nil {
		t.Fatalf("create sponsor: %v", err)
	}

	handler := newTestHandler(db)
	server := httptest.NewServer(handler)
	defer server.Close()

	xuid := "2535420712345678"
	inviteCode := createInvite(t, server.URL, sponsorUUID, "BedrockNew")
	if inviteCode == "" {
		t.Fatal("invite code empty")
	}

	result := checkBedrockWhitelist(t, server.URL, xuid, "BedrockNew")
	if !result.Allowed {
		t.Fatalf("expected allowed after invite, got %#v", result)
	}
	if result.Reason != "probation" {
		t.Fatalf("expected probation, got %s", result.Reason)
	}

	upsertBedrockPlayer(t, server.URL, xuid, "BedrockNew")

	expectedUUID := playersvc.BedrockMinecraftUUID(xuid)
	var player models.Player
	if err := db.Where("minecraft_uuid = ?", expectedUUID).First(&player).Error; err != nil {
		t.Fatalf("load bedrock player: %v", err)
	}
	if player.Username != "BedrockNew" {
		t.Fatalf("expected username BedrockNew, got %s", player.Username)
	}

	var identity models.PlayerIdentity
	if err := db.Where("platform = ? AND external_id = ?", models.PlatformBedrock, xuid).First(&identity).Error; err != nil {
		t.Fatalf("load bedrock identity: %v", err)
	}
	if identity.PlayerID != player.ID {
		t.Fatalf("identity player mismatch")
	}
}

func checkBedrockWhitelist(t *testing.T, baseURL, xuid, username string) whitelistResult {
	t.Helper()
	reqURL := fmt.Sprintf("%s/v1/internal/whitelist/bedrock/%s?username=%s", baseURL, xuid, username)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("X-Plugin-Key", testPluginKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get bedrock whitelist: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("bedrock whitelist status %d: %s", resp.StatusCode, string(b))
	}
	var out whitelistResult
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode bedrock whitelist: %v", err)
	}
	return out
}

func upsertBedrockPlayer(t *testing.T, baseURL, xuid, username string) {
	t.Helper()
	body := fmt.Sprintf(`{"xuid":%q,"username":%q,"serverSlug":"bedrock"}`, xuid, username)
	req, err := http.NewRequest(http.MethodPost, baseURL+"/v1/internal/players/bedrock/upsert", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Plugin-Key", testPluginKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("post bedrock upsert: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("bedrock upsert status %d: %s", resp.StatusCode, string(b))
	}
}
