package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	cityrepo "github.com/woragis/minecraft-campus-backend/server/internal/city/repository"
	claimrepo "github.com/woragis/minecraft-campus-backend/server/internal/claim/repository"
	gameserverrepo "github.com/woragis/minecraft-campus-backend/server/internal/gameserver/repository"
	guildrepo "github.com/woragis/minecraft-campus-backend/server/internal/guild/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	playerrepo "github.com/woragis/minecraft-campus-backend/server/internal/player/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/territory"
	"gorm.io/gorm"
)

type PermissionResult struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason"`
	ClaimID string `json:"claimId,omitempty"`
}

type Service struct {
	repo       *claimrepo.Repository
	playerRepo *playerrepo.Repository
	cityRepo   *cityrepo.Repository
	guildRepo  *guildrepo.Repository
	serverRepo *gameserverrepo.Repository
}

func New(repo *claimrepo.Repository, playerRepo *playerrepo.Repository, cityRepo *cityrepo.Repository, guildRepo *guildrepo.Repository, serverRepo *gameserverrepo.Repository) *Service {
	return &Service{repo: repo, playerRepo: playerRepo, cityRepo: cityRepo, guildRepo: guildRepo, serverRepo: serverRepo}
}

type CreateInput struct {
	OwnerMinecraftUUID uuid.UUID
	ServerSlug         string
	World              string
	MinX, MaxX         int
	MinZ, MaxZ         int
	MinY, MaxY         int
	ZoneType           string
	CityID             *uuid.UUID
}

func (s *Service) Create(ctx context.Context, in CreateInput) (*models.Claim, error) {
	owner, err := s.playerRepo.FindByMinecraftUUID(ctx, in.OwnerMinecraftUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeClaimPostV1ServiceOwnerNotFound, apperrors.MsgClaimPostV1ServiceOwnerNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeClaimPostV1ServiceCreateFailed, apperrors.MsgClaimPostV1ServiceCreateFailed, err)
	}
	if owner.Status == models.PlayerStatusProbation {
		return nil, apperrors.Forbidden(apperrors.CodeClaimPostV1ServiceProbation, apperrors.MsgClaimPostV1ServiceProbation)
	}
	if owner.Status == models.PlayerStatusBanned {
		return nil, apperrors.Forbidden(apperrors.CodeClaimPostV1ServiceBanned, apperrors.MsgClaimPostV1ServiceBanned)
	}

	minX, maxX, minZ, maxZ := territory.NormalizeBounds(in.MinX, in.MaxX, in.MinZ, in.MaxZ)
	area := territory.HorizontalArea(minX, maxX, minZ, maxZ)
	if area <= 0 {
		return nil, apperrors.Invalid(apperrors.CodeClaimPostV1ServiceBoundsInvalid, apperrors.MsgClaimPostV1ServiceBoundsInvalid)
	}

	maxAllowed := territory.MaxClaimArea(owner.CreatedAt, owner.TrustScore)
	used, err := s.repo.TotalAreaByOwner(ctx, owner.ID)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeClaimPostV1ServiceCreateFailed, apperrors.MsgClaimPostV1ServiceCreateFailed, err)
	}
	if used+area > maxAllowed {
		return nil, apperrors.Forbidden(apperrors.CodeClaimPostV1ServiceAreaLimit, apperrors.MsgClaimPostV1ServiceAreaLimit)
	}

	zone := strings.ToLower(strings.TrimSpace(in.ZoneType))
	if zone == "" {
		zone = models.ZoneUrban
	}
	if !territory.ValidZone(zone) {
		return nil, apperrors.Invalid(apperrors.CodeClaimPostV1ServiceZoneInvalid, apperrors.MsgClaimPostV1ServiceZoneInvalid)
	}

	world := strings.TrimSpace(in.World)
	if world == "" {
		world = "world"
	}
	gs, err := s.serverRepo.GetOrCreateBySlug(ctx, in.ServerSlug, in.ServerSlug)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeClaimPostV1ServiceCreateFailed, apperrors.MsgClaimPostV1ServiceCreateFailed, err)
	}

	overlap, err := s.repo.HasOverlap(ctx, gs.ID, world, minX, maxX, minZ, maxZ, nil)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeClaimPostV1ServiceCreateFailed, apperrors.MsgClaimPostV1ServiceCreateFailed, err)
	}
	if overlap {
		return nil, apperrors.ConflictErr(apperrors.CodeClaimPostV1ServiceOverlap, apperrors.MsgClaimPostV1ServiceOverlap)
	}

	minY, maxY := in.MinY, in.MaxY
	if minY == 0 && maxY == 0 {
		minY, maxY = -64, 320
	}
	if minY > maxY {
		minY, maxY = maxY, minY
	}

	var cityID *uuid.UUID
	var guildID *uuid.UUID
	if in.CityID != nil {
		city, err := s.cityRepo.FindByID(ctx, *in.CityID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apperrors.NotFound(apperrors.CodeClaimPostV1ServiceCityNotFound, apperrors.MsgClaimPostV1ServiceCityNotFound)
			}
			return nil, apperrors.InternalCause(apperrors.CodeClaimPostV1ServiceCreateFailed, apperrors.MsgClaimPostV1ServiceCreateFailed, err)
		}
		cityID = &city.ID
		guildID = city.GuildID
	} else if g, err := s.guildRepo.PlayerGuild(ctx, owner.ID); err == nil && g != nil {
		guildID = &g.ID
	}

	claim := &models.Claim{
		ID:         uuid.New(),
		OwnerID:    owner.ID,
		CityID:     cityID,
		GuildID:    guildID,
		ServerID:   gs.ID,
		World:      world,
		MinX:       minX,
		MaxX:       maxX,
		MinY:       minY,
		MaxY:       maxY,
		MinZ:       minZ,
		MaxZ:       maxZ,
		ZoneType:   zone,
		AreaBlocks: area,
	}
	if err := s.repo.Create(ctx, claim); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeClaimPostV1ServiceCreateFailed, apperrors.MsgClaimPostV1ServiceCreateFailed, err)
	}
	if cityID != nil {
		total, _ := s.repo.SumAreaByCity(ctx, *cityID)
		_ = s.cityRepo.UpdateArea(ctx, *cityID, total)
	}
	return claim, nil
}

func (s *Service) Delete(ctx context.Context, claimID uuid.UUID, ownerMinecraftUUID uuid.UUID) error {
	claim, err := s.repo.FindByID(ctx, claimID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeClaimDeleteV1ServiceNotFound, apperrors.MsgClaimDeleteV1ServiceNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeClaimDeleteV1ServiceFailed, apperrors.MsgClaimDeleteV1ServiceFailed, err)
	}
	owner, err := s.playerRepo.FindByMinecraftUUID(ctx, ownerMinecraftUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeClaimDeleteV1ServiceOwnerNotFound, apperrors.MsgClaimDeleteV1ServiceOwnerNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeClaimDeleteV1ServiceFailed, apperrors.MsgClaimDeleteV1ServiceFailed, err)
	}
	if claim.OwnerID != owner.ID {
		return apperrors.Forbidden(apperrors.CodeClaimDeleteV1ServiceNotOwner, apperrors.MsgClaimDeleteV1ServiceNotOwner)
	}
	cityID := claim.CityID
	if err := s.repo.Delete(ctx, claimID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeClaimDeleteV1ServiceNotFound, apperrors.MsgClaimDeleteV1ServiceNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeClaimDeleteV1ServiceFailed, apperrors.MsgClaimDeleteV1ServiceFailed, err)
	}
	if cityID != nil {
		total, _ := s.repo.SumAreaByCity(ctx, *cityID)
		_ = s.cityRepo.UpdateArea(ctx, *cityID, total)
	}
	return nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Claim, error) {
	claim, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeClaimGetV1ServiceNotFound, apperrors.MsgClaimGetV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeClaimGetV1ServiceLoadFailed, apperrors.MsgClaimGetV1ServiceLoadFailed, err)
	}
	return claim, nil
}

func (s *Service) ListByCity(ctx context.Context, cityID uuid.UUID) ([]models.Claim, error) {
	if _, err := s.cityRepo.FindByID(ctx, cityID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeClaimListCityV1ServiceCityNotFound, apperrors.MsgClaimListCityV1ServiceCityNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeClaimListCityV1ServiceListFailed, apperrors.MsgClaimListCityV1ServiceListFailed, err)
	}
	claims, err := s.repo.ListByCity(ctx, cityID)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeClaimListCityV1ServiceListFailed, apperrors.MsgClaimListCityV1ServiceListFailed, err)
	}
	return claims, nil
}

func (s *Service) CheckPermission(ctx context.Context, minecraftUUID uuid.UUID, serverSlug, world string, x, z int) (*PermissionResult, error) {
	gs, err := s.serverRepo.GetOrCreateBySlug(ctx, serverSlug, serverSlug)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeClaimPermV1ServiceFailed, apperrors.MsgClaimPermV1ServiceFailed, err)
	}
	claim, err := s.repo.FindAt(ctx, gs.ID, world, x, z)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &PermissionResult{Allowed: true, Reason: "wilderness"}, nil
		}
		return nil, apperrors.InternalCause(apperrors.CodeClaimPermV1ServiceFailed, apperrors.MsgClaimPermV1ServiceFailed, err)
	}
	player, err := s.playerRepo.FindByMinecraftUUID(ctx, minecraftUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &PermissionResult{Allowed: false, Reason: "denied", ClaimID: claim.ID.String()}, nil
		}
		return nil, apperrors.InternalCause(apperrors.CodeClaimPermV1ServiceFailed, apperrors.MsgClaimPermV1ServiceFailed, err)
	}
	if claim.OwnerID == player.ID {
		return &PermissionResult{Allowed: true, Reason: "owner", ClaimID: claim.ID.String()}, nil
	}
	return &PermissionResult{Allowed: false, Reason: "denied", ClaimID: claim.ID.String()}, nil
}
