package record

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/barashF/lms/service-course/internal/model"
)

type Validator struct {
	enrollmentRepository enrollmentRepository
}

func (v *Validator) Validate(ctx context.Context, courseID uuid.UUID, userID uuid.UUID) error {
	_, err := v.enrollmentRepository.FetchByCourseIDUserID(ctx, courseID, userID)
	if err != nil {
		if errors.Is(err, model.ErrNoEnrollmentFound) {
			return nil
		}
		return err
	}
	return ErrSecondRecord
}

type BaseStrategy struct {
	validator Validator
}

func (s *BaseStrategy) CanEnroll(ctx context.Context, course *model.Course, userID uuid.UUID, orderID uuid.UUID) error {
	return s.validator.Validate(ctx, course.ID, userID)
}

type FreeStrategy struct {
	validator Validator
}

func (s *FreeStrategy) CanEnroll(ctx context.Context, course *model.Course, userID uuid.UUID, orderID uuid.UUID) error {
	return s.validator.Validate(ctx, course.ID, userID)
}

type PaidStrategy struct {
	validator      Validator
	paymentGateway paymentGateway
}

func (s *PaidStrategy) CanEnroll(ctx context.Context, course *model.Course, userID uuid.UUID, orderID uuid.UUID) error {
	err := s.validator.Validate(ctx, course.ID, userID)
	if err != nil {
		return err
	}

	checker := StandardPaymentChecker{paymentGateway: s.paymentGateway}
	err = checker.CheckPayment(ctx, orderID)
	return err
}

type PremiumStrategy struct {
	validator      Validator
	paymentGateway paymentGateway
	profileGateway profileGateway
}

func (s *PremiumStrategy) CanEnroll(ctx context.Context, course *model.Course, userID uuid.UUID, orderID uuid.UUID) error {
	err := s.validator.Validate(ctx, course.ID, userID)
	if err != nil {
		return err
	}

	factory := NewCheckerFactory(s.paymentGateway)
	checker := factory.GetPaymentChecker(course.Price)
	err = checker.CheckPayment(ctx, orderID)
	if err != nil {
		return err
	}

	profile, err := s.profileGateway.GetProfileByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if profile.Tier != "premium" {
		return ErrTier
	}
	return nil
}

type NOPPaymentChecker struct{}

func (c *NOPPaymentChecker) CheckPayment(ctx context.Context, orderID uuid.UUID) error {
	return nil
}

type StandardPaymentChecker struct {
	paymentGateway paymentGateway
}

func (c *StandardPaymentChecker) CheckPayment(ctx context.Context, orderID uuid.UUID) error {
	payment, err := c.paymentGateway.GetPaymentByOrderID(ctx, orderID)
	if err != nil {
		return err
	}
	if payment.Status == "completed" {
		return nil
	}

	return ErrCourseNotBeenPaid
}

type CheckerFactory struct {
	paymentGateway paymentGateway
}

func NewCheckerFactory(p paymentGateway) *CheckerFactory {
	return &CheckerFactory{paymentGateway: p}
}

func (c *CheckerFactory) GetPaymentChecker(price int64) PaymentChecker {
	switch {
	case price > 0:
		return &StandardPaymentChecker{paymentGateway: c.paymentGateway}
	default:
		return &NOPPaymentChecker{}
	}
}
