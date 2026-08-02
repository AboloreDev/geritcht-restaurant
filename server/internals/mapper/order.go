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
		images := []dto.MenuImageResponse{}

		if len(order.OrderItems[i].Menu.Images) > 0 {
			images = append(images, dto.MenuImageResponse{
				ID:      order.OrderItems[i].Menu.Images[0].ID,
				URL:     order.OrderItems[i].Menu.Images[0].URL,
				AltText: order.OrderItems[i].Menu.Images[0].AltText,
			})
		}

		orderItems[i] = dto.OrderItemResponse{
			ID: order.OrderItems[i].ID,
			MenuItem: dto.MenuResponse{
				ID:     order.OrderItems[i].Menu.ID,
				Name:   order.OrderItems[i].Menu.Name,
				Price:  order.OrderItems[i].Menu.Price,
				Images: images,
			},
			Quantity: order.OrderItems[i].Quantity,
			Price:    order.OrderItems[i].Price,
		}
	}

	var paymentResponse *dto.PaymentResponse
	if order.Payment != nil {
		paymentResponse = PaymentResponse(order.Payment)
	}

	var userResponse *dto.UserResponse
	if order.UserID != nil {
		userResponse = &dto.UserResponse{
			ID:          *order.UserID,
			FirstName:   order.User.FirstName,
			LastName:    order.User.LastName,
			Email:       order.User.Email,
			PhoneNumber: order.User.PhoneNumber,
			Role:        string(order.User.Role),
			IsActive:    order.User.IsActive,
		}
	}

	return &dto.OrderResponse{
		ID:            order.ID,
		UserID:        order.UserID,
		User:          userResponse,
		Status:        string(order.Status),
		TotalAmount:   order.TotalAmount,
		OrderItems:    orderItems,
		PaymentStatus: string(order.PaymentStatus),
		Payment:       paymentResponse,
		CreatedAt:     order.CreatedAt,
		Notes:         order.Notes,
		Type:          string(order.Type),
	}
}
