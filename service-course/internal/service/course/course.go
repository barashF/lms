package course

import (
	"context"

	"github.com/google/uuid"

	"github.com/barashF/lms/service-course/internal/model"
)

type Service struct {
	repo courseRepo
}

func NewService(c courseRepo) *Service {
	return &Service{repo: c}
}

func (s *Service) Create(ctx context.Context, course model.Course) (uuid.UUID, error) {
	course.ID = uuid.New()
	return s.repo.Create(ctx, course)
}

func (s *Service) FetchByID(ctx context.Context, id uuid.UUID) (*model.Course, error) {
	return s.repo.FetchByID(ctx, id)
}

func (s *Service) FetchAll(ctx context.Context) ([]*model.Course, error) {
	return s.repo.FetchAll(ctx)
}
