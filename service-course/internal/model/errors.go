package model

import "errors"

var (
	ErrNoCourseFound         = errors.New("no course found")
	ErrNoEnrollmentFound     = errors.New("no course found")
	ErrMissingRequiredFields = errors.New("missing required fields")
	ErrInvalidTypeCourse     = errors.New("invalid type course")
	ErrNonFreeType           = errors.New("a course with a non-free type must have a price")
)
