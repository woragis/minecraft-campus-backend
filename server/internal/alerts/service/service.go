package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/alerts/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	"github.com/woragis/minecraft-campus-backend/server/internal/config"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	cfg  config.Config
	repo *repository.Repository
}

func New(cfg config.Config, repo *repository.Repository) *Service {
	return &Service{cfg: cfg, repo: repo}
}

func (s *Service) ScanGriefing(ctx context.Context) (int, error) {
	if !s.cfg.AlertsEnabled {
		return 0, nil
	}
	since := time.Now().UTC().Add(-time.Hour)
	threshold := int64(s.cfg.AlertsGriefingThresholdHour)
	if threshold <= 0 {
		threshold = 150
	}
	rows, err := s.repo.GriefingSpikes(ctx, since, threshold)
	if err != nil {
		return 0, err
	}
	created := 0
	for _, row := range rows {
		dup, err := s.repo.RecentDuplicate(ctx, models.AlertTypeGriefingSpike, row.PlayerID, since)
		if err != nil || dup {
			continue
		}
		payload, _ := json.Marshal(map[string]any{
			"breakCount": row.Count,
			"window":     "1h",
			"threshold":  threshold,
		})
		pid := row.PlayerID
		alert := &models.Alert{
			ID:        uuid.New(),
			AlertType: models.AlertTypeGriefingSpike,
			Severity:  models.AlertSeverityWarning,
			PlayerID:  &pid,
			Payload:   string(payload),
		}
		if err := s.repo.Create(ctx, alert); err != nil {
			return created, err
		}
		created++
		if s.cfg.AlertsDiscordWebhook != "" {
			go s.notifyDiscord(alert, row.Count)
		}
	}
	return created, nil
}

func (s *Service) ListUnacknowledged(ctx context.Context) ([]models.Alert, error) {
	alerts, err := s.repo.ListUnacknowledged(ctx, 50)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeAlertsListV1ServiceFailed, apperrors.MsgAlertsListV1ServiceFailed, err)
	}
	return alerts, nil
}

func (s *Service) Acknowledge(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Acknowledge(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeAlertsAckV1ServiceNotFound, apperrors.MsgAlertsAckV1ServiceNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeAlertsAckV1ServiceFailed, apperrors.MsgAlertsAckV1ServiceFailed, err)
	}
	return nil
}

func (s *Service) notifyDiscord(alert *models.Alert, count int64) {
	body := fmt.Sprintf(`{"content":"CampusWorld alert: %s — player %s — %d breaks in 1h"}`,
		alert.AlertType, alert.PlayerID.String(), count)
	req, err := http.NewRequest(http.MethodPost, s.cfg.AlertsDiscordWebhook, strings.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("discord alert: %v", err)
		return
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
}
