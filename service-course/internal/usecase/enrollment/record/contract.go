package record

import (
	"context"

	"github.com/google/uuid"

	"github.com/barashF/lms/service-course/internal/model"
	"github.com/barashF/lms/service-course/proto/orders"
	payments "github.com/barashF/lms/service-course/proto/payment"
	"github.com/barashF/lms/service-course/proto/profiles"
)

type paymentGateway interface {
	GetPaymentByOrderID(context.Context, uuid.UUID) (*payments.Payment, error)
}

type orderGateway interface {
	GetOrderByID(context.Context, uuid.UUID) (*orders.Order, error)
}

type profileGateway interface {
	GetProfileByUserID(context.Context, uuid.UUID) (*profiles.Profile, error)
}

type enrollmentRepository interface {
	Create(context.Context, model.Enrollment) (uuid.UUID, error)
	FetchByCourseIDUserID(ctx context.Context, courseID uuid.UUID, userID uuid.UUID) (*model.Enrollment, error)
}

type courseRepository interface {
	FetchByID(context.Context, uuid.UUID) (*model.Course, error)
}

type EnrollmentStrategy interface {
	CanEnroll(ctx context.Context, course *model.Course, userID uuid.UUID, orderID uuid.UUID) error
}

type PaymentChecker interface {
	CheckPayment(context.Context, uuid.UUID) error
}

type PaymentCheckerFactory interface {
	GetPaymentChecker(int64) PaymentChecker
}

type EnrollmentFactory interface {
	GetStrategy(model.TypeCourse) EnrollmentStrategy
}
