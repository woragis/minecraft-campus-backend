package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	"github.com/woragis/minecraft-campus-backend/server/internal/config"
	auditrepo "github.com/woragis/minecraft-campus-backend/server/internal/audit/repository"
	gameserverrepo "github.com/woragis/minecraft-campus-backend/server/internal/gameserver/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	playerrepo "github.com/woragis/minecraft-campus-backend/server/internal/player/repository"
	"gorm.io/gorm"
)

type Service struct {
	cfg        config.Config
	repo       *auditrepo.Repository
	playerRepo *playerrepo.Repository
	serverRepo *gameserverrepo.Repository
}

func New(cfg config.Config, repo *auditrepo.Repository, playerRepo *playerrepo.Repository, serverRepo *gameserverrepo.Repository) *Service {
	return &Service{cfg: cfg, repo: repo, playerRepo: playerRepo, serverRepo: serverRepo}
}

type IngestEvent struct {
	MinecraftUUID uuid.UUID
	ServerSlug    string
	World         string
	EventType     string
	BlockX        *int
	BlockY        *int
	BlockZ        *int
	BlockType     string
	ClaimID       *uuid.UUID
	OccurredAt    time.Time
}

type IngestResult struct {
	Accepted int `json:"accepted"`
}

func (s *Service) IngestBatch(ctx context.Context, events []IngestEvent) (*IngestResult, error) {
	if len(events) == 0 {
		return &IngestResult{}, nil
	}
	if len(events) > 500 {
		return nil, apperrors.Invalid(apperrors.CodeAuditIngestV1ServiceBatchTooLarge, apperrors.MsgAuditIngestV1ServiceBatchTooLarge)
	}

	records := make([]models.AuditEvent, 0, len(events))
	for _, in := range events {
		if !validAuditEventType(in.EventType) {
			continue
		}
		player, err := s.playerRepo.FindByMinecraftUUID(ctx, in.MinecraftUUID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return nil, apperrors.InternalCause(apperrors.CodeAuditIngestV1ServiceFailed, apperrors.MsgAuditIngestV1ServiceFailed, err)
		}
		world := strings.TrimSpace(in.World)
		if world == "" {
			world = "world"
		}
		gs, err := s.serverRepo.GetOrCreateBySlug(ctx, in.ServerSlug, in.ServerSlug)
		if err != nil {
			return nil, apperrors.InternalCause(apperrors.CodeAuditIngestV1ServiceFailed, apperrors.MsgAuditIngestV1ServiceFailed, err)
		}
		occurred := in.OccurredAt
		if occurred.IsZero() {
			occurred = time.Now().UTC()
		}
		records = append(records, models.AuditEvent{
			ID:         uuid.New(),
			ServerID:   gs.ID,
			World:      world,
			PlayerID:   player.ID,
			EventType:  in.EventType,
			BlockX:     in.BlockX,
			BlockY:     in.BlockY,
			BlockZ:     in.BlockZ,
			BlockType:  strings.TrimSpace(in.BlockType),
			ClaimID:    in.ClaimID,
			Metadata:   "{}",
			OccurredAt: occurred.UTC(),
		})
	}
	if err := s.repo.CreateBatch(ctx, records); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeAuditIngestV1ServiceFailed, apperrors.MsgAuditIngestV1ServiceFailed, err)
	}
	return &IngestResult{Accepted: len(records)}, nil
}

func (s *Service) ListByPlayer(ctx context.Context, playerID uuid.UUID, from, to *time.Time, eventType string) ([]models.AuditEvent, error) {
	if _, err := s.playerRepo.FindByID(ctx, playerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeAuditListV1ServicePlayerNotFound, apperrors.MsgAuditListV1ServicePlayerNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeAuditListV1ServiceListFailed, apperrors.MsgAuditListV1ServiceListFailed, err)
	}
	events, err := s.repo.ListByPlayer(ctx, auditrepo.ListFilter{
		PlayerID:  playerID,
		From:      from,
		To:        to,
		EventType: eventType,
		Limit:     200,
	})
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeAuditListV1ServiceListFailed, apperrors.MsgAuditListV1ServiceListFailed, err)
	}
	return events, nil
}

func (s *Service) PurgeOld(ctx context.Context) (int64, error) {
	if !s.cfg.AuditPurgeEnabled || s.cfg.AuditRetentionDays <= 0 {
		return 0, nil
	}
	cutoff := time.Now().UTC().Add(-time.Duration(s.cfg.AuditRetentionDays) * 24 * time.Hour)
	n, err := s.repo.PurgeBefore(ctx, cutoff)
	if err != nil {
		return 0, err
	}
	if n > 0 {
		log.Printf("audit purge: removed %d events older than %d days", n, s.cfg.AuditRetentionDays)
	}
	return n, nil
}

func validAuditEventType(t string) bool {
	switch t {
	case models.AuditEventBlockPlace, models.AuditEventBlockBreak, models.AuditEventPlayerJoin, models.AuditEventPlayerQuit:
		return true
	default:
		return false
	}
}
