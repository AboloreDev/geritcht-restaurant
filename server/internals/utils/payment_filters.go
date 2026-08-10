package utils

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"gorm.io/gorm"
)

func ApplyPaymentFilters(query *gorm.DB, filter *dto.PaymentFilterRequest) *gorm.DB {
	if filter.Amount > 0 {
		query = query.Where("amount = ?", filter.Amount)
	}

	if filter.Status > "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.Reference > "" {
		query = query.Where("reference = ?", filter.Reference)
	}

	return query
}
