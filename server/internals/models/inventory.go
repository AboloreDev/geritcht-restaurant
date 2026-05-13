package models

import (
	"time"

	"gorm.io/gorm"
)

type Ingredient struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"not null;uniqueIndex"`
	Unit         string         `json:"unit" gorm:"not null"`
	CurrentStock float64        `json:"current_stock" gorm:"default:0"`
	MinThreshold float64        `json:"min_threshold" gorm:"default:0"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	// Relationships
	StockMovements      []StockMovement      `json:"-"`
	MenuItemIngredients []MenuItemIngredient `json:"-"`
}

type MenuItemIngredient struct {
	MenuID       uint      `json:"menu_item_id" gorm:"primaryKey"`
	IngredientID uint      `json:"ingredient_id" gorm:"primaryKey"`
	Quantity     float64   `json:"quantity" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	// Relationships
	Menu       Menu       `json:"menu_item,omitempty"`
	Ingredient Ingredient `json:"ingredient,omitempty"`
}

type StockMovement struct {
	ID           uint              `json:"id" gorm:"primaryKey"`
	IngredientID uint              `json:"ingredient_id" gorm:"not null;index"`
	Type         StockMovementType `json:"type" gorm:"not null"`
	Quantity     float64           `json:"quantity" gorm:"not null"`
	Reason       string            `json:"reason"`
	CreatedBy    uint              `json:"created_by" gorm:"not null;index"`
	CreatedAt    time.Time         `json:"created_at"`
	// Relationships
	Ingredient Ingredient `json:"ingredient,omitempty"`
	User       User       `json:"user,omitempty" gorm:"foreignKey:CreatedBy"`
}

type StockMovementType string

const (
	StockMovementIn    StockMovementType = "in"
	StockMovementOut   StockMovementType = "out"
	StockMovementWaste StockMovementType = "waste"
)
