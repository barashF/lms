package model

import (
	"time"

	"github.com/google/uuid"
)

type Enrollment struct {
	ID        uuid.UUID
	CourseID  uuid.UUID
	UserID    uuid.UUID
	Status    EnrollmentStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EnrollmentStatus string

var (
	EnrollmentStatusActive   EnrollmentStatus = "active"
	EnrollmentStatusComleted EnrollmentStatus = "completed"
	EnrollmentStatusDroped   EnrollmentStatus = "droped"
)
