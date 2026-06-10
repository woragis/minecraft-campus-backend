package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	cityrepo "github.com/woragis/minecraft-campus-backend/server/internal/city/repository"
	gameserverrepo "github.com/woragis/minecraft-campus-backend/server/internal/gameserver/repository"
	guildrepo "github.com/woragis/minecraft-campus-backend/server/internal/guild/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	playerrepo "github.com/woragis/minecraft-campus-backend/server/internal/player/repository"
	"gorm.io/gorm"
)

type Service struct {
	repo       *cityrepo.Repository
	playerRepo *playerrepo.Repository
	guildRepo  *guildrepo.Repository
	serverRepo *gameserverrepo.Repository
}

func New(repo *cityrepo.Repository, playerRepo *playerrepo.Repository, guildRepo *guildrepo.Repository, serverRepo *gameserverrepo.Repository) *Service {
	return &Service{repo: repo, playerRepo: playerRepo, guildRepo: guildRepo, serverRepo: serverRepo}
}

type CreateInput struct {
	FounderMinecraftUUID uuid.UUID
	Name                 string
	ServerSlug           string
	World                string
	CenterX              int
	CenterZ              int
	GuildID              *uuid.UUID
}

func (s *Service) Create(ctx context.Context, in CreateInput) (*models.City, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || len(name) > 64 {
		return nil, apperrors.Invalid(apperrors.CodeCityPostV1ServiceNameInvalid, apperrors.MsgCityPostV1ServiceNameInvalid)
	}
	founder, err := s.playerRepo.FindByMinecraftUUID(ctx, in.FounderMinecraftUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeCityPostV1ServiceFounderNotFound, apperrors.MsgCityPostV1ServiceFounderNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeCityPostV1ServiceCreateFailed, apperrors.MsgCityPostV1ServiceCreateFailed, err)
	}
	if founder.Status == models.PlayerStatusProbation || founder.Status == models.PlayerStatusBanned {
		return nil, apperrors.Forbidden(apperrors.CodeCityPostV1ServiceFounderRestricted, apperrors.MsgCityPostV1ServiceFounderRestricted)
	}

	world := strings.TrimSpace(in.World)
	if world == "" {
		world = "world"
	}
	gs, err := s.serverRepo.GetOrCreateBySlug(ctx, in.ServerSlug, in.ServerSlug)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeCityPostV1ServiceCreateFailed, apperrors.MsgCityPostV1ServiceCreateFailed, err)
	}

	var guildID *uuid.UUID
	if in.GuildID != nil {
		if _, err := s.guildRepo.FindByID(ctx, *in.GuildID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apperrors.NotFound(apperrors.CodeCityPostV1ServiceGuildNotFound, apperrors.MsgCityPostV1ServiceGuildNotFound)
			}
			return nil, apperrors.InternalCause(apperrors.CodeCityPostV1ServiceCreateFailed, apperrors.MsgCityPostV1ServiceCreateFailed, err)
		}
		guildID = in.GuildID
	}

	slug, err := s.uniqueSlug(ctx, slugify(name))
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeCityPostV1ServiceCreateFailed, apperrors.MsgCityPostV1ServiceCreateFailed, err)
	}

	city := &models.City{
		ID:        uuid.New(),
		Slug:      slug,
		Name:      name,
		FounderID: founder.ID,
		GuildID:   guildID,
		ServerID:  gs.ID,
		World:     world,
		CenterX:   in.CenterX,
		CenterZ:   in.CenterZ,
	}
	if err := s.repo.Create(ctx, city); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeCityPostV1ServiceCreateFailed, apperrors.MsgCityPostV1ServiceCreateFailed, err)
	}
	return city, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.City, error) {
	city, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeCityGetV1ServiceNotFound, apperrors.MsgCityGetV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeCityGetV1ServiceLoadFailed, apperrors.MsgCityGetV1ServiceLoadFailed, err)
	}
	return s.enrich(ctx, city)
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*models.City, error) {
	city, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeCityGetV1ServiceNotFound, apperrors.MsgCityGetV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeCityGetV1ServiceLoadFailed, apperrors.MsgCityGetV1ServiceLoadFailed, err)
	}
	return s.enrich(ctx, city)
}

func (s *Service) List(ctx context.Context) ([]models.City, error) {
	cities, err := s.repo.List(ctx, 100)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeCityListV1ServiceListFailed, apperrors.MsgCityListV1ServiceListFailed, err)
	}
	for i := range cities {
		pop, _ := s.repo.Population(ctx, cities[i].ID)
		cities[i].Population = pop
	}
	return cities, nil
}

func (s *Service) enrich(ctx context.Context, city *models.City) (*models.City, error) {
	pop, err := s.repo.Population(ctx, city.ID)
	if err != nil {
		return nil, err
	}
	city.Population = pop
	return city, nil
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
	return "", fmt.Errorf("could not allocate slug for %q", base)
}

func slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = strings.ReplaceAll(s, " ", "-")
	if len(s) > 48 {
		s = s[:48]
	}
	return s
}
