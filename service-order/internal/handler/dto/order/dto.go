package order

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CourseID  uuid.UUID `json:"course_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	CourseID uuid.UUID `json:"course_id"`
}

type UpdateRequest Order
