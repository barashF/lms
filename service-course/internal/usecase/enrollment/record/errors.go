package record

import "errors"

var (
	ErrCheckPayment      = errors.New("error check payment")
	ErrSecondRecord      = errors.New("user already learning on the course")
	ErrCourseNotBeenPaid = errors.New("the course has not been paid for")
	ErrTier              = errors.New("unsuitable pricing plan")
)
