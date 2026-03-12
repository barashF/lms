package course

import (
	"context"

	"github.com/google/uuid"

	"github.com/barashF/lms/service-course/internal/model"
)

type courseService interface {
	Create(context.Context, model.Course) (uuid.UUID, error)
	FetchByID(context.Context, uuid.UUID) (*model.Course, error)
	FetchAll(context.Context) ([]*model.Course, error)
}
