package order

import (
	"context"

	"github.com/barashF/lms/service-order/internal/model"
	"github.com/google/uuid"
)

type orderService interface {
	Create(context.Context, *model.Order) (uuid.UUID, error)
}
