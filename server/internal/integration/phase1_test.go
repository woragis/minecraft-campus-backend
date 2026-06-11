package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	alliancerepo "github.com/woragis/minecraft-campus-backend/server/internal/alliance/repository"
	alliancesvc "github.com/woragis/minecraft-campus-backend/server/internal/alliance/service"
	alertsrepo "github.com/woragis/minecraft-campus-backend/server/internal/alerts/repository"
	alertssvc "github.com/woragis/minecraft-campus-backend/server/internal/alerts/service"
	auditrepo "github.com/woragis/minecraft-campus-backend/server/internal/audit/repository"
	auditsvc "github.com/woragis/minecraft-campus-backend/server/internal/audit/service"
	"github.com/woragis/minecraft-campus-backend/server/internal/config"
	cityrepo "github.com/woragis/minecraft-campus-backend/server/internal/city/repository"
	citysvc "github.com/woragis/minecraft-campus-backend/server/internal/city/service"
	claimrepo "github.com/woragis/minecraft-campus-backend/server/internal/claim/repository"
	claimsvc "github.com/woragis/minecraft-campus-backend/server/internal/claim/service"
	gameserverrepo "github.com/woragis/minecraft-campus-backend/server/internal/gameserver/repository"
	guildrepo "github.com/woragis/minecraft-campus-backend/server/internal/guild/repository"
	guildsvc "github.com/woragis/minecraft-campus-backend/server/internal/guild/service"
	"github.com/woragis/minecraft-campus-backend/server/internal/httpserver"
	inviterepo "github.com/woragis/minecraft-campus-backend/server/internal/invite/repository"
	invitesvc "github.com/woragis/minecraft-campus-backend/server/internal/invite/service"
	metricssvc "github.com/woragis/minecraft-campus-backend/server/internal/metrics/service"
	rollbackrepo "github.com/woragis/minecraft-campus-backend/server/internal/rollback/repository"
	rollbacksvc "github.com/woragis/minecraft-campus-backend/server/internal/rollback/service"
	trustrepo "github.com/woragis/minecraft-campus-backend/server/internal/trust/repository"
	trustsvc "github.com/woragis/minecraft-campus-backend/server/internal/trust/service"
	"github.com/woragis/minecraft-campus-backend/server/internal/middleware"
	"github.com/woragis/minecraft-campus-backend/server/internal/migrate"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	playerrepo "github.com/woragis/minecraft-campus-backend/server/internal/player/repository"
	playersvc "github.com/woragis/minecraft-campus-backend/server/internal/player/service"
	"github.com/woragis/minecraft-campus-backend/server/internal/platform/postgres"
	"gorm.io/gorm"
)

const testPluginKey = "integration-test-key"

type whitelistResult struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason"`
}

type inviteResult struct {
	Code string `json:"code"`
}

func TestPhase1InviteAndWhitelistFlow(t *testing.T) {
	db, sqlDB := openTestDB(t)
	defer func() { _ = sqlDB.Close() }()

	sponsorUUID := uuid.New()
	sponsor := &models.Player{
		ID:            uuid.New(),
		MinecraftUUID: sponsorUUID,
		Username:      "Sponsor",
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

	inviteCode := createInvite(t, server.URL, sponsorUUID, "NewPlayer")
	if inviteCode == "" {
		t.Fatal("invite code empty")
	}

	inviteeUUID := uuid.New()
	result := checkWhitelist(t, server.URL, inviteeUUID, "NewPlayer")
	if !result.Allowed {
		t.Fatalf("expected allowed after invite, got %#v", result)
	}
	if result.Reason != "probation" {
		t.Fatalf("expected probation, got %s", result.Reason)
	}

	upsertPlayer(t, server.URL, inviteeUUID, "NewPlayer")
}

func openTestDB(t *testing.T) (*gorm.DB, *sql.DB) {
	t.Helper()
	dsn := strings.TrimSpace(os.Getenv("TEST_DATABASE_URL"))
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set — skipping integration test")
	}

	db, err := postgres.Open(dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql db: %v", err)
	}

	dir := migrateDir(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err := migrate.Up(ctx, sqlDB, dir); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db, sqlDB
}

func migrateDir(t *testing.T) string {
	t.Helper()
	candidates := []string{
		filepath.Join("..", "..", "..", "migrations"),
		filepath.Join("..", "..", "migrations"),
		"migrations",
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			return c
		}
	}
	t.Fatal("migrations directory not found")
	return ""
}

func newTestHandler(db *gorm.DB) http.Handler {
	playerRepository := playerrepo.New(db)
	inviteRepository := inviterepo.New(db)
	gameServerRepository := gameserverrepo.New(db)
	guildRepository := guildrepo.New(db)
	trustRepository := trustrepo.New(db)
	cityRepository := cityrepo.New(db)
	claimRepository := claimrepo.New(db)
	allianceRepository := alliancerepo.New(db)
	auditRepository := auditrepo.New(db)
	rollbackRepository := rollbackrepo.New(db)
	trustService := trustsvc.New(trustRepository, playerRepository)
	playerService := playersvc.New(playerRepository, inviteRepository, gameServerRepository, 7, trustService)
	inviteService := invitesvc.New(inviteRepository, playerRepository)
	guildService := guildsvc.New(guildRepository, playerRepository)
	cityService := citysvc.New(cityRepository, playerRepository, guildRepository, gameServerRepository)
	claimService := claimsvc.New(claimRepository, playerRepository, cityRepository, guildRepository, gameServerRepository)
	allianceService := alliancesvc.New(allianceRepository, guildRepository, playerRepository)
	testCfg := config.Load()
	auditService := auditsvc.New(testCfg, auditRepository, playerRepository, gameServerRepository)
	rollbackService := rollbacksvc.New(rollbackRepository, auditRepository, playerRepository, gameServerRepository, trustService)
	metricsService := metricssvc.New(db)
	alertsService := alertssvc.New(testCfg, alertsrepo.New(db))
	app := &httpserver.App{
		DB:           db,
		PluginAPIKey: testPluginKey,
		Config:       testCfg,
		Players:      playerService,
		Invites:      inviteService,
		Guilds:       guildService,
		Trust:        trustService,
		Cities:       cityService,
		Claims:       claimService,
		Alliances:    allianceService,
		Audit:        auditService,
		Rollback:     rollbackService,
		Metrics:      metricsService,
		Alerts:       alertsService,
	}
	return httpserver.NewHandler(app, middleware.Config{})
}

func createInvite(t *testing.T, baseURL string, sponsorUUID uuid.UUID, target string) string {
	t.Helper()
	body := fmt.Sprintf(`{"sponsorUuid":%q,"targetUsername":%q}`, sponsorUUID.String(), target)
	req, err := http.NewRequest(http.MethodPost, baseURL+"/v1/internal/invites", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Plugin-Key", testPluginKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("post invite: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("invite status %d: %s", resp.StatusCode, string(b))
	}
	var out inviteResult
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode invite: %v", err)
	}
	return out.Code
}

func checkWhitelist(t *testing.T, baseURL string, mcUUID uuid.UUID, username string) whitelistResult {
	t.Helper()
	reqURL := fmt.Sprintf("%s/v1/internal/whitelist/%s?username=%s", baseURL, mcUUID.String(), url.QueryEscape(username))
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("X-Plugin-Key", testPluginKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get whitelist: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("whitelist status %d: %s", resp.StatusCode, string(b))
	}
	var out whitelistResult
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode whitelist: %v", err)
	}
	return out
}

func upsertPlayer(t *testing.T, baseURL string, mcUUID uuid.UUID, username string) {
	t.Helper()
	body := fmt.Sprintf(`{"minecraftUuid":%q,"username":%q,"serverSlug":"vanilla"}`, mcUUID.String(), username)
	req, err := http.NewRequest(http.MethodPost, baseURL+"/v1/internal/players/upsert", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Plugin-Key", testPluginKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("post upsert: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("upsert status %d: %s", resp.StatusCode, string(b))
	}
}
