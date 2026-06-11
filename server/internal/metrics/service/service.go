package service

import (
	"context"
	"fmt"
	"time"

	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	"gorm.io/gorm"
)

type Overview struct {
	TotalPlayers      int64 `json:"totalPlayers"`
	ActivePlayers7d   int64 `json:"activePlayers7d"`
	TotalCities       int64 `json:"totalCities"`
	TotalClaims       int64 `json:"totalClaims"`
	TotalClaimArea    int64 `json:"totalClaimArea"`
	Rollbacks7d       int64 `json:"rollbacks7d"`
	AuditBreaks24h    int64 `json:"auditBreaks24h"`
	UnackedAlerts     int64 `json:"unackedAlerts"`
	GeneratedAt       time.Time `json:"generatedAt"`
}

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Overview(ctx context.Context) (*Overview, error) {
	now := time.Now().UTC()
	since7d := now.Add(-7 * 24 * time.Hour)
	since24h := now.Add(-24 * time.Hour)

	out := &Overview{GeneratedAt: now}

	if err := s.db.WithContext(ctx).Model(&models.Player{}).Count(&out.TotalPlayers).Error; err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeMetricsOverviewV1ServiceFailed, apperrors.MsgMetricsOverviewV1ServiceFailed, err)
	}
	if err := s.db.WithContext(ctx).Model(&models.Player{}).Where("updated_at >= ?", since7d).Count(&out.ActivePlayers7d).Error; err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeMetricsOverviewV1ServiceFailed, apperrors.MsgMetricsOverviewV1ServiceFailed, err)
	}
	if err := s.db.WithContext(ctx).Model(&models.City{}).Count(&out.TotalCities).Error; err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeMetricsOverviewV1ServiceFailed, apperrors.MsgMetricsOverviewV1ServiceFailed, err)
	}
	if err := s.db.WithContext(ctx).Model(&models.Claim{}).Count(&out.TotalClaims).Error; err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeMetricsOverviewV1ServiceFailed, apperrors.MsgMetricsOverviewV1ServiceFailed, err)
	}

	type areaRow struct{ Total int64 }
	var area areaRow
	if err := s.db.WithContext(ctx).Model(&models.Claim{}).Select("COALESCE(SUM(area_blocks), 0) AS total").Scan(&area).Error; err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeMetricsOverviewV1ServiceFailed, apperrors.MsgMetricsOverviewV1ServiceFailed, err)
	}
	out.TotalClaimArea = area.Total

	if err := s.db.WithContext(ctx).Model(&models.Rollback{}).Where("created_at >= ?", since7d).Count(&out.Rollbacks7d).Error; err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeMetricsOverviewV1ServiceFailed, apperrors.MsgMetricsOverviewV1ServiceFailed, err)
	}
	if err := s.db.WithContext(ctx).Model(&models.AuditEvent{}).
		Where("event_type = ? AND occurred_at >= ?", models.AuditEventBlockBreak, since24h).
		Count(&out.AuditBreaks24h).Error; err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeMetricsOverviewV1ServiceFailed, apperrors.MsgMetricsOverviewV1ServiceFailed, err)
	}
	if err := s.db.WithContext(ctx).Model(&models.Alert{}).Where("acknowledged = ?", false).Count(&out.UnackedAlerts).Error; err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeMetricsOverviewV1ServiceFailed, apperrors.MsgMetricsOverviewV1ServiceFailed, err)
	}

	return out, nil
}

func (s *Service) Refresh(_ context.Context) error {
	// On-demand queries; no materialized views in budget profile.
	return nil
}

func (s *Service) TerritoryByServer(ctx context.Context) ([]map[string]any, error) {
	type row struct {
		ServerSlug string
		Claims     int64
		Area       int64
		Cities     int64
	}
	var rows []row
	err := s.db.WithContext(ctx).Raw(`
		SELECT gs.slug AS server_slug,
			COUNT(DISTINCT c.id) AS claims,
			COALESCE(SUM(c.area_blocks), 0) AS area,
			(SELECT COUNT(*) FROM cities ci WHERE ci.server_id = gs.id) AS cities
		FROM game_servers gs
		LEFT JOIN claims c ON c.server_id = gs.id
		GROUP BY gs.id, gs.slug
		ORDER BY gs.slug
	`).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("territory metrics: %w", err)
	}
	out := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		out = append(out, map[string]any{
			"serverSlug": r.ServerSlug,
			"claims":     r.Claims,
			"area":       r.Area,
			"cities":     r.Cities,
		})
	}
	return out, nil
}
