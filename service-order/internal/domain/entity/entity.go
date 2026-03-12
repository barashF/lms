package entity

import "time"

type Order struct {
	ID        string
	UserID    string
	CourseID  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
