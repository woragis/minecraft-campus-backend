package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	"github.com/woragis/minecraft-campus-backend/server/internal/guild/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	playerrepo "github.com/woragis/minecraft-campus-backend/server/internal/player/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/trust"
	"gorm.io/gorm"
)

type Service struct {
	repo       *repository.Repository
	playerRepo *playerrepo.Repository
}

func New(repo *repository.Repository, playerRepo *playerrepo.Repository) *Service {
	return &Service{repo: repo, playerRepo: playerRepo}
}

func (s *Service) Create(ctx context.Context, leaderMinecraftUUID uuid.UUID, name string) (*models.Guild, error) {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 64 {
		return nil, apperrors.Invalid(apperrors.CodeGuildPostV1ServiceNameInvalid, apperrors.MsgGuildPostV1ServiceNameInvalid)
	}

	leader, err := s.playerRepo.FindByMinecraftUUID(ctx, leaderMinecraftUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeGuildPostV1ServiceLeaderNotFound, apperrors.MsgGuildPostV1ServiceLeaderNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeGuildPostV1ServiceCreateFailed, apperrors.MsgGuildPostV1ServiceCreateFailed, err)
	}
	if err := s.assertCanManageGuild(ctx, leader); err != nil {
		return nil, err
	}
	if _, err := s.repo.PlayerGuild(ctx, leader.ID); err == nil {
		return nil, apperrors.ConflictErr(apperrors.CodeGuildPostV1ServiceAlreadyMember, apperrors.MsgGuildPostV1ServiceAlreadyMember)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.InternalCause(apperrors.CodeGuildPostV1ServiceCreateFailed, apperrors.MsgGuildPostV1ServiceCreateFailed, err)
	}

	slug, err := s.uniqueSlug(ctx, slugify(name))
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeGuildPostV1ServiceCreateFailed, apperrors.MsgGuildPostV1ServiceCreateFailed, err)
	}

	guild := &models.Guild{
		ID:         uuid.New(),
		Slug:       slug,
		Name:       name,
		LeaderID:   leader.ID,
		TrustScore: leader.TrustScore,
	}
	member := &models.GuildMember{
		GuildID:  guild.ID,
		PlayerID: leader.ID,
		Role:     models.GuildRoleLeader,
	}
	if err := s.repo.Create(ctx, guild, member); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeGuildPostV1ServiceCreateFailed, apperrors.MsgGuildPostV1ServiceCreateFailed, err)
	}
	return guild, nil
}

func (s *Service) Join(ctx context.Context, guildID, playerMinecraftUUID uuid.UUID) error {
	player, err := s.playerRepo.FindByMinecraftUUID(ctx, playerMinecraftUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeGuildJoinV1ServicePlayerNotFound, apperrors.MsgGuildJoinV1ServicePlayerNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeGuildJoinV1ServiceFailed, apperrors.MsgGuildJoinV1ServiceFailed, err)
	}
	if player.Status == models.PlayerStatusProbation {
		return apperrors.Forbidden(apperrors.CodeGuildJoinV1ServiceProbation, apperrors.MsgGuildJoinV1ServiceProbation)
	}
	if _, err := s.repo.PlayerGuild(ctx, player.ID); err == nil {
		return apperrors.ConflictErr(apperrors.CodeGuildJoinV1ServiceAlreadyMember, apperrors.MsgGuildJoinV1ServiceAlreadyMember)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrors.InternalCause(apperrors.CodeGuildJoinV1ServiceFailed, apperrors.MsgGuildJoinV1ServiceFailed, err)
	}
	if _, err := s.repo.FindByID(ctx, guildID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeGuildJoinV1ServiceGuildNotFound, apperrors.MsgGuildJoinV1ServiceGuildNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeGuildJoinV1ServiceFailed, apperrors.MsgGuildJoinV1ServiceFailed, err)
	}
	member := &models.GuildMember{
		GuildID:  guildID,
		PlayerID: player.ID,
		Role:     models.GuildRoleMember,
	}
	if err := s.repo.AddMember(ctx, member); err != nil {
		return apperrors.InternalCause(apperrors.CodeGuildJoinV1ServiceFailed, apperrors.MsgGuildJoinV1ServiceFailed, err)
	}
	return s.refreshGuildTrust(ctx, guildID)
}

func (s *Service) JoinBySlug(ctx context.Context, slug string, playerMinecraftUUID uuid.UUID) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return apperrors.Invalid(apperrors.CodeGuildGetV1HandlerIDInvalid, apperrors.MsgGuildGetV1HandlerIDInvalid)
	}
	guild, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeGuildJoinV1ServiceGuildNotFound, apperrors.MsgGuildJoinV1ServiceGuildNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeGuildJoinV1ServiceFailed, apperrors.MsgGuildJoinV1ServiceFailed, err)
	}
	return s.Join(ctx, guild.ID, playerMinecraftUUID)
}

func (s *Service) LeaveBySlug(ctx context.Context, slug string, playerMinecraftUUID uuid.UUID) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return apperrors.Invalid(apperrors.CodeGuildGetV1HandlerIDInvalid, apperrors.MsgGuildGetV1HandlerIDInvalid)
	}
	guild, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeGuildLeaveV1ServiceGuildNotFound, apperrors.MsgGuildLeaveV1ServiceGuildNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeGuildLeaveV1ServiceFailed, apperrors.MsgGuildLeaveV1ServiceFailed, err)
	}
	return s.Leave(ctx, guild.ID, playerMinecraftUUID)
}

func (s *Service) Leave(ctx context.Context, guildID, playerMinecraftUUID uuid.UUID) error {
	player, err := s.playerRepo.FindByMinecraftUUID(ctx, playerMinecraftUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeGuildLeaveV1ServicePlayerNotFound, apperrors.MsgGuildLeaveV1ServicePlayerNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeGuildLeaveV1ServiceFailed, apperrors.MsgGuildLeaveV1ServiceFailed, err)
	}
	guild, err := s.repo.FindByID(ctx, guildID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeGuildLeaveV1ServiceGuildNotFound, apperrors.MsgGuildLeaveV1ServiceGuildNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeGuildLeaveV1ServiceFailed, apperrors.MsgGuildLeaveV1ServiceFailed, err)
	}
	if guild.LeaderID == player.ID {
		return apperrors.Forbidden(apperrors.CodeGuildLeaveV1ServiceLeaderCannotLeave, apperrors.MsgGuildLeaveV1ServiceLeaderCannotLeave)
	}
	if err := s.repo.RemoveMember(ctx, guildID, player.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeGuildLeaveV1ServiceNotMember, apperrors.MsgGuildLeaveV1ServiceNotMember)
		}
		return apperrors.InternalCause(apperrors.CodeGuildLeaveV1ServiceFailed, apperrors.MsgGuildLeaveV1ServiceFailed, err)
	}
	return s.refreshGuildTrust(ctx, guildID)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Guild, error) {
	guild, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeGuildGetV1ServiceNotFound, apperrors.MsgGuildGetV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeGuildGetV1ServiceLoadFailed, apperrors.MsgGuildGetV1ServiceLoadFailed, err)
	}
	return s.enrichGuild(ctx, guild)
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*models.Guild, error) {
	guild, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeGuildGetV1ServiceNotFound, apperrors.MsgGuildGetV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeGuildGetV1ServiceLoadFailed, apperrors.MsgGuildGetV1ServiceLoadFailed, err)
	}
	return s.enrichGuild(ctx, guild)
}

func (s *Service) List(ctx context.Context) ([]models.Guild, error) {
	guilds, err := s.repo.List(ctx, 100)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeGuildListV1ServiceListFailed, apperrors.MsgGuildListV1ServiceListFailed, err)
	}
	for i := range guilds {
		count, _ := s.repo.MemberCount(ctx, guilds[i].ID)
		guilds[i].MemberCount = count
	}
	return guilds, nil
}

func (s *Service) ListMembers(ctx context.Context, guildID uuid.UUID) ([]models.GuildMember, error) {
	if _, err := s.repo.FindByID(ctx, guildID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeGuildMembersV1ServiceGuildNotFound, apperrors.MsgGuildMembersV1ServiceGuildNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeGuildMembersV1ServiceListFailed, apperrors.MsgGuildMembersV1ServiceListFailed, err)
	}
	members, err := s.repo.ListMembers(ctx, guildID)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeGuildMembersV1ServiceListFailed, apperrors.MsgGuildMembersV1ServiceListFailed, err)
	}
	return members, nil
}

func (s *Service) PlayerGuild(ctx context.Context, playerID uuid.UUID) (*models.Guild, error) {
	guild, err := s.repo.PlayerGuild(ctx, playerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, apperrors.InternalCause(apperrors.CodeGuildGetV1ServiceLoadFailed, apperrors.MsgGuildGetV1ServiceLoadFailed, err)
	}
	return s.enrichGuild(ctx, guild)
}

func (s *Service) assertCanManageGuild(ctx context.Context, player *models.Player) error {
	if player.Status == models.PlayerStatusBanned {
		return apperrors.Forbidden(apperrors.CodeGuildPostV1ServiceLeaderBanned, apperrors.MsgGuildPostV1ServiceLeaderBanned)
	}
	if player.Status == models.PlayerStatusProbation {
		return apperrors.Forbidden(apperrors.CodeGuildPostV1ServiceLeaderProbation, apperrors.MsgGuildPostV1ServiceLeaderProbation)
	}
	return nil
}

func (s *Service) uniqueSlug(ctx context.Context, base string) (string, error) {
	candidate := base
	for i := 0; i < 100; i++ {
		exists, err := s.repo.SlugExists(ctx, candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, i+2)
	}
	return "", fmt.Errorf("could not allocate unique slug for %q", base)
}

func (s *Service) refreshGuildTrust(ctx context.Context, guildID uuid.UUID) error {
	guild, err := s.repo.FindByID(ctx, guildID)
	if err != nil {
		return err
	}
	avg, err := s.repo.AverageMemberTrust(ctx, guildID)
	if err != nil {
		return err
	}
	guild.TrustScore = trust.ApplyDelta(avg, 0)
	return s.repo.Update(ctx, guild)
}

func (s *Service) enrichGuild(ctx context.Context, guild *models.Guild) (*models.Guild, error) {
	count, err := s.repo.MemberCount(ctx, guild.ID)
	if err != nil {
		return nil, err
	}
	guild.MemberCount = count
	return guild, nil
}
