package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	gameserverrepo "github.com/woragis/minecraft-campus-backend/server/internal/gameserver/repository"
	inviterepo "github.com/woragis/minecraft-campus-backend/server/internal/invite/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	"github.com/woragis/minecraft-campus-backend/server/internal/player/repository"
	"gorm.io/gorm"
)

type WhitelistResult struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
	Player  *models.Player `json:"player,omitempty"`
}

type Service struct {
	repo           *repository.Repository
	inviteRepo     *inviterepo.Repository
	gameServerRepo *gameserverrepo.Repository
	probationDays  int
}

func New(repo *repository.Repository, inviteRepo *inviterepo.Repository, gameServerRepo *gameserverrepo.Repository, probationDays int) *Service {
	if probationDays <= 0 {
		probationDays = 7
	}
	return &Service{
		repo:           repo,
		inviteRepo:     inviteRepo,
		gameServerRepo: gameServerRepo,
		probationDays:  probationDays,
	}
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Player, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodePlayerGetV1ServiceNotFound, apperrors.MsgPlayerGetV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodePlayerGetV1ServiceLoadFailed, apperrors.MsgPlayerGetV1ServiceLoadFailed, err)
	}
	return p, nil
}

func (s *Service) GetByMinecraftUUID(ctx context.Context, minecraftUUID uuid.UUID) (*models.Player, error) {
	p, err := s.repo.FindByMinecraftUUID(ctx, minecraftUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodePlayerGetMinecraftV1ServiceNotFound, apperrors.MsgPlayerGetMinecraftV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodePlayerGetMinecraftV1ServiceLoadFailed, apperrors.MsgPlayerGetMinecraftV1ServiceLoadFailed, err)
	}
	return p, nil
}

func (s *Service) ListDirectInvitees(ctx context.Context, playerID uuid.UUID) ([]models.Player, error) {
	if _, err := s.repo.FindByID(ctx, playerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodePlayerInvitesListV1ServiceNotFound, apperrors.MsgPlayerInvitesListV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodePlayerInvitesListV1ServiceListFailed, apperrors.MsgPlayerInvitesListV1ServiceListFailed, err)
	}
	players, err := s.repo.ListInvitedBy(ctx, playerID)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodePlayerInvitesListV1ServiceListFailed, apperrors.MsgPlayerInvitesListV1ServiceListFailed, err)
	}
	return players, nil
}

func (s *Service) CheckWhitelist(ctx context.Context, minecraftUUID uuid.UUID, username string) (*WhitelistResult, error) {
	username = strings.TrimSpace(username)

	player, err := s.repo.FindByMinecraftUUID(ctx, minecraftUUID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.InternalCause(apperrors.CodeWhitelistGetV1ServiceCheckFailed, apperrors.MsgWhitelistGetV1ServiceCheckFailed, err)
	}

	if player != nil {
		return s.whitelistExisting(player), nil
	}

	if username == "" {
		return &WhitelistResult{Allowed: false, Reason: "not_invited"}, nil
	}

	invite, err := s.inviteRepo.FindPendingByTargetUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &WhitelistResult{Allowed: false, Reason: "not_invited"}, nil
		}
		return nil, apperrors.InternalCause(apperrors.CodeWhitelistGetV1ServiceCheckFailed, apperrors.MsgWhitelistGetV1ServiceCheckFailed, err)
	}

	now := time.Now().UTC()
	probationUntil := now.Add(time.Duration(s.probationDays) * 24 * time.Hour)
	sponsorID := invite.SponsorID

	newPlayer := &models.Player{
		ID:             uuid.New(),
		MinecraftUUID:  minecraftUUID,
		Username:       username,
		Status:         models.PlayerStatusProbation,
		InvitedByID:    &sponsorID,
		TrustScore:     100,
		SponsorScore:   100,
		ProbationUntil: &probationUntil,
	}
	if err := s.repo.Create(ctx, newPlayer); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeWhitelistGetV1ServiceCheckFailed, apperrors.MsgWhitelistGetV1ServiceCheckFailed, err)
	}
	if err := s.inviteRepo.Accept(ctx, invite.ID, newPlayer.ID, now); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeWhitelistGetV1ServiceCheckFailed, apperrors.MsgWhitelistGetV1ServiceCheckFailed, err)
	}

	return &WhitelistResult{
		Allowed: true,
		Reason:  "probation",
		Player:  newPlayer,
	}, nil
}

func (s *Service) whitelistExisting(player *models.Player) *WhitelistResult {
	if player.Status == models.PlayerStatusBanned {
		return &WhitelistResult{Allowed: false, Reason: "banned", Player: player}
	}
	reason := player.Status
	if reason == "" {
		reason = "active"
	}
	return &WhitelistResult{Allowed: true, Reason: reason, Player: player}
}

func (s *Service) UpsertFromPlugin(ctx context.Context, minecraftUUID uuid.UUID, username, serverSlug string) (*models.Player, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, apperrors.Invalid(apperrors.CodePlayerUpsertV1HandlerBodyInvalid, apperrors.MsgPlayerUpsertV1HandlerBodyInvalid)
	}

	gs, err := s.gameServerRepo.GetOrCreateBySlug(ctx, serverSlug, serverSlug)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodePlayerUpsertV1ServiceFailed, apperrors.MsgPlayerUpsertV1ServiceFailed, err)
	}

	player, err := s.repo.FindByMinecraftUUID(ctx, minecraftUUID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.InternalCause(apperrors.CodePlayerUpsertV1ServiceFailed, apperrors.MsgPlayerUpsertV1ServiceFailed, err)
	}

	now := time.Now().UTC()
	if player == nil {
		// Upsert after whitelist should find the player; create only if invite flow already ran.
		invite, invErr := s.inviteRepo.FindPendingByTargetUsername(ctx, username)
		if invErr != nil {
			if errors.Is(invErr, gorm.ErrRecordNotFound) {
				return nil, apperrors.Forbidden(apperrors.CodeWhitelistGetV1ServiceCheckFailed, "Player is not whitelisted.")
			}
			return nil, apperrors.InternalCause(apperrors.CodePlayerUpsertV1ServiceFailed, apperrors.MsgPlayerUpsertV1ServiceFailed, invErr)
		}
		probationUntil := now.Add(time.Duration(s.probationDays) * 24 * time.Hour)
		sponsorID := invite.SponsorID
		player = &models.Player{
			ID:             uuid.New(),
			MinecraftUUID:  minecraftUUID,
			Username:       username,
			Status:         models.PlayerStatusProbation,
			InvitedByID:    &sponsorID,
			TrustScore:     100,
			SponsorScore:   100,
			ProbationUntil: &probationUntil,
		}
		if err := s.repo.Create(ctx, player); err != nil {
			return nil, apperrors.InternalCause(apperrors.CodePlayerUpsertV1ServiceFailed, apperrors.MsgPlayerUpsertV1ServiceFailed, err)
		}
		if err := s.inviteRepo.Accept(ctx, invite.ID, player.ID, now); err != nil {
			return nil, apperrors.InternalCause(apperrors.CodePlayerUpsertV1ServiceFailed, apperrors.MsgPlayerUpsertV1ServiceFailed, err)
		}
	} else if !strings.EqualFold(player.Username, username) {
		player.Username = username
		if err := s.repo.Update(ctx, player); err != nil {
			return nil, apperrors.InternalCause(apperrors.CodePlayerUpsertV1ServiceFailed, apperrors.MsgPlayerUpsertV1ServiceFailed, err)
		}
	}

	if err := s.gameServerRepo.TouchPlayerPresence(ctx, gs.ID, player.ID, now); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodePlayerUpsertV1ServiceFailed, apperrors.MsgPlayerUpsertV1ServiceFailed, err)
	}
	return player, nil
}
