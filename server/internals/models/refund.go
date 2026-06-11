package models

import "time"

type Refund struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	OrderID        uint       `json:"order_id" gorm:"not null;uniqueIndex"`
	PaymentID      uint       `json:"payment_id" gorm:"not null"`
	Reference      string     `json:"reference" gorm:"not null;uniqueIndex"`
	IdempotencyKey string     `json:"-" gorm:"not null;uniqueIndex"`
	Amount         float64    `json:"amount" gorm:"not null"`
	Currency       string     `json:"currency" gorm:"default:NGN"`
	Status         string     `json:"status" gorm:"default:pending"`
	Reason         string     `json:"reason"`
	ProcessedAt    *time.Time `json:"processed_at"`
	CreatedAt      time.Time  `json:"created_at"`
	// Relationships
	Order   Order   `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Payment Payment `json:"payment,omitempty" gorm:"foreignKey:PaymentID"`
}
