package utils

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"gorm.io/gorm"
)

func ApplyOrderFilters(query *gorm.DB, filter *dto.OrderFilterRequest) *gorm.DB {
	if filter.Date > "" {
		query = query.Where("date = ?", filter.Date)
	}

	if filter.Status > "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.Type > "" {
		query = query.Where("type = ?", filter.Type)
	}

	return query
}
