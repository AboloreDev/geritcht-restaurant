package dto

import "time"

type InitializePaymentRequest struct {
	OrderID uint `json:"order_id" binding:"required"`
}

type VerifyPaymentRequest struct {
	Reference string `json:"reference" binding:"required"`
}

type PaymentResponse struct {
	ID                uint       `json:"id"`
	OrderID           uint       `json:"order_id"`
	Reference         string     `json:"reference"`
	Amount            float64    `json:"amount"`
	Currency          string     `json:"currency"`
	Status            string     `json:"status"`
	Provider          string     `json:"provider"`
	ProviderReference string     `json:"provider_reference"`
	FailureReason     string     `json:"failure_reason,omitempty"`
	PaidAt            *time.Time `json:"paid_at"`
	CreatedAt         time.Time  `json:"created_at"`
}

type InitializePaymentResponse struct {
	AuthorizationURL string          `json:"authorization_url"`
	Reference        string          `json:"reference"`
	Amount           float64         `json:"amount"`
	Payment          PaymentResponse `json:"payment"`
}
