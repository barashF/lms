package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/barashF/lms/service-order/internal/logger"
	"github.com/barashF/lms/service-order/internal/model"
	"github.com/barashF/lms/service-order/internal/model/event"
)

type Service struct {
	orderRepo  orderRepository
	outboxRepo outboxRepository
	txManager  txManager
	logger     logger.Logger
}

func NewService(o orderRepository, outboxRepo outboxRepository, m txManager, logger logger.Logger) *Service {
	return &Service{
		orderRepo:  o,
		outboxRepo: outboxRepo,
		txManager:  m,
		logger:     logger,
	}
}

func (s *Service) Create(ctx context.Context, order *model.Order) (uuid.UUID, error) {
	order.ID = uuid.New()
	order.Status = model.OrderPending

	err := s.txManager.InTransaction(ctx, nil, func(dbContext context.Context) error {
		_, err := s.orderRepo.Create(dbContext, order)
		if err != nil {
			return err
		}

		event := event.EventCreatedOrder{
			ID:        order.ID,
			UserID:    order.UserID,
			CourseID:  order.CourseID,
			Status:    string(order.Status),
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.UpdatedAt,
		}
		message, err := model.NewOutboxMessage(model.AggregateOrder, order.ID, model.EventTypeOrderCreated, event)
		if err != nil {
			return err
		}
		_, err = s.outboxRepo.Create(dbContext, message)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return order.ID, nil
}
