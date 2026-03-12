package test

import (
	"context"
	"fmt"

	"github.com/barashF/lms/service-course/internal/model"
	"github.com/barashF/lms/service-course/proto/orders"
	payments "github.com/barashF/lms/service-course/proto/payment"
	"github.com/barashF/lms/service-course/proto/profiles"
	"github.com/google/uuid"
)

type MockCourseRepository struct{}

func (r *MockCourseRepository) FetchByID(ctx context.Context, id uuid.UUID) (*model.Course, error) {
	switch id {
	case uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"):
		return &model.Course{
			ID:    uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"),
			Price: 100,
			Type:  model.CoursePaid,
		}, nil
	default:
		return nil, model.ErrNoCourseFound
	}
}

type MockOrderGateway struct{}

func (o *MockOrderGateway) GetOrderByID(ctx context.Context, id uuid.UUID) (*orders.Order, error) {
	switch id {
	case uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"):
		return &orders.Order{
			Id:       "c8bfd113-1023-49eb-a077-6820dd7e7c9b",
			UserId:   "c8bfd113-1023-49eb-a077-6820dd7e7c9b",
			CourseId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b",
		}, nil
	default:
		return nil, fmt.Errorf("Failed get order")
	}
}

type MockPaymentGateway struct{}

func (p *MockPaymentGateway) GetPaymentByOrderID(ctx context.Context, id uuid.UUID) (*payments.Payment, error) {
	switch id {
	case uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"):
		return &payments.Payment{
			Id:      "c8bfd113-1023-49eb-a077-6820dd7e7c9b",
			UserId:  "c8bfd113-1023-49eb-a077-6820dd7e7c9b",
			OrderId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b",
			Status:  "completed",
		}, nil
	default:
		return nil, fmt.Errorf("Failed get Payment")
	}
}

type MockProfileGateway struct{}

func (p *MockProfileGateway) GetProfileByUserID(ctx context.Context, id uuid.UUID) (*profiles.Profile, error) {
	switch id {
	case uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"):
		return &profiles.Profile{
			UserId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b",
			Tier:   "premium",
		}, nil
	default:
		return nil, fmt.Errorf("Failed get Profile")
	}
}

type MockEnrollmentRepository struct{}

func (r *MockEnrollmentRepository) Create(ctx context.Context, enrollment model.Enrollment) (uuid.UUID, error) {
	return uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"), nil
}

func (r *MockEnrollmentRepository) FetchByCourseIDUserID(ctx context.Context, courseID uuid.UUID, userID uuid.UUID) (*model.Enrollment, error) {
	switch courseID {
	case uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"):
		return nil, model.ErrNoEnrollmentFound
	default:
		return &model.Enrollment{}, nil
	}
}
