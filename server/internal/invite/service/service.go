package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	"github.com/woragis/minecraft-campus-backend/server/internal/invite/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	playerrepo "github.com/woragis/minecraft-campus-backend/server/internal/player/repository"
	"gorm.io/gorm"
)

type Service struct {
	repo       *repository.Repository
	playerRepo *playerrepo.Repository
}

func New(repo *repository.Repository, playerRepo *playerrepo.Repository) *Service {
	return &Service{repo: repo, playerRepo: playerRepo}
}

func (s *Service) GetByCode(ctx context.Context, code string) (*models.Invite, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, apperrors.NotFound(apperrors.CodeInviteGetV1ServiceNotFound, apperrors.MsgInviteGetV1ServiceNotFound)
	}
	inv, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInviteGetV1ServiceNotFound, apperrors.MsgInviteGetV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeInviteGetV1ServiceLoadFailed, apperrors.MsgInviteGetV1ServiceLoadFailed, err)
	}
	return inv, nil
}

func (s *Service) CreateForSponsor(ctx context.Context, sponsorMinecraftUUID uuid.UUID, targetUsername string) (*models.Invite, error) {
	targetUsername = strings.TrimSpace(targetUsername)
	if targetUsername == "" {
		return nil, apperrors.Invalid(apperrors.CodeInvitePostInternalV1HandlerBodyInvalid, apperrors.MsgInvitePostInternalV1HandlerBodyInvalid)
	}

	sponsor, err := s.playerRepo.FindByMinecraftUUID(ctx, sponsorMinecraftUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInvitePostInternalV1ServiceSponsorNotFound, apperrors.MsgInvitePostInternalV1ServiceSponsorNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeInvitePostInternalV1ServiceCreateFailed, apperrors.MsgInvitePostInternalV1ServiceCreateFailed, err)
	}
	if sponsor.Status == models.PlayerStatusBanned {
		return nil, apperrors.Forbidden(apperrors.CodeInvitePostInternalV1ServiceSponsorBanned, apperrors.MsgInvitePostInternalV1ServiceSponsorBanned)
	}
	if sponsor.Status == models.PlayerStatusProbation {
		return nil, apperrors.Forbidden(apperrors.CodeInvitePostInternalV1ServiceSponsorProbation, apperrors.MsgInvitePostInternalV1ServiceSponsorProbation)
	}

	exists, err := s.repo.HasPendingForTargetUsername(ctx, targetUsername)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInvitePostInternalV1ServiceCreateFailed, apperrors.MsgInvitePostInternalV1ServiceCreateFailed, err)
	}
	if exists {
		return nil, apperrors.ConflictErr(apperrors.CodeInvitePostInternalV1ServicePendingExists, apperrors.MsgInvitePostInternalV1ServicePendingExists)
	}

	code, err := generateInviteCode()
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInvitePostInternalV1ServiceCreateFailed, apperrors.MsgInvitePostInternalV1ServiceCreateFailed, err)
	}

	inv := &models.Invite{
		ID:             uuid.New(),
		Code:           code,
		SponsorID:      sponsor.ID,
		TargetUsername: targetUsername,
		Status:         models.InviteStatusPending,
	}
	if err := s.repo.Create(ctx, inv); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInvitePostInternalV1ServiceCreateFailed, apperrors.MsgInvitePostInternalV1ServiceCreateFailed, err)
	}
	return inv, nil
}

func (s *Service) FindPendingByTargetUsername(ctx context.Context, username string) (*models.Invite, error) {
	return s.repo.FindPendingByTargetUsername(ctx, username)
}

func (s *Service) AcceptForPlayer(ctx context.Context, inviteID, playerID uuid.UUID) error {
	return s.repo.Accept(ctx, inviteID, playerID, time.Now().UTC())
}

func (s *Service) ListBySponsorPlayerID(ctx context.Context, sponsorPlayerID uuid.UUID) ([]models.Invite, error) {
	return s.repo.ListBySponsor(ctx, sponsorPlayerID)
}
