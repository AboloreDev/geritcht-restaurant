package server

import (
	"errors"
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Create an ingredient
// @Description Create a new ingredient. Admin access required.
// @Tags Ingredients
// @Accept json
// @Produce json
// @Param input body dto.CreateIngredientRequest true "Ingredient details"
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=dto.IngredientResponse} "Ingredient created successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 409 {object} utils.Response "Ingredient already exists"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /ingredients [post]
func (s *Server) CreateIngredientHandler(ctx *gin.Context) {
	var req dto.CreateIngredientRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.ingredientService.CreateIngredientService(ctx.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrConflict:
			utils.ConflictResponse(ctx, "Ingredient already exist", err)
		default:
			utils.InternalServerError(ctx, "Failed to create ingredient", err)
		}
		return
	}

	utils.CreatedResponse(ctx, "Ingredient created successfully", response)
}

// @Summary Update an ingredient
// @Description Update an existing ingredient by ID. Admin access required.
// @Tags Ingredients
// @Accept json
// @Produce json
// @Param id path int true "Ingredient ID"
// @Param input body dto.UpdateIngredientRequest true "Updated ingredient details"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.IngredientResponse} "Ingredient updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data or ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Ingredient not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /ingredients/{id} [patch]
func (s *Server) UpdateIngredientHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid ID", err)
		return
	}
	ingredientID := uint(id)

	var req dto.UpdateIngredientRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.ingredientService.UpdateIngredientService(ctx.Request.Context(), ingredientID, &req)
	if err != nil {
		switch err {
		case domain.ErrIngredientNotFound:
			utils.NotFound(ctx, "Ingredient not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to update ingredient", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Ingredient updated successfully", response)
}

// @Summary Delete an ingredient
// @Description Delete an ingredient by ID. Admin access required.
// @Tags Ingredients
// @Accept json
// @Produce json
// @Param id path int true "Ingredient ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response "Ingredient deleted successfully"
// @Failure 400 {object} utils.Response "Invalid ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Ingredient not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /ingredients/{id} [delete]
func (s *Server) DeleteIngredientHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid ID", err)
		return
	}
	ingredientID := uint(id)

	err = s.ingredientService.DeleteIngredientService(ctx.Request.Context(), ingredientID)
	if err != nil {
		switch err {
		case domain.ErrIngredientNotFound:
			utils.NotFound(ctx, "Ingredient not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to delete ingredient", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Ingredient deleted successfully", nil)
}

// @Summary Get an ingredient
// @Description Retrieve an ingredient by its ID Admin access required.
// @Tags Ingredients
// @Accept json
// @Produce json
// @Param id path int true "Ingredient ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.IngredientResponse} "Ingredient retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid ingredient ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Ingredient not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /ingredients/{id} [get]
func (s *Server) GetIngredientHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid ID", err)
		return
	}
	ingredientID := uint(id)

	response, err := s.ingredientService.GetIngredientService(ctx.Request.Context(), ingredientID)
	if err != nil {
		switch err {
		case domain.ErrIngredientNotFound:
			utils.NotFound(ctx, "Ingredient not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to fetch ingredient", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Ingredient fetched successfully", response)
}

// @Summary Get all ingredients
// @Description Retrieve a list of all ingredients with pagination.Admin access required.
// @Tags Ingredients
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Number of items per page"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]dto.IngredientResponse, meta=utils.PaginatedMeta} "Ingredients retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Ingredient not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /ingredients [get]
func (s *Server) GetAllIngredientHandler(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	response, meta, err := s.ingredientService.GetAllIngredientService(ctx.Request.Context(), page, pageSize)
	if err != nil {
		switch err {
		case domain.ErrIngredientNotFound:
			utils.NotFound(ctx, "Ingredient not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to fetch all ingredient", err)
		}
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Ingredient fetched successfully", response, *meta)
}

// @Summary Get low stock ingredients
// @Description Retrieve a list of ingredients with stock below their threshold Admin access required.
// @Tags Ingredients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]dto.IngredientResponse} "Low stock ingredients retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Ingredient not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /ingredients/low-stock [get]
func (s *Server) GetLowStockIngredientsHandler(ctx *gin.Context) {

	response, err := s.ingredientService.GetLowStockIngredientsService(ctx.Request.Context())
	if err != nil {
		switch err {
		case domain.ErrIngredientNotFound:
			utils.NotFound(ctx, "Ingredient not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to fetch ingredient", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Ingredient fetched successfully", response)
}

// @Summary Set threshold limit
// @Description Set the threshold limit for an ingredient. Admin access required.
// @Tags Ingredients
// @Accept json
// @Produce json
// @Param id path int true "Ingredient ID"
// @Param input body dto.ThresholdRequest true "Threshold limit"
// @Security BearerAuth
// @Success 200 {object} utils.Response "Threshold limit set successfully"
// @Failure 400 {object} utils.Response "Invalid request data or ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Ingredient not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /ingredients/{id}/limit [patch]
func (s *Server) SetThresholdLimitHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid ID", err)
		return
	}
	ingredientID := uint(id)

	var req dto.ThresholdRequest
	err = ctx.ShouldBind(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	err = s.ingredientService.SetThresholdLimit(ctx.Request.Context(), ingredientID, &req)
	if err != nil {
		switch err {
		case domain.ErrIngredientNotFound:
			utils.NotFound(ctx, "Ingredient not found", err)
		case domain.ErrNegativeThreshold:
			utils.BadRequest(ctx, "threshold can't be negative", err)
		default:
			utils.InternalServerError(ctx, "Failed to create ingredient", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Ingredient created successfully", nil)
}

// @Summary Check ingredient stock
// @Description Check whether an ingredient has fallen below its configured stock threshold and notify the authenticated user if necessary
// @Tags Ingredients
// @Accept json
// @Produce json
// @Param id path int true "Ingredient ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response "Low stock check completed successfully"
// @Failure 400 {object} utils.Response "Invalid ingredient ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Ingredient or user not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /ingredients/{id}/check-low-stock [post]
func (s *Server) CheckLowStockHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid ID", err)
		return
	}
	ingredientID := uint(id)

	err = s.ingredientService.CheckLowStock(ctx.Request.Context(), userID, ingredientID)
	if err != nil {
		switch err {
		case domain.ErrIngredientNotFound:
			utils.NotFound(ctx, "Ingredient not found", err)
		case domain.ErrUserNotFound:
			utils.NotFound(ctx, "User not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to fetch ingredient", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Ingredient fetched successfully", nil)
}

// @Summary Search ingredients
// @Description Search ingredients by name with pagination
// @Tags Ingredients
// @Accept json
// @Produce json
// @Param query query string false "Ingredient name"
// @Param page query int false "Page number"
// @Param pageSize query int false "Number of items per page"
// @Security BearerAuth
// @Success 200 {object} utils.PaginatedResponse{data=[]dto.IngredientResponse} "Ingredients retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid search parameters"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "No ingredients found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /ingredients/search [get]
func (s *Server) SearchIngredientHandler(ctx *gin.Context) {
	var req dto.IngredientSearchRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(ctx, "Invalid search parameters", err)
		return
	}

	response, meta, err := s.ingredientService.SearchIngredients(ctx.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrIngredientSearchNotFound:
			utils.NotFound(ctx, "Search returned no result", err)
		default:
			s.log.Error().Err(err).Msg("Ingredient search failed")
			utils.InternalServerError(ctx, "Something went wrong", errors.New("unable to complete search at this time"))

		}
		return
	}

	utils.PaginatedSuccessResponse(ctx, "ingredient retrieved successfully", response, *meta)
}
