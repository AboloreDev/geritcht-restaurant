package server

import (
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

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
	}

	utils.CreatedResponse(ctx, "Ingredient created successfully", response)
}

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
	}

	utils.SuccessResponse(ctx, "Ingredient updated successfully", response)
}

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
	}

	utils.SuccessResponse(ctx, "Ingredient deleted successfully", nil)
}

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
	}

	utils.SuccessResponse(ctx, "Ingredient fetched successfully", response)
}

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
	}

	utils.PaginatedSuccessResponse(ctx, "Ingredient fetched successfully", response, *meta)
}

func (s *Server) GetLowStockIngredientsHandler(ctx *gin.Context) {

	response, err := s.ingredientService.GetLowStockIngredientsService(ctx.Request.Context())
	if err != nil {
		switch err {
		case domain.ErrIngredientNotFound:
			utils.NotFound(ctx, "Ingredient not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to fetch ingredient", err)
		}
	}

	utils.SuccessResponse(ctx, "Ingredient fetched successfully", response)
}

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
	}

	utils.SuccessResponse(ctx, "Ingredient created successfully", nil)
}

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
	}

	utils.SuccessResponse(ctx, "Ingredient fetched successfully", nil)
}
