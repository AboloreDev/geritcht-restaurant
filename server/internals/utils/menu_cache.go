package utils

import (
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
)

func BuildMenuCacheKey(filter dto.MenuFilterRequest) string {
	return fmt.Sprintf(
		"menu:all:p%d:s%d:cat%d:min%.0f:max%.0f:spice%d:dietary%s:search%s:sort%s:%s",
		filter.Page,
		filter.PageSize,
		filter.CategoryID,
		filter.MinPrice,
		filter.MaxPrice,
		filter.SpiceLevel,
		filter.Dietary,
		filter.Search,
		filter.SortBy,
		filter.SortOrder,
	)
}

func GetCacheTTL(filter dto.MenuFilterRequest) time.Duration {
	hasFilter := filter.CategoryID > 0 ||
		filter.MinPrice > 0 ||
		filter.MaxPrice > 0 ||
		filter.Search != "" ||
		filter.Dietary != ""

	if hasFilter {
		return 30 * time.Minute
	}
	return 1 * time.Hour
}
