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

func BuildMenuFetchCacheKey(filter dto.MenuFilterRequest) string {
    return fmt.Sprintf(
        "menu:page:%d:size:%d:cat:%d:search:%s:minp:%.2f:maxp:%.2f:spice:%d:diet:%s:allergen:%s",
        filter.Page,
        filter.PageSize,
        filter.CategoryID,
        filter.Search,
        filter.MinPrice,
        filter.MaxPrice,
        filter.SpiceLevel,
        filter.Dietary,
        filter.AllergenExclude,
    )
}

func GetMenuCacheTTL(filter *dto.MenuFilterRequest) time.Duration {
	hasFilter := 
	filter.CategoryID != 0 ||
		filter.MinPrice != 0 ||
		filter.MaxPrice != 0 ||
		filter.SpiceLevel != 0 ||
		filter.Dietary != "" ||
		filter.AllergenExclude != ""

	if hasFilter {
		return 30 * time.Minute
	}
	return 1 * time.Hour
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
