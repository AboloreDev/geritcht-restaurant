package mapper

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
)

func OrderResponse(order *models.Order) *dto.OrderResponse {
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

	return &dto.OrderResponse{
		ID:     order.ID,
		UserID: order.UserID,
		User: &dto.UserResponse{
			FirstName:   order.User.FirstName,
			Email:       order.User.Email,
			PhoneNumber: order.User.PhoneNumber,
		},
		Status:        string(order.Status),
		TotalAmount:   order.TotalAmount,
		OrderItems:    orderItems,
		PaymentStatus: string(order.PaymentStatus),
		Payment: &dto.PaymentResponse{
			ID:     order.Payment.ID,
			Amount: order.Payment.Amount,
			Status: string(order.Payment.Status),
		},
		CreatedAt: order.CreatedAt,
		Notes:     order.Notes,
	}
}
