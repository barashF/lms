package publisher

import (
	"context"

	"github.com/google/uuid"

	"github.com/barashF/lms/service-order/internal/model"
)

type outboxRepository interface {
	FetchUnprocessed(context.Context, int64) ([]*model.OutBoxMessage, error)
	MarkAsProcessed(context.Context, uuid.UUID) error
}

type kafkaPublisher interface {
	Publish(topic, key string, msg []byte, headers map[string]string) error
}
