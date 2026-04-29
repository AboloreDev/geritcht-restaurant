package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        *uint          `json:"user_id" gorm:"index"`
	TableID       *uint          `json:"table_id" gorm:"index"`
	ReservationID *uint          `json:"reservation_id" gorm:"index"`
	Type          OrderType      `json:"type" gorm:"not null"`
	Status        OrderStatus    `json:"status" gorm:"default:pending"`
	TotalAmount   float64        `json:"total_amount" gorm:"not null"`
	PaymentStatus PaymentStatus  `json:"payment_status" gorm:"default:unpaid"`
	Notes         string         `json:"notes"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	// Relationships
	User        *User        `json:"user,omitempty"`
	Table       *Table       `json:"table,omitempty"`
	Reservation *Reservation `json:"reservation,omitempty"`
	OrderItems  []OrderItem  `json:"order_items,omitempty"`
	Payment     *Payment     `json:"payment,omitempty"`
}

type OrderItem struct {
	ID                  uint           `json:"id" gorm:"primaryKey"`
	OrderID             uint           `json:"order_id" gorm:"not null;index"`
	MenuItemID          uint           `json:"menu_item_id" gorm:"not null;index"`
	Quantity            int            `json:"quantity" gorm:"not null"`
	Price               float64        `json:"price" gorm:"not null"`
	SpecialInstructions string         `json:"special_instructions"`
	CreatedAt           time.Time      `json:"created_at"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`
	// Relationships
	Order Order `json:"-"`
	Menu  Menu  `json:"menu_item,omitempty"`
}

type OrderType string

const (
	OrderTypeTakeout OrderType = "takeout"
	OrderTypeDineIn  OrderType = "dine_in"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)
