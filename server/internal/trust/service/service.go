package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	playerrepo "github.com/woragis/minecraft-campus-backend/server/internal/player/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/trust"
	trustrepo "github.com/woragis/minecraft-campus-backend/server/internal/trust/repository"
	"gorm.io/gorm"
)

type Service struct {
	repo       *trustrepo.Repository
	playerRepo *playerrepo.Repository
}

func New(repo *trustrepo.Repository, playerRepo *playerrepo.Repository) *Service {
	return &Service{repo: repo, playerRepo: playerRepo}
}

func (s *Service) RecordEvent(ctx context.Context, playerID uuid.UUID, eventType, reason string, actorID *uuid.UUID) (*models.TrustEvent, *models.Player, error) {
	player, err := s.playerRepo.FindByID(ctx, playerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, apperrors.NotFound(apperrors.CodeTrustEventV1ServicePlayerNotFound, apperrors.MsgTrustEventV1ServicePlayerNotFound)
		}
		return nil, nil, apperrors.InternalCause(apperrors.CodeTrustEventV1ServiceRecordFailed, apperrors.MsgTrustEventV1ServiceRecordFailed, err)
	}

	delta := trust.DeltaForEvent(eventType)
	if eventType == models.TrustEventBan {
		player.Status = models.PlayerStatusBanned
	}
	if eventType == models.TrustEventUnban {
		player.Status = models.PlayerStatusActive
		delta = 0
	}

	player.TrustScore = trust.ApplyDelta(player.TrustScore, delta)
	event := &models.TrustEvent{
		ID:            uuid.New(),
		PlayerID:      playerID,
		EventType:     eventType,
		Delta:         delta,
		Reason:        strings.TrimSpace(reason),
		ActorPlayerID: actorID,
	}
	if err := s.repo.CreateEvent(ctx, event); err != nil {
		return nil, nil, apperrors.InternalCause(apperrors.CodeTrustEventV1ServiceRecordFailed, apperrors.MsgTrustEventV1ServiceRecordFailed, err)
	}
	if err := s.playerRepo.Update(ctx, player); err != nil {
		return nil, nil, apperrors.InternalCause(apperrors.CodeTrustEventV1ServiceRecordFailed, apperrors.MsgTrustEventV1ServiceRecordFailed, err)
	}

	if player.InvitedByID != nil && (eventType == models.TrustEventBan || eventType == models.TrustEventConfirmedReport) {
		_ = s.RecalculateSponsorScore(ctx, *player.InvitedByID)
	}

	return event, player, nil
}

func (s *Service) ListEvents(ctx context.Context, playerID uuid.UUID) ([]models.TrustEvent, error) {
	if _, err := s.playerRepo.FindByID(ctx, playerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeTrustListV1ServicePlayerNotFound, apperrors.MsgTrustListV1ServicePlayerNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeTrustListV1ServiceListFailed, apperrors.MsgTrustListV1ServiceListFailed, err)
	}
	events, err := s.repo.ListByPlayer(ctx, playerID, 50)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeTrustListV1ServiceListFailed, apperrors.MsgTrustListV1ServiceListFailed, err)
	}
	return events, nil
}

func (s *Service) RecalculateSponsorScore(ctx context.Context, sponsorID uuid.UUID) error {
	invitees, err := s.playerRepo.ListInvitedBy(ctx, sponsorID)
	if err != nil {
		return err
	}
	sponsor, err := s.playerRepo.FindByID(ctx, sponsorID)
	if err != nil {
		return err
	}

	scores := make([]int, 0, len(invitees))
	for _, inv := range invitees {
		scores = append(scores, inv.TrustScore)
	}
	since := time.Now().UTC().Add(-30 * 24 * time.Hour)
	bans, err := s.repo.CountRecentBansAmongInvitees(ctx, sponsorID, since)
	if err != nil {
		return err
	}
	sponsor.SponsorScore = trust.ComputeSponsorScore(scores, int(bans))
	return s.playerRepo.Update(ctx, sponsor)
}

func (s *Service) TryGraduateProbation(ctx context.Context, player *models.Player) (bool, error) {
	if player.Status != models.PlayerStatusProbation {
		return false, nil
	}
	if player.ProbationUntil == nil || time.Now().UTC().Before(*player.ProbationUntil) {
		return false, nil
	}
	player.Status = models.PlayerStatusActive
	player.ProbationUntil = nil
	if err := s.playerRepo.Update(ctx, player); err != nil {
		return false, err
	}
	_, _, err := s.RecordEvent(ctx, player.ID, models.TrustEventProbationClean, "probation completed", nil)
	return true, err
}

// SponsorTree returns direct invitees with trust info for moderation.
func (s *Service) SponsorTree(ctx context.Context, playerID uuid.UUID) ([]models.Player, error) {
	if _, err := s.playerRepo.FindByID(ctx, playerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeTrustTreeV1ServicePlayerNotFound, apperrors.MsgTrustTreeV1ServicePlayerNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeTrustTreeV1ServiceLoadFailed, apperrors.MsgTrustTreeV1ServiceLoadFailed, err)
	}
	return s.playerRepo.ListInvitedBy(ctx, playerID)
}
