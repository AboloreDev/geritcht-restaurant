package dto

import "time"

type CreateTakeoutOrderRequest struct {
	Notes string `json:"notes"`
}

type CreateDineInOrderRequest struct {
	TableID       uint                     `json:"table_id" binding:"required"`
	ReservationID *uint                    `json:"reservation_id"`
	Items         []CreateOrderItemRequest `json:"items" binding:"required,min=1"`
	Notes         string                   `json:"notes"`
}

type CreateOrderItemRequest struct {
	MenuItemID          uint   `json:"menu_item_id" binding:"required"`
	Quantity            int    `json:"quantity" binding:"required,gt=0"`
	SpecialInstructions string `json:"special_instructions"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type OrderFilterRequest struct {
	Status   string `form:"status"`
	Type     string `form:"type"`
	Date     string `form:"date"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

type OrderItemResponse struct {
	ID                  uint         `json:"id"`
	MenuItemID          uint         `json:"menu_item_id"`
	MenuItem            MenuResponse `json:"menu_item"`
	Quantity            int          `json:"quantity"`
	Price               float64      `json:"price"`
	Subtotal            float64      `json:"subtotal"`
	SpecialInstructions string       `json:"special_instructions"`
}

type OrderResponse struct {
	ID            uint                `json:"id"`
	UserID        *uint               `json:"user_id"`
	User          *UserResponse       `json:"user,omitempty"`
	TableID       *uint               `json:"table_id"`
	ReservationID *uint               `json:"reservation_id"`
	Type          string              `json:"type"`
	Status        string              `json:"status"`
	TotalAmount   float64             `json:"total_amount"`
	PaymentStatus string              `json:"payment_status"`
	Notes         string              `json:"notes"`
	OrderItems    []OrderItemResponse `json:"order_items"`
	Payment       *PaymentResponse    `json:"payment,omitempty"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}

type OrderListResponse struct {
	Orders   []OrderResponse `json:"orders"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}
