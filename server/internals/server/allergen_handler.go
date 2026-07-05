package server

import (
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateAllergenHandler(ctx *gin.Context) {
	var req dto.CreateAllergenRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.allergenServices.CreateAllergenServices(ctx.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrNameConflict:
			utils.ConflictResponse(ctx, "Allergen already exists", err)
		default:
			utils.InternalServerError(ctx, "Failed to create allergen", err)
		}
		return
	}

	utils.CreatedResponse(ctx, "Allergen created successfully", response)
}

func (s *Server) UpdateAllegenHandler(ctx *gin.Context) {
	var req dto.UpdateAllergenRequest

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid allergen ID", err)
		return
	}
	allergenID := uint(id)

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.allergenServices.UpdateAllergenService(ctx.Request.Context(), allergenID, &req)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "Allergen not found", err)
		case domain.ErrNameConflict:
			utils.ConflictResponse(ctx, "Allegen name already exists", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Allergen updated successfully", response)
}

func (s *Server) DeleteAllergenHandler(ctx *gin.Context) {

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid allergen ID", err)
		return
	}
	allergenID := uint(id)

	err = s.allergenServices.DeleteAllergenService(ctx.Request.Context(), allergenID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "Allergen not found", err)
		default:
			utils.ConflictResponse(ctx, err.Error(), err)
		}
		return
	}

	utils.CreatedResponse(ctx, "Allergen deleted successfully", nil)
}
func (s *Server) GetAllAllergenHandler(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	response, meta, err := s.allergenServices.GetAllAllergenService(ctx.Request.Context(), page, pageSize)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to fetch allergens", err)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Allergens fetched successfully", response, *meta)
}
