package server

import (
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) AddRecipeHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid id", err)
		return
	}
	menuID := uint(id)

	var req dto.LinkIngredientRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request body", err)
		return
	}

	response, err := s.recipesService.AddMenuRecipe(menuID, &req)
	if err != nil {
		switch err {
		case domain.ErrIngredientAlreadyLinked:
			utils.ConflictResponse(ctx, "Ingredient already linked to this recipe", err)
			return
		case domain.ErrIngredientNotFound:
			utils.NotFound(ctx, "Ingredient not found", err)
			return
		case domain.ErrInvalidQuantity:
			utils.BadRequest(ctx, "Invalid quantity", err)
			return
		case domain.ErrInvalidIngredientID:
			utils.BadRequest(ctx, "Invalid ingredient id", err)
			return
		default:
			utils.InternalServerError(ctx, "Failed to add recipe", err)
			return
		}
	}

	utils.CreatedResponse(ctx, "Recipe added successfully", response)
}

func (s *Server) UpdateRecipeHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid id", err)
		return
	}
	menuID := uint(id)

	ingredientIDStr := ctx.Param("ingredientID")
	ingredient_id, err := strconv.ParseUint(ingredientIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid id", err)
		return
	}
	ingredientID := uint(ingredient_id)

	var req dto.UpdateLinkItemRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request body", err)
		return
	}

	response, err := s.recipesService.UpdateMenuRecipe(menuID, ingredientID, &req)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to update recipe", err)
		return
	}

	utils.SuccessResponse(ctx, "Recipe updated successfully", response)
}

func (s *Server) DeleteRecipeHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid id", err)
		return
	}
	menuID := uint(id)

	ingredientIDStr := ctx.Param("ingredientID")
	ingredient_id, err := strconv.ParseUint(ingredientIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid id", err)
		return
	}
	ingredientID := uint(ingredient_id)

	err = s.recipesService.DeleteRecipe(menuID, ingredientID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "Ingredient not found", err)
			return
		default:
			utils.InternalServerError(ctx, "Failed to delete recipe", err)
			return
		}
	}

	utils.SuccessResponse(ctx, "Recipe deleted successfully", nil)
}

func (s *Server) GetRecipesHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid id", err)
		return
	}
	menuID := uint(id)

	response, err := s.recipesService.GetRecipes(menuID)
	if err != nil {
		switch err {
		case domain.ErrMenuNotFound:
			utils.NotFound(ctx, "Menu not found", err)
			return
		default:
			utils.InternalServerError(ctx, "Failed to fetch recipe", err)
			return
		}
	}

	utils.SuccessResponse(ctx, "Recipe fetched successfully", response)
}
