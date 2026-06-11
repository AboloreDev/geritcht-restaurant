package dto

import "time"

type RefundResponse struct {
	ID          uint            `json:"id"`
	OrderID     uint            `json:"order_id"`
	PaymentID   uint            `json:"payment_id"`
	Reference   string          `json:"reference"`
	Amount      float64         `json:"amount"`
	Currency    string          `json:"currency"`
	Status      string          `json:"status"`
	Reason      string          `json:"reason"`
	ProcessedAt *time.Time      `json:"processed_at"`
	CreatedAt   time.Time       `json:"created_at"`
	Order       OrderResponse   `json:"order,omitempty"`
	Payment     PaymentResponse `json:"payment,omitempty"`
}
