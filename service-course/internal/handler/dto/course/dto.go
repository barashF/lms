package course

import "github.com/google/uuid"

type Course struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	Type        string    `json:"type"`
	Author      string    `json:"author"`
}

type CreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Type        string `json:"type"`
	Author      string `json:"author"`
}

type UpdateRequest Course
