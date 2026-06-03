package dto

import "time"

type AddToCartRequest struct {
	MenuItemID          uint   `json:"menu_item_id" binding:"required"`
	Quantity            int    `json:"quantity" binding:"required,gt=0"`
	SpecialInstructions string `json:"special_instructions"`
}

type UpdateCartItemRequest struct {
	Quantity            int    `json:"quantity" binding:"required,gt=0"`
	SpecialInstructions string `json:"special_instructions"`
}

type CartItemResponse struct {
	ID                  uint         `json:"id"`
	MenuItemID          uint         `json:"menu_item_id"`
	MenuItem            MenuResponse `json:"menu_item"`
	Quantity            int          `json:"quantity"`
	SpecialInstructions string       `json:"special_instructions"`
	Subtotal            float64      `json:"subtotal"`
	CreatedAt           time.Time    `json:"created_at"`
}

type CartResponse struct {
	ID        uint               `json:"id"`
	UserID    uint               `json:"user_id"`
	CartItems []CartItemResponse `json:"cart_items"`
	Total     float64            `json:"total"`
	ItemCount int                `json:"item_count"`
	CreatedAt time.Time          `json:"created_at"`
}
