package mapper

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
)

func ConvertToCartResponse(cart *models.Cart) *dto.CartResponse {
	cartItem := make([]dto.CartItemResponse, len(cart.CartItems))
	var total float64

	for i := range cart.CartItems {
		subtotal := cart.CartItems[i].Menu.Price * float64(cart.CartItems[i].Quantity)
		total = total + subtotal
	
		cartItem[i] = dto.CartItemResponse{
			ID:       cart.CartItems[i].ID,
			MenuItemID: cart.CartItems[i].MenuID,
			MenuItem: dto.MenuResponse{
				ID:    cart.CartItems[i].Menu.ID,
				Name:  cart.CartItems[i].Menu.Name,
				Price: cart.CartItems[i].Menu.Price,
				PrepTimeMinutes: cart.CartItems[i].Menu.PrepTimeMinutes,
				SpiceLevel: cart.CartItems[i].Menu.SpiceLevel,
			},
			Quantity: cart.CartItems[i].Quantity,
			Subtotal:  subtotal,
			SpecialInstructions: cart.CartItems[i].SpecialInstructions,
		}
	}

	count := len(cartItem)

	return &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CartItems: cartItem,
		Total:     total,
		ItemCount: count,
		CreatedAt: cart.CreatedAt,
	}
}