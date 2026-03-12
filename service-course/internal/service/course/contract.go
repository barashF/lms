package course

import (
	"context"

	"github.com/google/uuid"

	"github.com/barashF/lms/service-course/internal/model"
)

type courseRepo interface {
	Create(ctx context.Context, course model.Course) (uuid.UUID, error)
	FetchByID(ctx context.Context, id uuid.UUID) (*model.Course, error)
	FetchAll(ctx context.Context) ([]*model.Course, error)
}
