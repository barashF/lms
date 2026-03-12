package record

import (
	"context"

	"github.com/google/uuid"

	"github.com/barashF/lms/service-course/internal/logger"
	"github.com/barashF/lms/service-course/internal/model"
)

type Usecase struct {
	enrollmentRepo enrollmentRepository
	courseRepo     courseRepository
	factory        EnrollmentFactory
	orderGateway   orderGateway
	logger         logger.Logger
}

func NewUsecase(enrollment enrollmentRepository, course courseRepository, factory EnrollmentFactory, order orderGateway, logger logger.Logger) *Usecase {
	return &Usecase{
		enrollmentRepo: enrollment,
		courseRepo:     course,
		factory:        factory,
		orderGateway:   order,
		logger:         logger,
	}
}

func (u *Usecase) Record(ctx context.Context, orderID uuid.UUID) (uuid.UUID, error) {
	order, err := u.orderGateway.GetOrderByID(ctx, orderID)
	if err != nil {
		u.logger.Error("Failed get order", logger.NewField("error", err))
		return uuid.Nil, err
	}
	courseID, err := uuid.Parse(order.CourseId)
	if err != nil {
		u.logger.Error("Invalid course id in order", logger.NewField("error", err))
		return uuid.Nil, err
	}
	userID, err := uuid.Parse(order.UserId)
	if err != nil {
		u.logger.Error("Invalid user id in order", logger.NewField("error", err))
		return uuid.Nil, err
	}
	course, err := u.courseRepo.FetchByID(ctx, courseID)
	if err != nil {
		u.logger.Error("Faled get course by id", logger.NewField("error", err))
		return uuid.Nil, err
	}
	checkerPayment := u.factory.GetStrategy(course.Type)
	err = checkerPayment.CanEnroll(ctx, course, userID, orderID)
	if err != nil {
		u.logger.Error("Failed check payment", logger.NewField("error", err))
		return uuid.Nil, err
	}

	enrollment := model.Enrollment{
		ID:       uuid.New(),
		CourseID: courseID,
		UserID:   uuid.MustParse(order.UserId),
		Status:   model.EnrollmentStatusActive,
	}
	enrollID, err := u.enrollmentRepo.Create(ctx, enrollment)
	if err != nil {
		u.logger.Error("Failed create enrollment", logger.NewField("error", err))
		return uuid.Nil, err
	}

	return enrollID, nil
}
