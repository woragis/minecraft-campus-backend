package presence

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	guildrepo "github.com/woragis/minecraft-campus-backend/server/internal/guild/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	playerrepo "github.com/woragis/minecraft-campus-backend/server/internal/player/repository"
	"gorm.io/gorm"
)

var defaultServerSlugs = []string{"bedrock", "crossplay", "vanilla"}

type Service struct {
	store      Store
	guildRepo  *guildrepo.Repository
	playerRepo *playerrepo.Repository
}

func New(store Store, guildRepo *guildrepo.Repository, playerRepo *playerrepo.Repository) *Service {
	return &Service{
		store:      store,
		guildRepo:  guildRepo,
		playerRepo: playerRepo,
	}
}

func (s *Service) Enabled() bool {
	return s.store.Enabled()
}

func (s *Service) MarkOnline(ctx context.Context, playerID uuid.UUID, username, serverSlug, platform string) error {
	if playerID == uuid.Nil {
		return apperrors.Invalid(apperrors.CodePresenceOnlineV1HandlerBodyInvalid, apperrors.MsgPresenceOnlineV1HandlerBodyInvalid)
	}
	var guildID *uuid.UUID
	if guild, err := s.guildRepo.PlayerGuild(ctx, playerID); err == nil && guild != nil {
		guildID = &guild.ID
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrors.InternalCause(apperrors.CodePresenceOnlineV1ServiceFailed, apperrors.MsgPresenceOnlineV1ServiceFailed, err)
	}
	return s.store.MarkOnline(ctx, MarkOnlineInput{
		PlayerID:   playerID,
		Username:   strings.TrimSpace(username),
		ServerSlug: serverSlug,
		Platform:   platform,
		GuildID:    guildID,
	})
}

func (s *Service) MarkOffline(ctx context.Context, playerID uuid.UUID, serverSlug string) error {
	if playerID == uuid.Nil {
		return apperrors.Invalid(apperrors.CodePresenceOfflineV1HandlerBodyInvalid, apperrors.MsgPresenceOfflineV1HandlerBodyInvalid)
	}
	var guildID *uuid.UUID
	if guild, err := s.guildRepo.PlayerGuild(ctx, playerID); err == nil && guild != nil {
		guildID = &guild.ID
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrors.InternalCause(apperrors.CodePresenceOfflineV1ServiceFailed, apperrors.MsgPresenceOfflineV1ServiceFailed, err)
	}
	return s.store.MarkOffline(ctx, playerID, serverSlug, guildID)
}

func (s *Service) Heartbeat(ctx context.Context, playerID uuid.UUID, serverSlug string) error {
	if playerID == uuid.Nil || strings.TrimSpace(serverSlug) == "" {
		return apperrors.Invalid(apperrors.CodePresenceHeartbeatV1HandlerBodyInvalid, apperrors.MsgPresenceHeartbeatV1HandlerBodyInvalid)
	}
	return s.store.Heartbeat(ctx, playerID, serverSlug)
}

func (s *Service) Overview(ctx context.Context) (*Overview, error) {
	servers, err := s.store.ListServersOnline(ctx, defaultServerSlugs)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodePresenceOverviewV1ServiceFailed, apperrors.MsgPresenceOverviewV1ServiceFailed, err)
	}
	total := 0
	for _, srv := range servers {
		total += srv.OnlineCount
	}
	return &Overview{
		Enabled:     s.store.Enabled(),
		TotalOnline: total,
		Servers:     servers,
	}, nil
}

func (s *Service) Server(ctx context.Context, slug string) (*ServerPresence, error) {
	slug = strings.ToLower(strings.TrimSpace(slug))
	if slug == "" {
		return nil, apperrors.Invalid(apperrors.CodePresenceServerV1HandlerSlugInvalid, apperrors.MsgPresenceServerV1HandlerSlugInvalid)
	}
	players, err := s.store.ListServerOnline(ctx, slug)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodePresenceServerV1ServiceFailed, apperrors.MsgPresenceServerV1ServiceFailed, err)
	}
	return &ServerPresence{
		Slug:        slug,
		OnlineCount: len(players),
		Players:     players,
	}, nil
}

func (s *Service) Guild(ctx context.Context, guildID uuid.UUID) (*GuildPresence, error) {
	if guildID == uuid.Nil {
		return nil, apperrors.Invalid(apperrors.CodeGuildGetV1HandlerIDInvalid, apperrors.MsgGuildGetV1HandlerIDInvalid)
	}
	if _, err := s.guildRepo.FindByID(ctx, guildID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeGuildGetV1ServiceNotFound, apperrors.MsgGuildGetV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodePresenceGuildV1ServiceFailed, apperrors.MsgPresenceGuildV1ServiceFailed, err)
	}
	members, err := s.guildRepo.ListMembers(ctx, guildID)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodePresenceGuildV1ServiceFailed, apperrors.MsgPresenceGuildV1ServiceFailed, err)
	}
	memberIDs := make([]uuid.UUID, 0, len(members))
	for _, m := range members {
		memberIDs = append(memberIDs, m.PlayerID)
	}
	online, err := s.store.CountGuildOnline(ctx, guildID, memberIDs)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodePresenceGuildV1ServiceFailed, apperrors.MsgPresenceGuildV1ServiceFailed, err)
	}
	return &GuildPresence{
		GuildID:     guildID,
		OnlineCount: len(online),
		Members:     online,
	}, nil
}

func (s *Service) MarkOnlineForPlayer(ctx context.Context, player *models.Player, serverSlug, platform string) error {
	if player == nil {
		return nil
	}
	return s.MarkOnline(ctx, player.ID, player.Username, serverSlug, platform)
}
