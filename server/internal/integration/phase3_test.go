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

func TestPhase3TerritoryClaimsAndCities(t *testing.T) {
	db, sqlDB := openTestDB(t)
	defer func() { _ = sqlDB.Close() }()

	ownerUUID := uuid.New()
	owner := &models.Player{
		ID:            uuid.New(),
		MinecraftUUID: ownerUUID,
		Username:      "Owner",
		Status:        models.PlayerStatusActive,
		TrustScore:    100,
		SponsorScore:  100,
		CreatedAt:     time.Now().UTC().Add(-8 * 24 * time.Hour),
	}
	if err := db.Create(owner).Error; err != nil {
		t.Fatalf("create owner: %v", err)
	}

	probationUUID := uuid.New()
	until := time.Now().UTC().Add(24 * time.Hour)
	probation := &models.Player{
		ID:             uuid.New(),
		MinecraftUUID:  probationUUID,
		Username:       "Prob",
		Status:         models.PlayerStatusProbation,
		TrustScore:     100,
		SponsorScore:   100,
		ProbationUntil: &until,
	}
	if err := db.Create(probation).Error; err != nil {
		t.Fatalf("create probation: %v", err)
	}

	server := httptestFromDB(t, db)
	defer server.Close()

	cityID := createCity(t, server.URL, ownerUUID, "Campus Central")
	if cityID == "" {
		t.Fatal("city id empty")
	}

	claimID := createClaim(t, server.URL, ownerUUID, 0, 30, 0, 30, &cityID)
	if claimID == "" {
		t.Fatal("claim id empty")
	}

	// overlap rejected
	status, _ := postJSON(t, server.URL+"/v1/internal/claims", fmt.Sprintf(
		`{"ownerUuid":%q,"serverSlug":"vanilla","world":"world","minX":10,"maxX":20,"minZ":10,"maxZ":20}`,
		ownerUUID.String(),
	))
	if status != http.StatusConflict {
		t.Fatalf("expected 409 for overlap, got %d", status)
	}

	// probation cannot claim
	status, _ = postJSON(t, server.URL+"/v1/internal/claims", fmt.Sprintf(
		`{"ownerUuid":%q,"serverSlug":"vanilla","world":"world","minX":100,"maxX":110,"minZ":100,"maxZ":110}`,
		probationUUID.String(),
	))
	if status != http.StatusForbidden {
		t.Fatalf("expected 403 for probation claim, got %d", status)
	}

	// owner allowed in claim
	perm := claimPermission(t, server.URL, ownerUUID, 15, 15)
	if !perm.Allowed || perm.Reason != "owner" {
		t.Fatalf("expected owner permission, got %#v", perm)
	}

	// stranger denied in claim
	strangerUUID := uuid.New()
	perm = claimPermission(t, server.URL, strangerUUID, 15, 15)
	if perm.Allowed {
		t.Fatalf("expected denied for stranger, got %#v", perm)
	}

	// wilderness allowed
	perm = claimPermission(t, server.URL, strangerUUID, 500, 500)
	if !perm.Allowed || perm.Reason != "wilderness" {
		t.Fatalf("expected wilderness, got %#v", perm)
	}

	guildAID := createGuild(t, server.URL, ownerUUID, "Alpha")
	guildBLeader := uuid.New()
	guildBLeaderPlayer := &models.Player{
		ID:            uuid.New(),
		MinecraftUUID: guildBLeader,
		Username:      "BetaLead",
		Status:        models.PlayerStatusActive,
		TrustScore:    100,
		SponsorScore:  100,
	}
	if err := db.Create(guildBLeaderPlayer).Error; err != nil {
		t.Fatalf("create guild b leader: %v", err)
	}
	guildBID := createGuild(t, server.URL, guildBLeader, "Beta")

	status, _ = postJSON(t, server.URL+"/v1/internal/alliances", fmt.Sprintf(
		`{"leaderUuid":%q,"guildAId":%q,"guildBId":%q}`,
		ownerUUID.String(), guildAID, guildBID,
	))
	if status != http.StatusCreated {
		t.Fatalf("expected 201 for alliance, got %d", status)
	}
}

type claimPermResult struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason"`
	ClaimID string `json:"claimId,omitempty"`
}

func createCity(t *testing.T, baseURL string, founderUUID uuid.UUID, name string) string {
	t.Helper()
	status, body := postJSON(t, baseURL+"/v1/internal/cities", fmt.Sprintf(
		`{"founderUuid":%q,"name":%q,"serverSlug":"vanilla","world":"world","centerX":15,"centerZ":15}`,
		founderUUID.String(), name,
	))
	if status != http.StatusCreated {
		t.Fatalf("create city status %d: %s", status, body)
	}
	var out struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal([]byte(body), &out); err != nil {
		t.Fatalf("decode city: %v", err)
	}
	return out.ID
}

func createClaim(t *testing.T, baseURL string, ownerUUID uuid.UUID, minX, maxX, minZ, maxZ int, cityID *string) string {
	t.Helper()
	payload := fmt.Sprintf(
		`{"ownerUuid":%q,"serverSlug":"vanilla","world":"world","minX":%d,"maxX":%d,"minZ":%d,"maxZ":%d,"zoneType":"urban"`,
		ownerUUID.String(), minX, maxX, minZ, maxZ,
	)
	if cityID != nil {
		payload += fmt.Sprintf(`,"cityId":%q`, *cityID)
	}
	payload += "}"
	status, body := postJSON(t, baseURL+"/v1/internal/claims", payload)
	if status != http.StatusCreated {
		t.Fatalf("create claim status %d: %s", status, body)
	}
	var out struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal([]byte(body), &out); err != nil {
		t.Fatalf("decode claim: %v", err)
	}
	return out.ID
}

func claimPermission(t *testing.T, baseURL string, mcUUID uuid.UUID, x, z int) claimPermResult {
	t.Helper()
	url := fmt.Sprintf("%s/v1/internal/claims/permission?minecraftUuid=%s&serverSlug=vanilla&world=world&x=%d&z=%d",
		baseURL, mcUUID.String(), x, z)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("X-Plugin-Key", testPluginKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get permission: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("permission status %d", resp.StatusCode)
	}
	var out claimPermResult
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode permission: %v", err)
	}
	return out
}
