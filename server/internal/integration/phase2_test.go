package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	"gorm.io/gorm"
)

func TestPhase2GuildAndProbationRestrictions(t *testing.T) {
	db, sqlDB := openTestDB(t)
	defer func() { _ = sqlDB.Close() }()

	leaderUUID := uuid.New()
	leader := &models.Player{
		ID:            uuid.New(),
		MinecraftUUID: leaderUUID,
		Username:      "Leader",
		Status:        models.PlayerStatusActive,
		TrustScore:    100,
		SponsorScore:  100,
	}
	if err := db.Create(leader).Error; err != nil {
		t.Fatalf("create leader: %v", err)
	}

	probationUUID := uuid.New()
	until := time.Now().UTC().Add(24 * time.Hour)
	probation := &models.Player{
		ID:             uuid.New(),
		MinecraftUUID:  probationUUID,
		Username:       "Prob",
		Status:         models.PlayerStatusProbation,
		TrustScore:     100,
		SponsorScore:     100,
		ProbationUntil: &until,
	}
	if err := db.Create(probation).Error; err != nil {
		t.Fatalf("create probation player: %v", err)
	}

	server := httptestFromDB(t, db)
	defer server.Close()

	guildID := createGuild(t, server.URL, leaderUUID, "Estatistica")
	if guildID == "" {
		t.Fatal("guild id empty")
	}

	// probation cannot create guild
	status, _ := postJSON(t, server.URL+"/v1/internal/guilds", fmt.Sprintf(`{"leaderUuid":%q,"name":"Fail"}`, probationUUID.String()))
	if status != http.StatusForbidden {
		t.Fatalf("expected 403 for probation guild create, got %d", status)
	}

	// probation cannot invite
	status, _ = postJSON(t, server.URL+"/v1/internal/invites", fmt.Sprintf(`{"sponsorUuid":%q,"targetUsername":"X"}`, probationUUID.String()))
	if status != http.StatusForbidden {
		t.Fatalf("expected 403 for probation invite, got %d", status)
	}

	// active member can join
	memberUUID := uuid.New()
	member := &models.Player{
		ID:            uuid.New(),
		MinecraftUUID: memberUUID,
		Username:      "Member",
		Status:        models.PlayerStatusActive,
		TrustScore:    100,
		SponsorScore:  100,
	}
	if err := db.Create(member).Error; err != nil {
		t.Fatalf("create member: %v", err)
	}
	joinGuild(t, server.URL, guildID, memberUUID)
}

func httptestFromDB(t *testing.T, db *gorm.DB) *httptest.Server {
	t.Helper()
	return httptest.NewServer(newTestHandler(db))
}

func createGuild(t *testing.T, baseURL string, leaderUUID uuid.UUID, name string) string {
	t.Helper()
	status, body := postJSON(t, baseURL+"/v1/internal/guilds", fmt.Sprintf(`{"leaderUuid":%q,"name":%q}`, leaderUUID.String(), name))
	if status != http.StatusCreated {
		t.Fatalf("create guild status %d: %s", status, body)
	}
	var out struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal([]byte(body), &out); err != nil {
		t.Fatalf("decode guild: %v", err)
	}
	return out.ID
}

func joinGuild(t *testing.T, baseURL, guildID string, playerUUID uuid.UUID) {
	t.Helper()
	status, body := postJSON(t, baseURL+"/v1/internal/guilds/"+guildID+"/join", fmt.Sprintf(`{"playerUuid":%q}`, playerUUID.String()))
	if status != http.StatusOK {
		t.Fatalf("join guild status %d: %s", status, body)
	}
}

func postJSON(t *testing.T, url, payload string) (int, string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Plugin-Key", testPluginKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(b)
}
