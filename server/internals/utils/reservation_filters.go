package utils

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"gorm.io/gorm"
)

func ApplyReservationFilters(query *gorm.DB, filter *dto.ReservationFilterRequest) *gorm.DB {
	if filter.Date > "" {
		query = query.Where("date = ?", filter.Date)
	}

	if filter.Status > "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.TimeSlot > "" {
		query = query.Where("time_slot = ?", filter.TimeSlot)
	}

	return query
}
