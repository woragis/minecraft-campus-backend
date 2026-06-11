package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/config"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
)

func TestPhase5MetricsAndAlerts(t *testing.T) {
	db, sqlDB := openTestDB(t)
	defer func() { _ = sqlDB.Close() }()

	playerID := uuid.New()
	player := &models.Player{
		ID:            playerID,
		MinecraftUUID: uuid.New(),
		Username:      "MetricsPlayer",
		Status:        models.PlayerStatusActive,
		TrustScore:    100,
		SponsorScore:  100,
	}
	if err := db.Create(player).Error; err != nil {
		t.Fatalf("create player: %v", err)
	}

	gs := &models.GameServer{ID: uuid.New(), Slug: "vanilla", Name: "vanilla"}
	if err := db.Create(gs).Error; err != nil {
		t.Fatalf("create server: %v", err)
	}

	claimID := uuid.New()
	if err := db.Create(&models.Claim{
		ID: claimID, OwnerID: playerID, ServerID: gs.ID, World: "world",
		MinX: 0, MaxX: 10, MinY: -64, MaxY: 320, MinZ: 0, MaxZ: 10,
		ZoneType: models.ZoneUrban, AreaBlocks: 121,
	}).Error; err != nil {
		t.Fatalf("create claim: %v", err)
	}

	now := time.Now().UTC()
	for i := 0; i < 5; i++ {
		x, y, z := 10+i, 64, 10+i
		if err := db.Create(&models.AuditEvent{
			ID: uuid.New(), ServerID: gs.ID, World: "world", PlayerID: playerID,
			EventType: models.AuditEventBlockBreak, BlockX: &x, BlockY: &y, BlockZ: &z,
			BlockType: "STONE", ClaimID: &claimID, Metadata: "{}", OccurredAt: now,
		}).Error; err != nil {
			t.Fatalf("audit event: %v", err)
		}
	}

	server := httptestFromDB(t, db)
	defer server.Close()

	overview := getMetricsOverview(t, server.URL)
	if overview.TotalPlayers < 1 || overview.TotalClaims < 1 {
		t.Fatalf("unexpected overview: %#v", overview)
	}
	if overview.AuditBreaks24h < 5 {
		t.Fatalf("expected breaks in 24h, got %d", overview.AuditBreaks24h)
	}

	if err := db.Create(&models.Alert{
		ID: uuid.New(), AlertType: models.AlertTypeGriefingSpike,
		Severity: models.AlertSeverityWarning, PlayerID: &playerID, Payload: "{}",
	}).Error; err != nil {
		t.Fatalf("create alert: %v", err)
	}

	alerts := listAlerts(t, server.URL)
	if len(alerts) < 1 {
		t.Fatal("expected alerts")
	}
}

type metricsOverview struct {
	TotalPlayers   int64 `json:"totalPlayers"`
	TotalClaims    int64 `json:"totalClaims"`
	AuditBreaks24h int64 `json:"auditBreaks24h"`
}

func getMetricsOverview(t *testing.T, baseURL string) metricsOverview {
	t.Helper()
	resp, err := http.Get(baseURL + "/v1/metrics/overview")
	if err != nil {
		t.Fatalf("get overview: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("overview status %d", resp.StatusCode)
	}
	var out metricsOverview
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return out
}

func listAlerts(t *testing.T, baseURL string) []map[string]any {
	t.Helper()
	resp, err := http.Get(baseURL + "/v1/alerts")
	if err != nil {
		t.Fatalf("get alerts: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("alerts status %d", resp.StatusCode)
	}
	var out struct {
		Alerts []map[string]any `json:"alerts"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return out.Alerts
}

func TestPhase5BackupSkippedWhenDisabled(t *testing.T) {
	t.Setenv("CAMPUSWORLD_PROFILE", "budget")
	t.Setenv("BACKUP_ENABLED", "0")
	cfg := config.Load()
	if cfg.BackupActive() {
		t.Fatal("budget should not run backups")
	}
}
