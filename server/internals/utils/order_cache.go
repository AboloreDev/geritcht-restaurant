package utils

import (
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
)

func BuildOrderCacheKey(filter *dto.OrderFilterRequest) string {
	return fmt.Sprintf(
		"orders:all:page:%d:size:%d:date:%s:status:%s:type:%s",
		filter.Page,
		filter.PageSize,
		normalize(filter.Date),
		normalize(filter.Status),
		normalize(filter.Type),
	)
}

func GetOrderCacheTTL(filter *dto.OrderFilterRequest) time.Duration {
	hasFilter := filter.Date != "" || filter.Status != "" || filter.Type != ""

	if hasFilter {
		return 60 * time.Second
	}
	return 5 * time.Minute
}

func BuildUserOrderCacheKey(userID uint, filter *dto.OrderFilterRequest) string {
	return fmt.Sprintf(
		"orders:user:%d:page:%d:size:%d:date:%s:status:%s:type:%s",
		userID,
		filter.Page,
		filter.PageSize,
		normalize(filter.Date),
		normalize(filter.Status),
		normalize(filter.Type),
	)
}
