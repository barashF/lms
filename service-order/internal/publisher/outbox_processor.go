package publisher

import (
	"context"
	"fmt"
	"time"

	"github.com/barashF/lms/service-order/internal/logger"
	"github.com/barashF/lms/service-order/internal/model"
)

type Publisher struct {
	outboxRepo     outboxRepository
	kafkaPublisher kafkaPublisher
	logger         logger.Logger
	pollInterval   time.Duration
	batchSize      int64
}

func NewPublisher(outbox outboxRepository, publisher kafkaPublisher, logger logger.Logger, pollInterval time.Duration, batchSize int64) *Publisher {
	return &Publisher{
		outboxRepo:     outbox,
		kafkaPublisher: publisher,
		logger:         logger,
		pollInterval:   pollInterval,
		batchSize:      batchSize,
	}
}

func (p *Publisher) Start(ctx context.Context) {
	p.logger.Info("publisher started",
		logger.NewField("poll_interval", p.pollInterval),
		logger.NewField("batch_size", p.batchSize),
	)

	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("publisher stopped", logger.NewField("reason", ctx.Err()))
			return
		case <-ticker.C:
			p.processMessages(ctx)
		}
	}
}

func (p *Publisher) processMessages(ctx context.Context) {
	messages, err := p.outboxRepo.FetchUnprocessed(ctx, p.batchSize)
	if err != nil {
		p.logger.Error("failed to fetch unprocessed messages",
			logger.NewField("batch_size", p.batchSize),
			logger.NewField("error", err),
		)
		return
	}

	p.logger.Info("fetched messages for processing",
		logger.NewField("count", len(messages)),
	)

	for _, msg := range messages {
		if err = p.processMessage(ctx, msg); err != nil {
			p.logger.Warn("failed to process message, skipped",
				logger.NewField("message_id", msg.ID),
				logger.NewField("aggregate_id", msg.AggregateID),
				logger.NewField("error", err),
			)
			continue
		}
	}
}

func (p *Publisher) processMessage(ctx context.Context, msg *model.OutBoxMessage) error {
	p.logger.Debug("processing message",
		logger.NewField("message_id", msg.ID),
		logger.NewField("aggregate_type", msg.AggregateType),
		logger.NewField("aggregate_id", msg.AggregateID),
	)

	key := fmt.Sprintf("%s:%s", msg.AggregateType, msg.AggregateID)
	topic := getTopic(msg.AggregateType)
	if err := p.kafkaPublisher.Publish(topic, key, msg.Payload, map[string]string{
		"event_type": string(msg.EventType),
	}); err != nil {
		return fmt.Errorf("failed to publish to Kafka: %w", err)
	}

	p.logger.Debug("message published to Kafka",
		logger.NewField("message_id", msg.ID),
		logger.NewField("key", key),
	)

	if err := p.outboxRepo.MarkAsProcessed(ctx, msg.ID); err != nil {
		return fmt.Errorf("failed to mark as processed: %w", err)
	}

	p.logger.Info("message marked as processed",
		logger.NewField("message_id", msg.ID),
	)

	return nil
}

func getTopic(aggregateType model.AggregateType) string {
	switch aggregateType {
	case model.AggregateOrder:
		return "order.status.changed"
	default:
		return ""
	}
}
