package utils

import (
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
)

func BuildReservationCacheKey(filter *dto.ReservationFilterRequest) string {
	return fmt.Sprintf(
		"reservations:all:page:%d:size:%d:date:%s:status:%s:timeslot:%s",
		filter.Page,
		filter.PageSize,
		normalize(filter.Date),
		normalize(filter.Status),
		normalize(filter.TimeSlot),
	)
}

func GetReservationCacheTTL(filter *dto.ReservationFilterRequest) time.Duration {
	hasFilter := filter.Date != "" || filter.Status != "" || filter.TimeSlot != ""

	if hasFilter {
		return 50 * time.Second
	}
	return 5 * time.Minute
}

func BuildUserReservationCacheKey(userID uint, filter *dto.ReservationFilterRequest) string {
	return fmt.Sprintf(
		"reservations:user:%d:page:%d:size:%d:date:%s:status:%s",
		userID,
		filter.Page,
		filter.PageSize,
		normalize(filter.Date),
		normalize(filter.Status),
	)
}

func normalize(value interface{}) interface{} {
	if value == "" {
		return "all"
	}
	return value
}
