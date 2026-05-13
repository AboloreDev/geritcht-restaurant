package server

import (
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateDietaryTagHandler(ctx *gin.Context) {
	var req dto.CreateDietaryTagRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.dietaryTagsService.CreateDietaryTagService(&req)
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

	response, err := s.dietaryTagsService.UpdateDietaryTagService(tagID, &req)
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

func (s *Server) DeleteDietaryTagHandler(ctx *gin.Context) {

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid dietary tag ID", err)
		return
	}
	tagID := uint(id)

	err = s.dietaryTagsService.DeleteDietaryTagService(tagID)
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

func (s *Server) GetAllDietaryTagHandler(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	response, meta, err := s.dietaryTagsService.GetAllDietaryTagService(page, pageSize)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to fetch tags", err)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Dietary tags fetched successfully", response, *meta)
}
