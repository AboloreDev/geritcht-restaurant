package dto

import "time"

type CreateCategoryRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"display_order"`
}

type UpdateCategoryRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

type CreateMenuRequest struct {
	CategoryID      uint    `json:"category_id" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	Description     string  `json:"description"`
	Price           float64 `json:"price" binding:"required,gt=0"`
	PrepTimeMinutes int     `json:"prep_time_minutes"`
	SpiceLevel      int     `json:"spice_level"`
	AllergenIDs     []uint  `json:"allergen_ids"`
	DietaryTagIDs   []uint  `json:"dietary_tag_ids"`
}

type UpdateMenuRequest struct {
	CategoryID      uint    `json:"category_id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	PrepTimeMinutes int     `json:"prep_time_minutes"`
	SpiceLevel      int     `json:"spice_level"`
	IsAvailable     *bool   `json:"is_available"`
	AllergenIDs     []uint  `json:"allergen_ids"`
	DietaryTagIDs   []uint  `json:"dietary_tag_ids"`
}

type ToggleAvailabilityRequest struct {
	IsAvailable bool `json:"is_available" binding:"required"`
}

type CreateAllergenRequest struct {
	Name string `json:"name" binding:"required"`
}
type CreateDietaryTagRequest struct {
	Name string `json:"name" binding:"required"`
}

type MenuFilterRequest struct {
	CategoryID    uint    `form:"category_id"`
	MinPrice      float64 `form:"min_price"`
	MaxPrice      float64 `form:"max_price"`
	IsAvailable   *bool   `form:"is_available"`
	AllergenIDs   []uint  `form:"allergen_ids"`
	DietaryTagIDs []uint  `form:"dietary_tag_ids"`
	SpiceLevel    int     `form:"spice_level"`
	Search        string  `form:"search"`
}

type AllergenResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type DietaryTagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type MenuImageResponse struct {
	ID        uint      `json:"id"`
	URL       string    `json:"url"`
	AltText   string    `json:"alt_text"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
}

type MenuCategoryResponse struct {
	ID           uint           `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	ImageURL     string         `json:"image_url"`
	DisplayOrder int            `json:"display_order"`
	IsActive     bool           `json:"is_active"`
	MenuItems    []MenuResponse `json:"menu_items,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

type MenuResponse struct {
	ID              uint                  `json:"id"`
	CategoryID      uint                  `json:"category_id"`
	Category        *MenuCategoryResponse `json:"category,omitempty"`
	Name            string                `json:"name"`
	Description     string                `json:"description"`
	Price           float64               `json:"price"`
	ImageURL        string                `json:"image_url"`
	IsAvailable     bool                  `json:"is_available"`
	PrepTimeMinutes int                   `json:"prep_time_minutes"`
	SpiceLevel      int                   `json:"spice_level"`
	Allergens       []AllergenResponse    `json:"allergens,omitempty"`
	DietaryTags     []DietaryTagResponse  `json:"dietary_tags,omitempty"`
	Images          []MenuImageResponse   `json:"images,omitempty"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}
