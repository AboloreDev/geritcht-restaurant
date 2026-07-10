package utils

import (
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
)

func BuildMenuCacheKey(filter *dto.MenuSearchRequest) string {
	return fmt.Sprintf(
		"menu:all:p%d:s%d:c%d:s%d:m%d:M%d:t%d:q%s",
		filter.Page,
		filter.Limit,
		filter.CategoryID,
		filter.SpiceLevel,
		filter.MaxPrice,
		filter.MinPrice,
		filter.PrepTimeMinutes,
		filter.Query,
	)
}

func GetCacheTTL(filter *dto.MenuSearchRequest) time.Duration {
	hasFilter := filter.CategoryID != nil ||
		filter.MinPrice != nil ||
		filter.MaxPrice != nil ||
		filter.SpiceLevel != nil ||
		filter.PrepTimeMinutes != nil

	if hasFilter {
		return 30 * time.Minute
	}
	return 1 * time.Hour
}
