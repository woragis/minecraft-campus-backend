package service

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	"gorm.io/gorm"
)

// BedrockMinecraftUUID returns a deterministic UUID for a Bedrock XUID so existing
// player tables keyed by minecraft_uuid keep working.
func BedrockMinecraftUUID(xuid string) uuid.UUID {
	xuid = strings.TrimSpace(xuid)
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte("bedrock:"+xuid))
}

func validXUID(xuid string) bool {
	xuid = strings.TrimSpace(xuid)
	if xuid == "" {
		return false
	}
	for _, r := range xuid {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func (s *Service) CheckBedrockWhitelist(ctx context.Context, xuid, username string) (*WhitelistResult, error) {
	xuid = strings.TrimSpace(xuid)
	if !validXUID(xuid) {
		return nil, apperrors.Invalid(apperrors.CodeBedrockWhitelistGetV1HandlerXUIDInvalid, apperrors.MsgBedrockWhitelistGetV1HandlerXUIDInvalid)
	}
	username = strings.TrimSpace(username)

	identity, err := s.repo.FindIdentityByPlatformExternalID(ctx, models.PlatformBedrock, xuid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.InternalCause(apperrors.CodeBedrockWhitelistGetV1ServiceCheckFailed, apperrors.MsgBedrockWhitelistGetV1ServiceCheckFailed, err)
	}
	if identity != nil {
		player, err := s.repo.FindByID(ctx, identity.PlayerID)
		if err != nil {
			return nil, apperrors.InternalCause(apperrors.CodeBedrockWhitelistGetV1ServiceCheckFailed, apperrors.MsgBedrockWhitelistGetV1ServiceCheckFailed, err)
		}
		s.maybeGraduate(ctx, player)
		return s.whitelistExisting(player), nil
	}

	mcUUID := BedrockMinecraftUUID(xuid)
	player, err := s.repo.FindByMinecraftUUID(ctx, mcUUID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.InternalCause(apperrors.CodeBedrockWhitelistGetV1ServiceCheckFailed, apperrors.MsgBedrockWhitelistGetV1ServiceCheckFailed, err)
	}
	if player != nil {
		s.maybeGraduate(ctx, player)
		_ = s.repo.UpsertIdentity(ctx, &models.PlayerIdentity{
			PlayerID:   player.ID,
			Platform:   models.PlatformBedrock,
			ExternalID: xuid,
			Username:   player.Username,
		})
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
		return nil, apperrors.InternalCause(apperrors.CodeBedrockWhitelistGetV1ServiceCheckFailed, apperrors.MsgBedrockWhitelistGetV1ServiceCheckFailed, err)
	}

	now := time.Now().UTC()
	probationUntil := now.Add(time.Duration(s.probationDaysFor(invite)) * 24 * time.Hour)
	sponsorID := invite.SponsorID

	newPlayer := &models.Player{
		ID:             uuid.New(),
		MinecraftUUID:  mcUUID,
		Username:       username,
		Status:         models.PlayerStatusProbation,
		InvitedByID:    &sponsorID,
		TrustScore:     100,
		SponsorScore:   100,
		ProbationUntil: &probationUntil,
	}
	applyAffiliationFromInvite(newPlayer, invite)
	if err := s.repo.Create(ctx, newPlayer); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeBedrockWhitelistGetV1ServiceCheckFailed, apperrors.MsgBedrockWhitelistGetV1ServiceCheckFailed, err)
	}
	if err := s.repo.UpsertIdentity(ctx, &models.PlayerIdentity{
		PlayerID:   newPlayer.ID,
		Platform:   models.PlatformBedrock,
		ExternalID: xuid,
		Username:   username,
	}); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeBedrockWhitelistGetV1ServiceCheckFailed, apperrors.MsgBedrockWhitelistGetV1ServiceCheckFailed, err)
	}
	if err := s.inviteRepo.Accept(ctx, invite.ID, newPlayer.ID, now); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeBedrockWhitelistGetV1ServiceCheckFailed, apperrors.MsgBedrockWhitelistGetV1ServiceCheckFailed, err)
	}

	return &WhitelistResult{
		Allowed: true,
		Reason:  "probation",
		Player:  newPlayer,
	}, nil
}

func (s *Service) UpsertBedrockFromServer(ctx context.Context, xuid, username, serverSlug string) (*models.Player, error) {
	xuid = strings.TrimSpace(xuid)
	if !validXUID(xuid) {
		return nil, apperrors.Invalid(apperrors.CodeBedrockPlayerUpsertV1HandlerBodyInvalid, apperrors.MsgBedrockPlayerUpsertV1HandlerBodyInvalid)
	}
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, apperrors.Invalid(apperrors.CodeBedrockPlayerUpsertV1HandlerBodyInvalid, apperrors.MsgBedrockPlayerUpsertV1HandlerBodyInvalid)
	}

	gs, err := s.gameServerRepo.GetOrCreateBySlug(ctx, serverSlug, serverSlug)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeBedrockPlayerUpsertV1ServiceFailed, apperrors.MsgBedrockPlayerUpsertV1ServiceFailed, err)
	}

	mcUUID := BedrockMinecraftUUID(xuid)
	now := time.Now().UTC()

	var player *models.Player
	identity, idErr := s.repo.FindIdentityByPlatformExternalID(ctx, models.PlatformBedrock, xuid)
	if idErr != nil && !errors.Is(idErr, gorm.ErrRecordNotFound) {
		return nil, apperrors.InternalCause(apperrors.CodeBedrockPlayerUpsertV1ServiceFailed, apperrors.MsgBedrockPlayerUpsertV1ServiceFailed, idErr)
	}
	if identity != nil {
		player, err = s.repo.FindByID(ctx, identity.PlayerID)
		if err != nil {
			return nil, apperrors.InternalCause(apperrors.CodeBedrockPlayerUpsertV1ServiceFailed, apperrors.MsgBedrockPlayerUpsertV1ServiceFailed, err)
		}
	} else {
		player, err = s.repo.FindByMinecraftUUID(ctx, mcUUID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.InternalCause(apperrors.CodeBedrockPlayerUpsertV1ServiceFailed, apperrors.MsgBedrockPlayerUpsertV1ServiceFailed, err)
		}
	}

	if player == nil {
		invite, invErr := s.inviteRepo.FindPendingByTargetUsername(ctx, username)
		if invErr != nil {
			if errors.Is(invErr, gorm.ErrRecordNotFound) {
				return nil, apperrors.Forbidden(apperrors.CodeBedrockWhitelistGetV1ServiceCheckFailed, "Player is not whitelisted.")
			}
			return nil, apperrors.InternalCause(apperrors.CodeBedrockPlayerUpsertV1ServiceFailed, apperrors.MsgBedrockPlayerUpsertV1ServiceFailed, invErr)
		}
		probationUntil := now.Add(time.Duration(s.probationDaysFor(invite)) * 24 * time.Hour)
		sponsorID := invite.SponsorID
		player = &models.Player{
			ID:             uuid.New(),
			MinecraftUUID:  mcUUID,
			Username:       username,
			Status:         models.PlayerStatusProbation,
			InvitedByID:    &sponsorID,
			TrustScore:     100,
			SponsorScore:   100,
			ProbationUntil: &probationUntil,
		}
		applyAffiliationFromInvite(player, invite)
		if err := s.repo.Create(ctx, player); err != nil {
			return nil, apperrors.InternalCause(apperrors.CodeBedrockPlayerUpsertV1ServiceFailed, apperrors.MsgBedrockPlayerUpsertV1ServiceFailed, err)
		}
		if err := s.inviteRepo.Accept(ctx, invite.ID, player.ID, now); err != nil {
			return nil, apperrors.InternalCause(apperrors.CodeBedrockPlayerUpsertV1ServiceFailed, apperrors.MsgBedrockPlayerUpsertV1ServiceFailed, err)
		}
	} else if !strings.EqualFold(player.Username, username) {
		player.Username = username
		if err := s.repo.Update(ctx, player); err != nil {
			return nil, apperrors.InternalCause(apperrors.CodeBedrockPlayerUpsertV1ServiceFailed, apperrors.MsgBedrockPlayerUpsertV1ServiceFailed, err)
		}
	}

	if err := s.repo.UpsertIdentity(ctx, &models.PlayerIdentity{
		PlayerID:   player.ID,
		Platform:   models.PlatformBedrock,
		ExternalID: xuid,
		Username:   username,
	}); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeBedrockPlayerUpsertV1ServiceFailed, apperrors.MsgBedrockPlayerUpsertV1ServiceFailed, err)
	}

	if err := s.gameServerRepo.TouchPlayerPresence(ctx, gs.ID, player.ID, now); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeBedrockPlayerUpsertV1ServiceFailed, apperrors.MsgBedrockPlayerUpsertV1ServiceFailed, err)
	}
	s.maybeGraduate(ctx, player)
	return player, nil
}
