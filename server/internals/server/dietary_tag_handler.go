package server

import (
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Create a dietary tag
// @Description Create a new dietary tag. Admin access required.
// @Tags Dietary Tags
// @Accept json
// @Produce json
// @Param input body dto.CreateDietaryTagRequest true "Dietary tag details"
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=dto.DietaryTagResponse} "Dietary tag created successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 409 {object} utils.Response "Dietary tag already exists"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /tags [post]
func (s *Server) CreateDietaryTagHandler(ctx *gin.Context) {
	var req dto.CreateDietaryTagRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.dietaryTagsService.CreateDietaryTagService(ctx.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrNameConflict:
			utils.ConflictResponse(ctx, "Dietary tag already exists", err)
		default:
			utils.InternalServerError(ctx, "Failed to create dietary tag", err)
		}
		return
	}

	utils.CreatedResponse(ctx, "Dietary Tag created successfully", response)
}

// @Summary Update a dietary tag
// @Description Update an existing dietary tag by ID. Admin access required.
// @Tags Dietary Tags
// @Accept json
// @Produce json
// @Param id path string true "Dietary Tag ID"
// @Param input body dto.UpdateDietaryTagRequest true "Updated dietary tag details"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.DietaryTagResponse} "Dietary tag updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data or ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Dietary tag not found"
// @Failure 409 {object} utils.Response "Dietary tag name already exists"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /tags/{id} [patch]
func (s *Server) UpdateDietaryTagHandler(ctx *gin.Context) {
	var req dto.UpdateDietaryTagRequest

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid deitary tag ID", err)
		return
	}
	tagID := uint(id)

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.dietaryTagsService.UpdateDietaryTagService(ctx.Request.Context(), tagID, &req)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "Dietary Tag not found", err)
		case domain.ErrNameConflict:
			utils.ConflictResponse(ctx, "Dietary tag already exists", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Dietary tag updated successfully", response)
}

// @Summary Delete a dietary tag
// @Description Delete a dietary tag by ID. Admin access required.
// @Tags Dietary Tags
// @Accept json
// @Produce json
// @Param id path string true "Dietary Tag ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response "Dietary tag deleted successfully"
// @Failure 400 {object} utils.Response "Invalid dietary tag ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Dietary tag not found"
// @Failure 409 {object} utils.Response "Cannot delete dietary tag with associated meals"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /tags/{id} [delete]
func (s *Server) DeleteDietaryTagHandler(ctx *gin.Context) {

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid dietary tag ID", err)
		return
	}
	tagID := uint(id)

	err = s.dietaryTagsService.DeleteDietaryTagService(ctx.Request.Context(), tagID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "tag not found", err)
		default:
			utils.ConflictResponse(ctx, err.Error(), err)
		}
		return
	}

	utils.CreatedResponse(ctx, "Tag deleted successfully", nil)
}

// @Summary Get all dietary tags
// @Description Retrieve a list of all dietary tags with pagination.
// @Tags Dietary Tags
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Number of items per page"
// @Success 200 {object} utils.Response{data=[]dto.DietaryTagResponse, meta=utils.PaginatedMeta} "Dietary tags retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /tags [get]
func (s *Server) GetAllDietaryTagHandler(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	response, meta, err := s.dietaryTagsService.GetAllDietaryTagService(ctx.Request.Context(), page, pageSize)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to fetch tags", err)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Dietary tags fetched successfully", response, *meta)
}
