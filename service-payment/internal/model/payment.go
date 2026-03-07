package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Payment struct {
	ID                uuid.UUID
	OrderID           uuid.UUID
	UserID            uuid.UUID
	CourseID          uuid.UUID
	Value             decimal.Decimal
	Currency          string
	Status            string
	PaymentID         uuid.UUID
	ConfirmationToken string
	IdempotencyKey    string
	Paid              bool
	GatewayPayload    json.RawMessage
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
