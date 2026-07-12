package server

import (
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Create an allergen
// @Description Create a new allergen. Admin access required.
// @Tags Allergens
// @Accept json
// @Produce json
// @Param input body dto.CreateAllergenRequest true "Allergen details"
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=dto.AllergenResponse} "Allergen created successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 409 {object} utils.Response "Allergen already exists"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /allergens [post]
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

// @Summary Update an allergen
// @Description Update an existing allergen. Admin access required.
// @Tags Allergens
// @Accept json
// @Produce json
// @Param id path int true "Allergen ID"
// @Param input body dto.UpdateAllergenRequest true "Updated allergen details"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.AllergenResponse} "Allergen updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data or allergen ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Allergen not found"
// @Failure 409 {object} utils.Response "Allergen name already exists"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /allergens/{id} [patch]
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

// @Summary Delete an allergen
// @Description Delete an allergen by its ID. Admin access required.
// @Tags Allergens
// @Accept json
// @Produce json
// @Param id path int true "Allergen ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response "Allergen deleted successfully"
// @Failure 400 {object} utils.Response "Invalid allergen ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Allergen not found"
// @Failure 409 {object} utils.Response "Allergen cannot be deleted because it is in use"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /allergens/{id} [delete]
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

	utils.SuccessResponse(ctx, "Allergen deleted successfully", nil)
}

// @Summary List allergens
// @Description Retrieve a paginated list of all allergens.
// @Tags Allergens
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Security BearerAuth
// @Success 200 {object} utils.PaginatedResponse{data=[]dto.AllergenResponse} "Allergens retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /allergens [get]
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
