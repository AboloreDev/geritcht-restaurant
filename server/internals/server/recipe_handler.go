package server

import (
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Add ingredient to a recipe
// @Description Link an ingredient to a menu item's recipe with the specified quantity. Admin access required.
// @Tags Recipes
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Param input body dto.LinkIngredientRequest true "Recipe ingredient details"
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=dto.RecipeResponse} "Ingredient added to recipe successfully"
// @Failure 400 {object} utils.Response "Invalid request data, ingredient ID, or quantity"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Ingredient not found"
// @Failure 409 {object} utils.Response "Ingredient already linked to this recipe"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /recipes/{id} [post]
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

	response, err := s.recipesService.AddMenuRecipe(ctx.Request.Context(), menuID, &req)
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

// @Summary Update recipe ingredient
// @Description Update the quantity or details of an ingredient linked to a menu item's recipe. Admin access required.
// @Tags Recipes
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Param ingredientID path int true "Ingredient ID"
// @Param input body dto.UpdateLinkItemRequest true "Updated recipe ingredient details"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.RecipeResponse} "Recipe updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Recipe or ingredient not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /recipes/{id}/{ingredientID} [patch]
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

	response, err := s.recipesService.UpdateMenuRecipe(ctx.Request.Context(), menuID, ingredientID, &req)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to update recipe", err)
		return
	}

	utils.SuccessResponse(ctx, "Recipe updated successfully", response)
}

// @Summary Remove ingredient from recipe
// @Description Remove an ingredient from a menu item's recipe. Admin access required.
// @Tags Recipes
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Param ingredientID path int true "Ingredient ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response "Recipe deleted successfully"
// @Failure 400 {object} utils.Response "Invalid ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Recipe or ingredient not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /recipes/{id}/{ingredientID} [delete]
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

	err = s.recipesService.DeleteRecipe(ctx.Request.Context(), menuID, ingredientID)
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

// @Summary Get recipe
// @Description Retrieve all ingredients linked to a menu item's recipe.
// @Tags Recipes
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]dto.RecipeResponse} "Recipe retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid menu ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Menu not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /recipes/{id} [get]
func (s *Server) GetRecipesHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid id", err)
		return
	}
	menuID := uint(id)

	response, err := s.recipesService.GetAllRecipes(ctx.Request.Context(), menuID)
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
