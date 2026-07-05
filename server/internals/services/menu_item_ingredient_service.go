package services

import (
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"gorm.io/gorm"
)

type MenuItemIngredientService struct {
	db *gorm.DB
}

func NewMenuItemIngredientService(db *gorm.DB) *MenuItemIngredientService {
	return &MenuItemIngredientService{
		db: db,
	}
}

func (s *MenuItemIngredientService) AddMenuRecipe(menuItemID uint, req *dto.LinkIngredientRequest) (*dto.MenuItemIngredientResponse, error) {
	var recipe models.MenuItemIngredient

	err := s.db.Preload("Ingredient").Where("menu_item_id = ? AND ingredient_id = ?", menuItemID, req.IngredientID).First(&recipe).Error
	if err == nil {
		return nil, domain.ErrIngredientAlreadyLinked
	}

	var ingredient models.Ingredient
	err = s.db.First(&ingredient, req.IngredientID).Error
	if err != nil {
		return nil, domain.ErrIngredientNotFound
	}

	if req.Quantity <= 0 {
		return nil, domain.ErrInvalidQuantity
	}

	if req.IngredientID == 0 {
		return nil, domain.ErrInvalidIngredientID
	}

	recipe = models.MenuItemIngredient{
		MenuID:       menuItemID,
		IngredientID: req.IngredientID,
		Quantity:     req.Quantity,
		CreatedAt:    time.Now(),
	}

	err = s.db.Create(&recipe).Error
	if err != nil {
		return nil, err
	}

	return &dto.MenuItemIngredientResponse{
		Ingredient: dto.IngredientResponse{
			ID:           recipe.Ingredient.ID,
			Name:         recipe.Ingredient.Name,
			Unit:         recipe.Ingredient.Unit,
			CurrentStock: recipe.Ingredient.CurrentStock,
		},
		IngredientID: recipe.IngredientID,
		Quantity:     recipe.Quantity,
	}, nil
}

func (s *MenuItemIngredientService) UpdateMenuRecipe(menuItemID uint, ingredientID uint, req *dto.UpdateLinkItemRequest) (*dto.MenuItemIngredientResponse, error) {
	var recipe models.MenuItemIngredient

	err := s.db.Preload("Ingredient").Where("menu_item_id = ? AND ingredient_id = ?", menuItemID, ingredientID).First(&recipe).Error
	if err != nil {
		return nil, err
	}

	if req.Quantity != 0 {
		recipe.Quantity = req.Quantity
	}

	err = s.db.Save(&recipe).Error
	if err != nil {
		return nil, err
	}

	return &dto.MenuItemIngredientResponse{
		Ingredient: dto.IngredientResponse{
			ID:           recipe.Ingredient.ID,
			Name:         recipe.Ingredient.Name,
			Unit:         recipe.Ingredient.Unit,
			CurrentStock: recipe.Ingredient.CurrentStock,
		},
		IngredientID: recipe.IngredientID,
		Quantity:     recipe.Quantity,
	}, nil
}

func (s *MenuItemIngredientService) DeleteRecipe(menuItemID uint, ingredientID uint) error {
	var recipe models.MenuItemIngredient

	err := s.db.Where("menu_item_id = ? AND ingredient_id = ?",
		menuItemID, ingredientID).First(&recipe).Error
	if err != nil {
		return domain.ErrNotFound
	}

	result := s.db.Delete(&recipe)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (s *MenuItemIngredientService) GetRecipes(menuItemID uint) ([]*dto.MenuItemIngredientResponse, error) {
	var recipes []models.MenuItemIngredient

	err := s.db.Preload("Ingredient").Where("menu_item_id = ?", menuItemID).Find(&recipes).Error
	if err != nil {
		return nil, domain.ErrMenuNotFound
	}

	response := make([]*dto.MenuItemIngredientResponse, len(recipes))

	for i, recipe := range recipes {
		response[i] = &dto.MenuItemIngredientResponse{
			Ingredient: dto.IngredientResponse{
				ID:           recipe.Ingredient.ID,
				Name:         recipe.Ingredient.Name,
				Unit:         recipe.Ingredient.Unit,
				CurrentStock: recipe.Ingredient.CurrentStock,
			},
			IngredientID: recipe.IngredientID,
			Quantity:     recipe.Quantity,
		}
	}

	return response, nil
}
