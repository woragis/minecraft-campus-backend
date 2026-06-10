package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	auditrepo "github.com/woragis/minecraft-campus-backend/server/internal/audit/repository"
	gameserverrepo "github.com/woragis/minecraft-campus-backend/server/internal/gameserver/repository"
	"github.com/woragis/minecraft-campus-backend/server/internal/models"
	playerrepo "github.com/woragis/minecraft-campus-backend/server/internal/player/repository"
	rollbackrepo "github.com/woragis/minecraft-campus-backend/server/internal/rollback/repository"
	trustsvc "github.com/woragis/minecraft-campus-backend/server/internal/trust/service"
	"gorm.io/gorm"
)

const maxRollbackWindow = 24 * time.Hour

type Service struct {
	repo       *rollbackrepo.Repository
	auditRepo  *auditrepo.Repository
	playerRepo *playerrepo.Repository
	serverRepo *gameserverrepo.Repository
	trustSvc   *trustsvc.Service
}

func New(
	repo *rollbackrepo.Repository,
	auditRepo *auditrepo.Repository,
	playerRepo *playerrepo.Repository,
	serverRepo *gameserverrepo.Repository,
	trustSvc *trustsvc.Service,
) *Service {
	return &Service{
		repo:       repo,
		auditRepo:  auditRepo,
		playerRepo: playerRepo,
		serverRepo: serverRepo,
		trustSvc:   trustSvc,
	}
}

type CreateInput struct {
	TargetMinecraftUUID uuid.UUID
	ActorMinecraftUUID  uuid.UUID
	ServerSlug          string
	World               string
	WindowMinutes       int
}

type CreateResult struct {
	Rollback *models.Rollback `json:"rollback"`
}

func (s *Service) Create(ctx context.Context, in CreateInput) (*CreateResult, error) {
	if in.WindowMinutes <= 0 || in.WindowMinutes > 1440 {
		return nil, apperrors.Invalid(apperrors.CodeRollbackPostV1ServiceWindowInvalid, apperrors.MsgRollbackPostV1ServiceWindowInvalid)
	}
	target, err := s.playerRepo.FindByMinecraftUUID(ctx, in.TargetMinecraftUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeRollbackPostV1ServiceTargetNotFound, apperrors.MsgRollbackPostV1ServiceTargetNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeRollbackPostV1ServiceCreateFailed, apperrors.MsgRollbackPostV1ServiceCreateFailed, err)
	}
	actor, err := s.playerRepo.FindByMinecraftUUID(ctx, in.ActorMinecraftUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeRollbackPostV1ServiceActorNotFound, apperrors.MsgRollbackPostV1ServiceActorNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeRollbackPostV1ServiceCreateFailed, apperrors.MsgRollbackPostV1ServiceCreateFailed, err)
	}

	world := strings.TrimSpace(in.World)
	if world == "" {
		world = "world"
	}
	gs, err := s.serverRepo.GetOrCreateBySlug(ctx, in.ServerSlug, in.ServerSlug)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeRollbackPostV1ServiceCreateFailed, apperrors.MsgRollbackPostV1ServiceCreateFailed, err)
	}

	toAt := time.Now().UTC()
	fromAt := toAt.Add(-time.Duration(in.WindowMinutes) * time.Minute)
	if toAt.Sub(fromAt) > maxRollbackWindow {
		fromAt = toAt.Add(-maxRollbackWindow)
	}

	events, err := s.auditRepo.ListForRollback(ctx, target.ID, gs.ID, world, fromAt, toAt)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeRollbackPostV1ServiceCreateFailed, apperrors.MsgRollbackPostV1ServiceCreateFailed, err)
	}

	items := buildRollbackItems(events)
	rollbackID := uuid.New()
	actorID := actor.ID
	rb := &models.Rollback{
		ID:             rollbackID,
		TargetPlayerID: target.ID,
		ActorPlayerID:  &actorID,
		ServerID:       gs.ID,
		World:          world,
		FromAt:         fromAt,
		ToAt:           toAt,
		Status:         models.RollbackStatusPending,
		ItemCount:      len(items),
	}
	for i := range items {
		items[i].ID = uuid.New()
		items[i].RollbackID = rollbackID
	}
	if err := s.repo.Create(ctx, rb, items); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeRollbackPostV1ServiceCreateFailed, apperrors.MsgRollbackPostV1ServiceCreateFailed, err)
	}
	if len(items) > 0 {
		_, _, _ = s.trustSvc.RecordEvent(ctx, target.ID, models.TrustEventRollbackApplied, "rollback created", &actorID)
	}
	return &CreateResult{Rollback: rb}, nil
}

func buildRollbackItems(events []models.AuditEvent) []models.RollbackItem {
	items := make([]models.RollbackItem, 0, len(events))
	for _, ev := range events {
		if ev.BlockX == nil || ev.BlockY == nil || ev.BlockZ == nil {
			continue
		}
		item := models.RollbackItem{
			AuditEventID: &ev.ID,
			BlockX:       *ev.BlockX,
			BlockY:       *ev.BlockY,
			BlockZ:       *ev.BlockZ,
			World:        ev.World,
		}
		switch ev.EventType {
		case models.AuditEventBlockPlace:
			item.Action = models.RollbackActionRemove
			item.BlockType = ev.BlockType
		case models.AuditEventBlockBreak:
			if ev.BlockType == "" {
				continue
			}
			item.Action = models.RollbackActionRestore
			item.BlockType = ev.BlockType
		default:
			continue
		}
		items = append(items, item)
	}
	return items
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Rollback, error) {
	rb, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeRollbackGetV1ServiceNotFound, apperrors.MsgRollbackGetV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeRollbackGetV1ServiceLoadFailed, apperrors.MsgRollbackGetV1ServiceLoadFailed, err)
	}
	return rb, nil
}

func (s *Service) ListItems(ctx context.Context, id uuid.UUID) ([]models.RollbackItem, error) {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeRollbackItemsV1ServiceNotFound, apperrors.MsgRollbackItemsV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeRollbackItemsV1ServiceListFailed, apperrors.MsgRollbackItemsV1ServiceListFailed, err)
	}
	items, err := s.repo.ListItems(ctx, id)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeRollbackItemsV1ServiceListFailed, apperrors.MsgRollbackItemsV1ServiceListFailed, err)
	}
	return items, nil
}

func (s *Service) Complete(ctx context.Context, id uuid.UUID, appliedCount int) (*models.Rollback, error) {
	rb, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeRollbackCompleteV1ServiceNotFound, apperrors.MsgRollbackCompleteV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeRollbackCompleteV1ServiceFailed, apperrors.MsgRollbackCompleteV1ServiceFailed, err)
	}
	status := models.RollbackStatusCompleted
	if appliedCount < rb.ItemCount {
		status = models.RollbackStatusFailed
	}
	if err := s.repo.MarkCompleted(ctx, id, appliedCount, status); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeRollbackCompleteV1ServiceNotFound, apperrors.MsgRollbackCompleteV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeRollbackCompleteV1ServiceFailed, apperrors.MsgRollbackCompleteV1ServiceFailed, err)
	}
	rb.Status = status
	rb.AppliedCount = appliedCount
	now := time.Now().UTC()
	rb.CompletedAt = &now
	return rb, nil
}

func (s *Service) MarkApplying(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.UpdateStatus(ctx, id, models.RollbackStatusApplying); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeRollbackGetV1ServiceNotFound, apperrors.MsgRollbackGetV1ServiceNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeRollbackGetV1ServiceLoadFailed, apperrors.MsgRollbackGetV1ServiceLoadFailed, err)
	}
	return nil
}
