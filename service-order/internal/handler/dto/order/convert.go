package order

import "github.com/barashF/lms/service-order/internal/model"

func (r CreateRequest) ToModel() *model.Order {
	return &model.Order{
		UserID:   r.UserID,
		CourseID: r.CourseID,
	}
}

func (r UpdateRequest) ToModel() *model.Order {
	return &model.Order{
		ID:       r.ID,
		CourseID: r.CourseID,
		UserID:   r.UserID,
		Status:   model.OrderStatus(r.Status),
	}
}

func ModelToResponse(order model.Order) *Order {
	return &Order{
		ID:        order.ID,
		UserID:    order.UserID,
		CourseID:  order.CourseID,
		Status:    string(order.Status),
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}
