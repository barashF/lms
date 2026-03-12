package course

import "github.com/barashF/lms/service-course/internal/model"

func (r CreateRequest) ToModel() model.Course {
	return model.Course{
		Title:       r.Title,
		Description: r.Description,
		Price:       r.Price,
		Type:        model.TypeCourse(r.Type),
		Author:      r.Author,
	}
}

func (r UpdateRequest) ToModel() model.Course {
	return model.Course{
		ID:          r.ID,
		Title:       r.Title,
		Description: r.Description,
		Price:       r.Price,
		Type:        model.TypeCourse(r.Type),
		Author:      r.Author,
	}
}

func ModelToResponse(course model.Course) Course {
	return Course{
		ID:          course.ID,
		Title:       course.Title,
		Description: course.Description,
		Price:       course.Price,
		Type:        string(course.Type),
		Author:      course.Author,
	}
}
