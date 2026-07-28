package utils

import (
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
)

func BuildReservationCacheKey(filter *dto.ReservationFilterRequest) string {
	return fmt.Sprintf(
		"reservations:all:page:%d:size:%d:date:%s:status:%s",
		filter.Page,
		filter.PageSize,
		normalize(filter.Date),
		normalize(filter.Status),
	)
}

func GetReservationCacheTTL(filter *dto.ReservationFilterRequest) time.Duration {
	hasFilter := filter.Date != "" || filter.Status != ""

	if hasFilter {
		return 50 * time.Second
	}
	return 20 * time.Minute
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

func normalize(value string) string {
	if value == "" {
		return "all"
	}
	return value
}
