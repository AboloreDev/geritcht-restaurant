package services

import (
	"context"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
)

type MenuItemIngredientService struct {
	recipesRepo    repositories.RecipesRepositoryInterface
	ingredeintRepo repositories.IngredientRepositoryInterface
}

func NewMenuItemIngredientService(
	recipesRepo repositories.RecipesRepositoryInterface,
	ingredeintRepo repositories.IngredientRepositoryInterface) *MenuItemIngredientService {
	return &MenuItemIngredientService{
		recipesRepo:    recipesRepo,
		ingredeintRepo: ingredeintRepo,
	}
}

func (s *MenuItemIngredientService) AddMenuRecipe(ctx context.Context, menuItemID uint, req *dto.LinkIngredientRequest) (*dto.MenuItemIngredientResponse, error) {

	_, err := s.recipesRepo.CheckForLinkedIngredient(ctx, menuItemID, req.IngredientID)
	if err != nil {
		return nil, domain.ErrIngredientAlreadyLinked
	}

	_, err = s.ingredeintRepo.GetIngredientByID(ctx, req.IngredientID)
	if err != nil {
		return nil, domain.ErrIngredientNotFound
	}

	if req.Quantity <= 0 {
		return nil, domain.ErrInvalidQuantity
	}

	if req.IngredientID == 0 {
		return nil, domain.ErrInvalidIngredientID
	}

	recipe := models.MenuItemIngredient{
		MenuID:       menuItemID,
		IngredientID: req.IngredientID,
		Quantity:     req.Quantity,
		CreatedAt:    time.Now(),
	}

	err = s.recipesRepo.CreateRecipe(ctx, &recipe)
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

func (s *MenuItemIngredientService) UpdateMenuRecipe(ctx context.Context, menuItemID uint, ingredientID uint, req *dto.UpdateLinkItemRequest) (*dto.MenuItemIngredientResponse, error) {
	recipe, err := s.recipesRepo.GetLinkedIngredient(ctx, menuItemID, ingredientID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if req.Quantity != 0 {
		recipe.Quantity = req.Quantity
	}

	err = s.recipesRepo.UpdateLinkedIngredients(ctx, recipe)
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

func (s *MenuItemIngredientService) DeleteRecipe(ctx context.Context, menuItemID uint, ingredientID uint) error {
	err := s.recipesRepo.DeleteLinkedIngredient(ctx, menuItemID, ingredientID)
	if err != nil {
		return domain.ErrNotFound
	}

	return nil
}

func (s *MenuItemIngredientService) GetAllRecipes(ctx context.Context, menuItemID uint) ([]*dto.MenuItemIngredientResponse, error) {
	recipes, err := s.recipesRepo.GetRecipesByMenuItemID(ctx, menuItemID)
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

// func (s *MenuItemIngredientService) GetRecipe(menuItemID uint, ingredientID uint) (*dto.MenuItemIngredientResponse, error) {
// 	var recipe models.MenuItemIngredient

// 	err := s.db.Preload("Ingredient").Where("menu_item_id = ? AND ingredient_id = ?", menuItemID, ingredientID).First(&recipe).Error
// 	if err != nil {
// 		return nil, domain.ErrNotFound
// 	}

// 	return &dto.MenuItemIngredientResponse{
// 		Ingredient: dto.IngredientResponse{
// 			ID:           recipe.Ingredient.ID,
// 			Name:         recipe.Ingredient.Name,
// 			Unit:         recipe.Ingredient.Unit,
// 			CurrentStock: recipe.Ingredient.CurrentStock,
// 		},
// 		IngredientID: recipe.IngredientID,
// 		Quantity:     recipe.Quantity,
// 	}, nil
// }
