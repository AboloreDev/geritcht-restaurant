package models

import (
	"time"

	"gorm.io/gorm"
)

type MenuCategory struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"not null"`
	Description  string         `json:"description"`
	ImageURL     string         `json:"image_url"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	DisplayOrder int            `json:"display_order" gorm:"default:0"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Menu []Menu `json:"-"`
}

type Menu struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	MenuCategoryID  uint           `json:"menu_category_id" gorm:"not null;index"`
	Name            string         `json:"name" gorm:"not null"`
	Description     string         `json:"description"`
	Price           float64        `json:"price" gorm:"not null"`
	ImageURL        string         `json:"image_url"`
	IsAvailable     bool           `json:"is_available" gorm:"default:true"`
	PrepTimeMinutes int            `json:"prep_time_minutes" gorm:"default:15"`
	SpiceLevel      int            `json:"spice_level" gorm:"default:0"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
	DisplayOrder    int            `json:"display_order" gorm:"default:0"`
	// Relationships
	MenuCategory        MenuCategory         `json:"-"`
	Allergens           []Allergen           `json:"allergens,omitempty" gorm:"many2many:menu_item_allergens"`
	DietaryTags         []DietaryTag         `json:"dietary_tags,omitempty" gorm:"many2many:menu_item_dietary"`
	MenuItemIngredients []MenuItemIngredient `json:"-"`
	Images              []MenuImage          `json:"images,omitempty"`
}

type MenuImage struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	MenuID    uint           `json:"menu_id" gorm:"not null;index"`
	URL       string         `json:"url" gorm:"not null"`
	AltText   string         `json:"alt_text"`
	IsPrimary bool           `json:"is_primary" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	// Relationships
	Menu Menu `json:"-"`
}

type Allergen struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null;uniqueIndex"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	// Relationships
	Menu []Menu `json:"-" gorm:"many2many:menu_item_allergens"`
}

type DietaryTag struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null;uniqueIndex"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	// Relationships
	Menu []Menu `json:"-" gorm:"many2many:menu_item_dietary"`
}

type MenuWithRank struct {
	Menu
	Rank float32 `gorm:"column:rank"`
}

type MenuCategoryWithRank struct {
	MenuCategory
	Rank float32 `gorm:"column:rank"`
}