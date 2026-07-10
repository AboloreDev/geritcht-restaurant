package repositories

import (
	"context"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"gorm.io/gorm"
)

type Recipe struct {
	db *gorm.DB
}

func NewRecipe(db *gorm.DB) *Recipe {
	return &Recipe{db: db}
}

func (r *Recipe) CheckForLinkedIngredient(ctx context.Context, menuItemID uint, ingredientID uint) (*models.MenuItemIngredient, error) {
	var recipe models.MenuItemIngredient

	err := r.db.Preload("Ingredient").
		WithContext(ctx).
		Where("menu_item_id = ? AND ingredient_id = ?", menuItemID, ingredientID).First(&recipe).Error
	if err == nil {
		return nil, domain.ErrIngredientAlreadyLinked
	}

	return &recipe, nil
}

func (r *Recipe) CreateRecipe(ctx context.Context, recipe *models.MenuItemIngredient) error {
	return r.db.WithContext(ctx).Create(recipe).Error
}

func (r *Recipe) GetLinkedIngredient(ctx context.Context, menuItemID uint, ingredientID uint) (*models.MenuItemIngredient, error) {
	var recipe models.MenuItemIngredient

	err := r.db.Preload("Ingredient").
		WithContext(ctx).
		Where("menu_item_id = ? AND ingredient_id = ?", menuItemID, ingredientID).First(&recipe).Error
	if err != nil {
		return nil, domain.ErrLinkedIngredeintNotFound
	}

	return &recipe, nil
}

func (r *Recipe) UpdateLinkedIngredients(ctx context.Context, recipe *models.MenuItemIngredient) error {
	return r.db.WithContext(ctx).Save(recipe).Error
}

func (r *Recipe) DeleteLinkedIngredient(ctx context.Context, menuItemID uint, ingredientID uint) error {
	result := r.db.WithContext(ctx).
		Where("menu_item_id = ? AND ingredient_id = ?", menuItemID, ingredientID).
		Delete(&models.MenuItemIngredient{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *Recipe) GetRecipesByMenuItemID(ctx context.Context, menuItemID uint) ([]models.MenuItemIngredient, error) {
	var recipes []models.MenuItemIngredient

	err := r.db.Preload("Ingredient").
		WithContext(ctx).
		Where("menu_item_id = ?", menuItemID).
		Find(&recipes).Error
	if err != nil {
		return nil, domain.ErrMenuNotFound
	}

	return recipes, nil
}
