package validation

import (
	dto "github.com/barashF/lms/service-course/internal/handler/dto/course"
	"github.com/barashF/lms/service-course/internal/model"
)

func ValidateRequest[T any](req *T) error {
	switch v := any(req).(type) {
	case *dto.CreateRequest:
		return validateCreate(v)
	case *dto.UpdateRequest:
		return validateUpdate(v)
	default:
		return nil
	}
}

func validateCreate(req *dto.CreateRequest) error {
	if req.Author == "" || req.Description == "" || req.Title == "" || req.Type == "" {
		return model.ErrMissingRequiredFields
	}

	if err := validateType(req.Type); err != nil {
		return err
	}

	if req.Type != string(model.CourseFree) && req.Price == 0 {
		return model.ErrNonFreeType
	}

	return nil
}

func validateUpdate(req *dto.UpdateRequest) error {
	if req.Author == "" && req.Description == "" && req.Price == 0 && req.Title == "" {
		return model.ErrMissingRequiredFields
	}

	return nil
}

func validateType(courseType string) error {
	switch courseType {
	case string(model.CourseFree), string(model.CoursePaid), string(model.CoursePremium), string(model.CourseCorporate):
		return nil
	default:
		return model.ErrInvalidTypeCourse
	}
}
