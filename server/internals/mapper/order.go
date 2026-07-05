package mapper

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
)

func OrderResponse(order *models.Order) *dto.OrderResponse {
	if order == nil {
		return nil
	}

	orderItems := make([]dto.OrderItemResponse, len(order.OrderItems))

	for i := range order.OrderItems {
		orderItems[i] = dto.OrderItemResponse{
			ID: order.OrderItems[i].ID,
			MenuItem: dto.MenuResponse{
				ID:    order.OrderItems[i].Menu.ID,
				Name:  order.OrderItems[i].Menu.Name,
				Price: order.OrderItems[i].Menu.Price,
			},
			Quantity: order.OrderItems[i].Quantity,
			Price:    order.OrderItems[i].Price,
		}
	}

	// guard Payment — it's a pointer (*Payment)
	var paymentResponse *dto.PaymentResponse
	if order.Payment != nil {
		paymentResponse = PaymentResponse(order.Payment)
	}

	// guard UserID — it's *uint
	var userID uint
	if order.UserID != nil {
		userID = *order.UserID
	}

	return &dto.OrderResponse{
		ID:     order.ID,
		UserID: &userID,
		User: &dto.UserResponse{
			ID:          order.User.ID,
			FirstName:   order.User.FirstName,
			Email:       order.User.Email,
			PhoneNumber: order.User.PhoneNumber,
		},
		Status:        string(order.Status),
		TotalAmount:   order.TotalAmount,
		OrderItems:    orderItems,
		PaymentStatus: string(order.PaymentStatus),
		Payment:       paymentResponse,
		CreatedAt:     order.CreatedAt,
		Notes:         order.Notes,
	}
}
