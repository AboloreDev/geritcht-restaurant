package server

import (
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateCategoryHandler(ctx *gin.Context) {
	var req dto.CreateCategoryRequest
	image, err := ctx.FormFile("image")
	if err != nil {
		utils.BadRequest(ctx, "No file found", err)
		return
	}

	err = ctx.ShouldBind(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	imageUrl, err := s.uploadServices.UploadCategoryImage(image)
	if err != nil {
		utils.BadRequest(ctx, "Error uploading image", err)
		return
	}

	response, err := s.categoryServices.CreateCategoryService(&req, imageUrl)
	if err != nil {
		switch err {
		case domain.ErrNameConflict:
			utils.ConflictResponse(ctx, "Category with the same name already exists", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
		}
		return
	}

	utils.CreatedResponse(ctx, "Category created successfully", response)
}

func (s *Server) UpdateCategoryHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	categoryID := uint(id)
	if err != nil {
		utils.BadRequest(ctx, "Invalid category ID", err)
		return
	}

	var req dto.UpdateCategoryRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	response, err := s.categoryServices.UpdateCategoryService(categoryID, &req)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "Category not found", err)
		case domain.ErrNameConflict:
			utils.ConflictResponse(ctx, "Category name already exists", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Category updated successfully", response)
}

func (s *Server) GetCategory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	categoryID := uint(id)
	if err != nil {
		utils.BadRequest(ctx, "Invalid category ID", err)
		return
	}

	response, err := s.categoryServices.GetCategoryService(categoryID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.BadRequest(ctx, "Category not found", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Category fetched successfully", response)
}

func (s *Server) DeleteCategory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	categoryID := uint(id)
	if err != nil {
		utils.BadRequest(ctx, "Invalid category ID", err)
		return
	}

	err = s.categoryServices.DeleteCategoryService(categoryID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "Category not found", err)
		default:
			utils.ConflictResponse(ctx, err.Error(), err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Category deleted successfully", nil)
}

func (s *Server) GetCategoriesHandler(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	response, meta, err := s.categoryServices.GetCategoriesService(page, pageSize)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to fetch categories", err)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Categories fetched successfully", response, *meta)
}
