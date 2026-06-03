package models

import (
	"time"

	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentStatusUnpaid   PaymentStatus = "unpaid"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
	PaymentStatusPending  PaymentStatus = "pending"
)

type Payment struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	OrderID           uint           `json:"order_id" gorm:"not null;uniqueIndex"`
	UserID            uint           `json:"user_id" gorm:"not null;index"`
	Reference         string         `json:"reference" gorm:"not null;uniqueIndex"`
	IdempotencyKey    string         `json:"-" gorm:"not null;uniqueIndex"`
	Amount            float64        `json:"amount" gorm:"not null"`
	Currency          string         `json:"currency" gorm:"default:NGN"`
	Status            PaymentStatus  `json:"status" gorm:"default:unpaid"`
	Provider          string         `json:"provider" gorm:"default:paystack"`
	ProviderReference string         `json:"provider_reference"`
	FailureReason     string         `json:"failure_reason"`
	PaidAt            *time.Time     `json:"paid_at"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
	// Relationships
	Order Order `json:"order,omitempty"`
	User  User  `json:"user,omitempty"`
}
