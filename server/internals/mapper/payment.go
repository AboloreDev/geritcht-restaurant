package mapper

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
)

func PaymentResponse(payment *models.Payment) *dto.PaymentResponse {

	return &dto.PaymentResponse{
		ID:                payment.ID,
		OrderID:           payment.OrderID,
		Status:            string(payment.Status),
		Reference:         payment.Reference,
		Currency:          payment.Currency,
		Amount:            payment.Amount,
		CreatedAt:         payment.CreatedAt,
		PaidAt:            payment.PaidAt,
		ProviderReference: payment.ProviderReference,
		FailureReason:     payment.FailureReason,
		Provider:          payment.Provider,
	}
}
