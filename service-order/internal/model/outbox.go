package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type OutBoxMessage struct {
	ID            uuid.UUID
	AggregateType AggregateType
	AggregateID   uuid.UUID
	EventType     EventType
	Payload       json.RawMessage
	ProcessedAt   *time.Time
	Processed     bool
}

type (
	AggregateType string
	EventType     string
)

const (
	AggregateOrder AggregateType = "Order"
)

const (
	EventTypeOrderCreated   EventType = "OrderCreated"
	EventTypeOrderCancelled EventType = "OrderCancelled"
)

func NewOutboxMessage(aggregateType AggregateType, aggregateID uuid.UUID, eventType EventType, payload any) (*OutBoxMessage, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &OutBoxMessage{
		ID:            uuid.New(),
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		Payload:       payloadJSON,
	}, nil
}
