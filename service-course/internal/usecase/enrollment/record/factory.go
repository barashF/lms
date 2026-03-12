package record

import "github.com/barashF/lms/service-course/internal/model"

type Factory struct {
	courseRepo     courseRepository
	enrollmentRepo enrollmentRepository
	payment        paymentGateway
	order          orderGateway
	profile        profileGateway
	validator      *Validator
}

func NewFactory(e enrollmentRepository, p paymentGateway, o orderGateway, pr profileGateway) *Factory {
	return &Factory{
		enrollmentRepo: e,
		payment:        p,
		order:          o,
		profile:        pr,
		validator:      &Validator{enrollmentRepository: e},
	}
}

func (f *Factory) GetStrategy(typeCourse model.TypeCourse) EnrollmentStrategy {
	switch typeCourse {
	case model.CourseFree:
		return &FreeStrategy{validator: *f.validator}
	case model.CoursePaid:
		return &PaidStrategy{
			validator:      *f.validator,
			paymentGateway: f.payment,
		}
	case model.CoursePremium:
		return &PremiumStrategy{
			validator:      *f.validator,
			paymentGateway: f.payment,
			profileGateway: f.profile,
		}
	default:
		return &BaseStrategy{validator: *f.validator}
	}
}
