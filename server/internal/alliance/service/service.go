package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/alliance/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	guildrepo "github.com/woragis/minecraft-campus-backend/server/internal/guild/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	playerrepo "github.com/woragis/minecraft-campus-backend/server/internal/player/repository"
	"gorm.io/gorm"
)

type Service struct {
	repo       *repository.Repository
	guildRepo  *guildrepo.Repository
	playerRepo *playerrepo.Repository
}

func New(repo *repository.Repository, guildRepo *guildrepo.Repository, playerRepo *playerrepo.Repository) *Service {
	return &Service{repo: repo, guildRepo: guildRepo, playerRepo: playerRepo}
}

func (s *Service) Create(ctx context.Context, leaderMinecraftUUID, guildAID, guildBID uuid.UUID) (*models.Alliance, error) {
	if guildAID == guildBID {
		return nil, apperrors.Invalid(apperrors.CodeAlliancePostV1ServiceSameGuild, apperrors.MsgAlliancePostV1ServiceSameGuild)
	}
	leader, err := s.playerRepo.FindByMinecraftUUID(ctx, leaderMinecraftUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeAlliancePostV1ServiceLeaderNotFound, apperrors.MsgAlliancePostV1ServiceLeaderNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeAlliancePostV1ServiceCreateFailed, apperrors.MsgAlliancePostV1ServiceCreateFailed, err)
	}
	guildA, err := s.guildRepo.FindByID(ctx, guildAID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeAlliancePostV1ServiceGuildNotFound, apperrors.MsgAlliancePostV1ServiceGuildNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeAlliancePostV1ServiceCreateFailed, apperrors.MsgAlliancePostV1ServiceCreateFailed, err)
	}
	if guildA.LeaderID != leader.ID {
		return nil, apperrors.Forbidden(apperrors.CodeAlliancePostV1ServiceNotLeader, apperrors.MsgAlliancePostV1ServiceNotLeader)
	}
	if _, err := s.guildRepo.FindByID(ctx, guildBID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeAlliancePostV1ServiceGuildNotFound, apperrors.MsgAlliancePostV1ServiceGuildNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeAlliancePostV1ServiceCreateFailed, apperrors.MsgAlliancePostV1ServiceCreateFailed, err)
	}
	exists, err := s.repo.Exists(ctx, guildAID, guildBID)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeAlliancePostV1ServiceCreateFailed, apperrors.MsgAlliancePostV1ServiceCreateFailed, err)
	}
	if exists {
		return nil, apperrors.ConflictErr(apperrors.CodeAlliancePostV1ServiceExists, apperrors.MsgAlliancePostV1ServiceExists)
	}

	alliance := &models.Alliance{
		ID:       uuid.New(),
		GuildAID: guildAID,
		GuildBID: guildBID,
		Status:   models.AllianceStatusActive,
	}
	if err := s.repo.Create(ctx, alliance); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeAlliancePostV1ServiceCreateFailed, apperrors.MsgAlliancePostV1ServiceCreateFailed, err)
	}
	return alliance, nil
}

func (s *Service) ListByGuild(ctx context.Context, guildID uuid.UUID) ([]models.Alliance, error) {
	if _, err := s.guildRepo.FindByID(ctx, guildID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeAllianceListV1ServiceGuildNotFound, apperrors.MsgAllianceListV1ServiceGuildNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeAllianceListV1ServiceListFailed, apperrors.MsgAllianceListV1ServiceListFailed, err)
	}
	alliances, err := s.repo.ListByGuild(ctx, guildID)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeAllianceListV1ServiceListFailed, apperrors.MsgAllianceListV1ServiceListFailed, err)
	}
	return alliances, nil
}
