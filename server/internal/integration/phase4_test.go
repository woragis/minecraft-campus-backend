package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
)

func TestPhase4AuditIngestAndRollback(t *testing.T) {
	db, sqlDB := openTestDB(t)
	defer func() { _ = sqlDB.Close() }()

	grieferUUID := uuid.New()
	griefer := &models.Player{
		ID:            uuid.New(),
		MinecraftUUID: grieferUUID,
		Username:      "Griefer",
		Status:        models.PlayerStatusActive,
		TrustScore:    100,
		SponsorScore:  100,
	}
	if err := db.Create(griefer).Error; err != nil {
		t.Fatalf("create griefer: %v", err)
	}

	modUUID := uuid.New()
	mod := &models.Player{
		ID:            uuid.New(),
		MinecraftUUID: modUUID,
		Username:      "Mod",
		Status:        models.PlayerStatusActive,
		TrustScore:    100,
		SponsorScore:  100,
	}
	if err := db.Create(mod).Error; err != nil {
		t.Fatalf("create mod: %v", err)
	}

	claimID := uuid.New()
	gs := &models.GameServer{ID: uuid.New(), Slug: "vanilla", Name: "vanilla"}
	if err := db.Create(gs).Error; err != nil {
		t.Fatalf("create server: %v", err)
	}
	claim := &models.Claim{
		ID:         claimID,
		OwnerID:    griefer.ID,
		ServerID:   gs.ID,
		World:      "world",
		MinX:       0,
		MaxX:       50,
		MinY:       -64,
		MaxY:       320,
		MinZ:       0,
		MaxZ:       50,
		ZoneType:   models.ZoneUrban,
		AreaBlocks: 2601,
	}
	if err := db.Create(claim).Error; err != nil {
		t.Fatalf("create claim: %v", err)
	}

	server := httptestFromDB(t, db)
	defer server.Close()

	now := time.Now().UTC()
	ingestAuditEvents(t, server.URL, fmt.Sprintf(`{
		"events":[
			{"minecraftUuid":%q,"serverSlug":"vanilla","world":"world","eventType":"block_break","blockX":10,"blockY":64,"blockZ":10,"blockType":"STONE","claimId":%q,"occurredAt":%q},
			{"minecraftUuid":%q,"serverSlug":"vanilla","world":"world","eventType":"block_place","blockX":11,"blockY":64,"blockZ":11,"blockType":"DIRT","claimId":%q,"occurredAt":%q}
		]
	}`, grieferUUID.String(), claimID.String(), now.Format(time.RFC3339), grieferUUID.String(), claimID.String(), now.Add(time.Minute).Format(time.RFC3339)))

	status, body := postJSON(t, server.URL+"/v1/internal/rollbacks", fmt.Sprintf(
		`{"targetUuid":%q,"actorUuid":%q,"serverSlug":"vanilla","world":"world","windowMinutes":60}`,
		grieferUUID.String(), modUUID.String(),
	))
	if status != http.StatusCreated {
		t.Fatalf("create rollback status %d: %s", status, body)
	}

	var rollbackOut struct {
		Rollback struct {
			ID        string `json:"id"`
			ItemCount int    `json:"itemCount"`
		} `json:"rollback"`
	}
	if err := json.Unmarshal([]byte(body), &rollbackOut); err != nil {
		t.Fatalf("decode rollback: %v", err)
	}
	if rollbackOut.Rollback.ItemCount != 2 {
		t.Fatalf("expected 2 rollback items, got %d", rollbackOut.Rollback.ItemCount)
	}

	items := listRollbackItems(t, server.URL, rollbackOut.Rollback.ID)
	if len(items) != 2 {
		t.Fatalf("expected 2 items from list, got %d", len(items))
	}

	completeRollback(t, server.URL, rollbackOut.Rollback.ID, 2)

	var grieferAfter models.Player
	if err := db.Where("id = ?", griefer.ID).First(&grieferAfter).Error; err != nil {
		t.Fatalf("reload griefer: %v", err)
	}
	if grieferAfter.TrustScore != 95 {
		t.Fatalf("expected trust 95 after rollback_applied, got %d", grieferAfter.TrustScore)
	}
}

func ingestAuditEvents(t *testing.T, baseURL, payload string) {
	t.Helper()
	status, body := postJSON(t, baseURL+"/v1/internal/audit/events", payload)
	if status != http.StatusCreated {
		t.Fatalf("ingest audit status %d: %s", status, body)
	}
}

func listRollbackItems(t *testing.T, baseURL, rollbackID string) []map[string]any {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, baseURL+"/v1/internal/rollbacks/"+rollbackID+"/items", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("X-Plugin-Key", testPluginKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get items: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("items status %d", resp.StatusCode)
	}
	var out struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode items: %v", err)
	}
	return out.Items
}

func completeRollback(t *testing.T, baseURL, rollbackID string, applied int) {
	t.Helper()
	status, body := postJSON(t, baseURL+"/v1/internal/rollbacks/"+rollbackID+"/complete", fmt.Sprintf(`{"appliedCount":%d}`, applied))
	if status != http.StatusOK {
		t.Fatalf("complete rollback status %d: %s", status, body)
	}
}
