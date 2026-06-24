package stats

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	affiliationsvc "github.com/woragis/minecraft-campus-backend/server/internal/affiliation/service"
	gameserverrepo "github.com/woragis/minecraft-campus-backend/server/internal/gameserver/repository"
	guildrepo "github.com/woragis/minecraft-campus-backend/server/internal/guild/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/presence"
	playerrepo "github.com/woragis/minecraft-campus-backend/server/internal/player/repository"
	"gorm.io/gorm"
	"time"
)

type ServerStat struct {
	ServerSlug   string `json:"serverSlug"`
	PlayTimeSecs int64  `json:"playTimeSecs"`
	MobKills     int64  `json:"mobKills"`
}

type PlayerStats struct {
	PlayerID          uuid.UUID    `json:"playerId"`
	TotalPlayTimeSecs int64        `json:"totalPlayTimeSecs"`
	TotalMobKills     int64        `json:"totalMobKills"`
	ByServer          []ServerStat `json:"byServer"`
}

type HUD struct {
	Username         string `json:"username"`
	Status           string `json:"status"`
	AffiliationType  string `json:"affiliationType"`
	UniversitySlug   string `json:"universitySlug,omitempty"`
	FacultySlug      string `json:"facultySlug,omitempty"`
	CourseSlug       string `json:"courseSlug,omitempty"`
	UniversityName   string `json:"universityName,omitempty"`
	UniversityHex    string `json:"universityHex,omitempty"`
	FacultyName      string `json:"facultyName,omitempty"`
	FacultyAbbr      string `json:"facultyAbbr,omitempty"`
	FacultyHex       string `json:"facultyHex,omitempty"`
	CourseName       string `json:"courseName,omitempty"`
	CourseAbbr       string `json:"courseAbbr,omitempty"`
	CourseHex        string `json:"courseHex,omitempty"`
	GuildID          string `json:"guildId,omitempty"`
	GuildName        string `json:"guildName,omitempty"`
	GuildSlug        string `json:"guildSlug,omitempty"`
	GuildOnlineCount int    `json:"guildOnlineCount"`
}

type IngestInput struct {
	PlayerID       uuid.UUID
	ServerSlug     string
	SessionSeconds int64
	MobKills       int64
}

type Service struct {
	gameServerRepo *gameserverrepo.Repository
	playerRepo     *playerrepo.Repository
	guildRepo      *guildrepo.Repository
	presence       *presence.Service
	affiliation    *affiliationsvc.Service
}

func New(
	gameServerRepo *gameserverrepo.Repository,
	playerRepo *playerrepo.Repository,
	guildRepo *guildrepo.Repository,
	presenceSvc *presence.Service,
	affiliationSvc *affiliationsvc.Service,
) *Service {
	return &Service{
		gameServerRepo: gameServerRepo,
		playerRepo:     playerRepo,
		guildRepo:      guildRepo,
		presence:       presenceSvc,
		affiliation:    affiliationSvc,
	}
}

func (s *Service) Ingest(ctx context.Context, in IngestInput) error {
	if in.PlayerID == uuid.Nil {
		return apperrors.Invalid(apperrors.CodeStatsIngestV1HandlerBodyInvalid, apperrors.MsgStatsIngestV1HandlerBodyInvalid)
	}
	serverSlug := strings.TrimSpace(in.ServerSlug)
	if serverSlug == "" {
		return apperrors.Invalid(apperrors.CodeStatsIngestV1HandlerBodyInvalid, apperrors.MsgStatsIngestV1HandlerBodyInvalid)
	}
	if in.SessionSeconds <= 0 && in.MobKills <= 0 {
		return apperrors.Invalid(apperrors.CodeStatsIngestV1HandlerBodyInvalid, apperrors.MsgStatsIngestV1HandlerBodyInvalid)
	}
	if _, err := s.playerRepo.FindByID(ctx, in.PlayerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodePlayerGetV1ServiceNotFound, apperrors.MsgPlayerGetV1ServiceNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeStatsIngestV1ServiceFailed, apperrors.MsgStatsIngestV1ServiceFailed, err)
	}
	gs, err := s.gameServerRepo.GetOrCreateBySlug(ctx, serverSlug, serverSlug)
	if err != nil {
		return apperrors.InternalCause(apperrors.CodeStatsIngestV1ServiceFailed, apperrors.MsgStatsIngestV1ServiceFailed, err)
	}
	if err := s.gameServerRepo.AddPlayerStats(ctx, gs.ID, in.PlayerID, in.SessionSeconds, in.MobKills, time.Now().UTC()); err != nil {
		return apperrors.InternalCause(apperrors.CodeStatsIngestV1ServiceFailed, apperrors.MsgStatsIngestV1ServiceFailed, err)
	}
	return nil
}

func (s *Service) GetPlayerStats(ctx context.Context, playerID uuid.UUID) (*PlayerStats, error) {
	if playerID == uuid.Nil {
		return nil, apperrors.Invalid(apperrors.CodePlayerGetV1HandlerIDInvalid, apperrors.MsgPlayerGetV1HandlerIDInvalid)
	}
	if _, err := s.playerRepo.FindByID(ctx, playerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodePlayerGetV1ServiceNotFound, apperrors.MsgPlayerGetV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeStatsGetV1ServiceFailed, apperrors.MsgStatsGetV1ServiceFailed, err)
	}
	rows, err := s.gameServerRepo.ListPlayerStats(ctx, playerID)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeStatsGetV1ServiceFailed, apperrors.MsgStatsGetV1ServiceFailed, err)
	}
	out := &PlayerStats{PlayerID: playerID, ByServer: make([]ServerStat, 0, len(rows))}
	for _, row := range rows {
		out.TotalPlayTimeSecs += row.PlayTimeSecs
		out.TotalMobKills += row.MobKills
		out.ByServer = append(out.ByServer, ServerStat{
			ServerSlug:   row.ServerSlug,
			PlayTimeSecs: row.PlayTimeSecs,
			MobKills:     row.MobKills,
		})
	}
	return out, nil
}

func (s *Service) GetHUD(ctx context.Context, playerID uuid.UUID) (*HUD, error) {
	if playerID == uuid.Nil {
		return nil, apperrors.Invalid(apperrors.CodePlayerGetV1HandlerIDInvalid, apperrors.MsgPlayerGetV1HandlerIDInvalid)
	}
	player, err := s.playerRepo.FindByID(ctx, playerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodePlayerGetV1ServiceNotFound, apperrors.MsgPlayerGetV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeStatsHUDV1ServiceFailed, apperrors.MsgStatsHUDV1ServiceFailed, err)
	}
	out := &HUD{
		Username:        player.Username,
		Status:          player.Status,
		AffiliationType: player.AffiliationType,
	}
	if player.UniversitySlug != nil {
		out.UniversitySlug = *player.UniversitySlug
	}
	if player.FacultySlug != nil {
		out.FacultySlug = *player.FacultySlug
	}
	if player.CourseSlug != nil {
		out.CourseSlug = *player.CourseSlug
	}
	if s.affiliation != nil {
		labels := s.affiliation.ResolveLabels(ctx, player)
		out.UniversityName = labels.UniversityName
		out.UniversityHex = labels.UniversityHex
		out.FacultyName = labels.FacultyName
		out.FacultyAbbr = labels.FacultyAbbr
		out.FacultyHex = labels.FacultyHex
		out.CourseName = labels.CourseName
		out.CourseAbbr = labels.CourseAbbr
		out.CourseHex = labels.CourseHex
	}
	guild, err := s.guildRepo.PlayerGuild(ctx, playerID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.InternalCause(apperrors.CodeStatsHUDV1ServiceFailed, apperrors.MsgStatsHUDV1ServiceFailed, err)
	}
	if guild != nil {
		out.GuildID = guild.ID.String()
		out.GuildName = guild.Name
		out.GuildSlug = guild.Slug
		if s.presence != nil {
			gp, err := s.presence.Guild(ctx, guild.ID)
			if err == nil && gp != nil {
				out.GuildOnlineCount = gp.OnlineCount
			}
		}
	}
	return out, nil
}
