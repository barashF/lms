package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/barashF/lms/service-order/internal/domain/entity"
	"github.com/barashF/lms/service-order/internal/model"
)

type orderRepository interface {
	Create(context.Context, *model.Order) (uuid.UUID, error)
	FetchByID(context.Context, uuid.UUID) (*model.Order, error)
}

type outboxRepository interface {
	Create(context.Context, *model.OutBoxMessage) (uuid.UUID, error)
}

type paymentEventProducer interface {
	OrderCreated(context.Context, *entity.Order) error
}

type txManager interface {
	InTransaction(context.Context, *model.TransactionOptions, func(context.Context) error) error
}
