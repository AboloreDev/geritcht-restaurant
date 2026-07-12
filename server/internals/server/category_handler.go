package server

import (
	"errors"
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Create a category
// @Description Create a new product category with a name, description, and image: admin only
// @Tags Categories
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Category name"
// @Param description formData string false "Category description"
// @Param image formData file true "Category image"
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=dto.MenuCategoryResponse} "Category created successfully"
// @Failure 400 {object} utils.Response "Invalid request data or image upload failed"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 409 {object} utils.Response "Category already exists"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /categories [post]
func (s *Server) CreateCategoryHandler(ctx *gin.Context) {
	c := ctx.Request.Context()

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

	response, err := s.categoryServices.CreateCategoryService(c, &req, imageUrl)
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

// @Summary Update a category
// @Description Update an existing product category's name, description, and/or image admin only
// @Tags Categories
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Category ID"
// @Param name formData string false "Category name"
// @Param description formData string false "Category description"
// @Param image formData file false "Category image"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.MenuCategoryResponse} "Category updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data or image upload failed"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Category not found"
// @Failure 409 {object} utils.Response "Category name already exists"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /categories/{id} [patch]
func (s *Server) UpdateCategoryHandler(ctx *gin.Context) {
	c := ctx.Request.Context()

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

	response, err := s.categoryServices.UpdateCategoryService(c, categoryID, &req)
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

// @Summary Get a category by ID
// @Description Retrieve details of a specific product category by its ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} utils.Response{data=dto.MenuCategoryResponse} "Category fetched successfully"
// @Failure 400 {object} utils.Response "Invalid category ID"
// @Failure 404 {object} utils.Response "Category not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /categories/{id} [get]
func (s *Server) GetCategory(ctx *gin.Context) {
	c := ctx.Request.Context()

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	categoryID := uint(id)
	if err != nil {
		utils.BadRequest(ctx, "Invalid category ID", err)
		return
	}

	response, err := s.categoryServices.GetCategoryService(c, categoryID)
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

// @Summary Delete a category
// @Description Remove a product category by its ID. Admin only
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response "Category deleted successfully"
// @Failure 400 {object} utils.Response "Invalid category ID"
// @Failure 404 {object} utils.Response "Category not found"
// @Failure 409 {object} utils.Response "Category is associated with products"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /categories/{id} [delete]
func (s *Server) DeleteCategory(ctx *gin.Context) {
	c := ctx.Request.Context()

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	categoryID := uint(id)
	if err != nil {
		utils.BadRequest(ctx, "Invalid category ID", err)
		return
	}

	err = s.categoryServices.DeleteCategoryService(c, categoryID)
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

// @Summary Get all categories
// @Description Retrieve a paginated list of all product categories
// @Tags Categories
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Number of items per page"
// @Success 200 {object} utils.Response{data=[]dto.MenuCategoryResponse, meta=utils.PaginatedMeta} "Categories fetched successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /categories [get]
func (s *Server) GetCategoriesHandler(ctx *gin.Context) {
	c := ctx.Request.Context()

	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	response, meta, err := s.categoryServices.GetCategoriesService(c, page, pageSize)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to fetch categories", err)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Categories fetched successfully", response, *meta)
}

// @Summary Search categories
// @Description Search for categories by name or description
// @Tags Categories
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param page query int false "Page number"
// @Param pageSize query int false "Number of items per page"
// @Success 200 {object} utils.Response{data=[]dto.MenuCategoryResponse, meta=utils.PaginatedMeta} "Categories retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid search parameters"
// @Failure 404 {object} utils.Response "No categories found matching the search criteria"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /categories/search [get]
func (s *Server) SearchCategoryHandler(ctx *gin.Context) {
	var req dto.CategorySearchRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(ctx, "Invalid search parameters", err)
		return
	}

	response, meta, err := s.categoryServices.SearchCategory(ctx.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrCategoriesSearchNotFound:
			utils.NotFound(ctx, "Search returned no result", err)
		case domain.ErrInternalServerError:
			s.log.Error().Err(err).Msg("category search failed")
			utils.InternalServerError(ctx, "Something went wrong", errors.New("unable to complete search at this time"))
			return
		}

	}

	utils.PaginatedSuccessResponse(ctx, "category retrieved successfully", response, *meta)
}
