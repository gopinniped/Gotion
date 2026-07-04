package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gopinniped/gotion/internal/modules/slots/entity"
)

type SlotRepositoryI interface {
	Create(ctx context.Context, userID uuid.UUID, slot *entity.Slot) (*entity.Slot, error)
	Update(ctx context.Context, slot *entity.Slot) (*entity.Slot, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Slot, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Slot, error)
	Delete(ctx context.Context, id uuid.UUID) error
	IsOwner(ctx context.Context, userID, slotID uuid.UUID) (bool, error)

	GetByDateRange(ctx context.Context, userID uuid.UUID, start time.Time, end time.Time) ([]*entity.Slot, error)
}
