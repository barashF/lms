package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	CourseID  uuid.UUID
	Status    OrderStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderStatus string

var (
	OrderPending   OrderStatus = "pending"
	OrderPaid      OrderStatus = "paid"
	OrderCompleted OrderStatus = "completed"
	OrderCancelled OrderStatus = "cancelled"
)
