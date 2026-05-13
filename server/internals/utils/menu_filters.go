package utils

import (
	"strings"

	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"gorm.io/gorm"
)

func ApplyMenuFilters(query *gorm.DB, filter dto.MenuFilterRequest) *gorm.DB {
	if filter.CategoryID > 0 {
		query = query.Where("menu_category_id = ?", filter.CategoryID)
	}

	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}

	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}

	if filter.SpiceLevel > 0 {
		query = query.Where("spice_level = ?", filter.SpiceLevel)
	}

	if filter.Search != "" {
		query = query.Where(
			"name ILIKE ? OR description ILIKE ?",
			"%"+filter.Search+"%",
			"%"+filter.Search+"%",
		)
	}

	if filter.Dietary != "" {
		query = query.Joins(
			"JOIN menu_item_dietary ON menu_item_dietary.menu_id = menus.id",
		).Joins(
			"JOIN dietary_tags ON dietary_tags.id = menu_item_dietary.dietary_tag_id",
		).Where("dietary_tags.name ILIKE ?", filter.Dietary)
	}

	if filter.AllergenExclude != "" {
		allergens := strings.Split(filter.AllergenExclude, ",")
		query = query.Where(
			"menus.id NOT IN (SELECT menu_id FROM menu_item_allergens JOIN allergens ON allergens.id = menu_item_allergens.allergen_id WHERE allergens.name IN ?)",
			allergens,
		)
	}

	return query
}
