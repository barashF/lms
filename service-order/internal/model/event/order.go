package event

import (
	"time"

	"github.com/google/uuid"
)

type EventCreatedOrder struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CourseID  uuid.UUID `json:"course_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EventCancelledOrder struct {
	OrderID     uuid.UUID `json:"order_id"`
	CancelledAt time.Time `json:"cancelled_at"`
}
