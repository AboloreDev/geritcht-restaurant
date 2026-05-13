package utils

import (
	"fmt"
	"strings"

	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"gorm.io/gorm"
)

var allowedSortColumns = map[string]string{
	"price":       "price",
	"name":        "name",
	"created_at":  "created_at",
	"spice_level": "spice_level",
}

func ApplyMenuSorting(query *gorm.DB, filter dto.MenuFilterRequest) *gorm.DB {

	if filter.SortBy == "" {
		return query.Order("created_at DESC")
	}

	column, ok := allowedSortColumns[filter.SortBy]
	if !ok {
		return query.Order("created_at DESC")
	}

	order := "ASC"
	if strings.ToUpper(filter.SortOrder) == "DESC" {
		order = "DESC"
	}

	return query.Order(fmt.Sprintf("%s %s", column, order))
}
