package model

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID          uuid.UUID
	Title       string
	Description string
	Type        TypeCourse
	Price       int64
	Author      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TypeCourse string

var (
	CourseFree      TypeCourse = "free"
	CoursePaid      TypeCourse = "paid"
	CoursePremium   TypeCourse = "premium"
	CourseCorporate TypeCourse = "corporate"
)
