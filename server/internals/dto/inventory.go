package dto

import "time"

type CreateIngredientRequest struct {
	Name         string  `json:"name" binding:"required"`
	Unit         string  `json:"unit" binding:"required"`
	CurrentStock float64 `json:"current_stock"`
	MinThreshold float64 `json:"min_threshold" binding:"required"`
}

type UpdateIngredientRequest struct {
	Name         string  `json:"name"`
	Unit         string  `json:"unit"`
	MinThreshold float64 `json:"min_threshold"`
}

type UpdateStockRequest struct {
	Quantity float64 `json:"quantity" binding:"required"`
	Reason   string  `json:"reason" binding:"required"`
	Type     string  `json:"type" binding:"required,oneof=in out waste"`
}

type LinkIngredientRequest struct {
	IngredientID uint    `json:"ingredient_id" binding:"required"`
	Quantity     float64 `json:"quantity" binding:"required,gt=0"`
}

type IngredientResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Unit         string    `json:"unit"`
	CurrentStock float64   `json:"current_stock"`
	MinThreshold float64   `json:"min_threshold"`
	IsLow        bool      `json:"is_low"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type StockMovementResponse struct {
	ID           uint               `json:"id"`
	IngredientID uint               `json:"ingredient_id"`
	Ingredient   IngredientResponse `json:"ingredient"`
	Type         string             `json:"type"`
	Quantity     float64            `json:"quantity"`
	Reason       string             `json:"reason"`
	CreatedBy    uint               `json:"created_by"`
	User         UserResponse       `json:"user"`
	CreatedAt    time.Time          `json:"created_at"`
}

type MenuItemIngredientResponse struct {
	IngredientID uint               `json:"ingredient_id"`
	Ingredient   IngredientResponse `json:"ingredient"`
	Quantity     float64            `json:"quantity"`
}

type InventoryAlertResponse struct {
	LowStockIngredients []IngredientResponse `json:"low_stock_ingredients"`
	OutOfStockItems     []MenuResponse       `json:"out_of_stock_items"`
	TotalLowStock       int                  `json:"total_low_stock"`
	TotalOutOfStock     int                  `json:"total_out_of_stock"`
}

type ThresholdRequest struct {
	Threshold float64 `json:"min_threshold" binding:"required"`
}
type UpdateLinkItemRequest struct {
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
}
