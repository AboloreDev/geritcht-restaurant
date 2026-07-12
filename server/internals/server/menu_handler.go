package server

import (
	"errors"
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Create a menu item
// @Description Create a new menu item. The referenced category, dietary tags, and ingredients must already exist. Admin access required.
// @Tags Menu
// @Accept json
// @Produce json
// @Param input body dto.CreateMenuRequest true "Menu item details"
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=dto.MenuResponse} "Menu item created successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Related resource not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /menu [post]
func (s *Server) CreateMenuHandler(ctx *gin.Context) {
	var req dto.CreateMenuRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	response, err := s.menuServices.CreateMenuService(ctx.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "Related resource not found", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
		}
		return
	}

	utils.CreatedResponse(ctx, "Menu created successfully", response)
}

// @Summary Update a menu item
// @Description Update an existing menu item by ID. The referenced category, dietary tags, and ingredients must already exist. Admin access required.
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "Menu Item ID"
// @Param input body dto.UpdateMenuRequest true "Updated menu item details"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.MenuResponse} "Menu item updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data or menu ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Menu item not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /menu/{id} [patch]
func (s *Server) UpdateMenuHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid menu ID", err)
		return
	}
	menuID := uint(id)

	var req dto.UpdateMenuRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	response, err := s.menuServices.UpdateMenuService(ctx.Request.Context(), menuID, &req)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "Related resource not found", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Menu updated successfully", response)
}

// @Summary Delete a menu item
// @Description Delete an existing menu item by ID. Admin access required.
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "Menu Item ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response "Menu item deleted successfully"
// @Failure 400 {object} utils.Response "Invalid menu ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Menu item not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /menu/{id} [delete]
func (s *Server) DeleteMenuHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid menu ID", err)
		return
	}
	menuID := uint(id)

	err = s.menuServices.DeleteMenu(ctx.Request.Context(), menuID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "Menu not found", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Menu deleted successfully", nil)
}

// @Summary Get a menu item
// @Description Retrieve an existing menu item by ID. Admin access required.
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "Menu Item ID"
// @Success 200 {object} utils.Response{data=dto.MenuResponse} "Menu item retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid menu ID"
// @Failure 404 {object} utils.Response "Menu item not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /menu/{id} [get]
func (s *Server) GetMenuHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid menu ID", err)
		return
	}
	menuID := uint(id)

	response, err := s.menuServices.GetMenu(ctx.Request.Context(), menuID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "Menu not found", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Menu retrieved successfully", response)
}

// @Summary Get all menu items with filtering and pagination
// @Description Retrieve menu items with optional filtering by category, dietary tags, and ingredients, along with pagination. Admin access required.
// @Tags Menu
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param category query string false "Category name"
// @Param dietary_tags query string false "Comma-separated list of dietary tags"
// @Param ingredients query string false "Comma-separated list of ingredient names"
// @Success 200 {object} utils.Response{data=[]dto.MenuResponse, meta=utils.PaginatedMeta} "Menu items retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid filter params"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /menu [get]
func (s *Server) GetAllMenuHandler(ctx *gin.Context) {
	var filter dto.MenuFilterRequest

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		utils.BadRequest(ctx, "Invalid filter params", err)
		return
	}

	response, meta, err := s.menuServices.GetAllMenuService(ctx.Request.Context(), filter)
	if err != nil {
		utils.InternalServerError(ctx, "Something went wrong", err)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Menu retrieved successfully", response, *meta)
}

// @Summary Upload a menu item image
// @Description Upload an image for an existing menu item. Admin access required.
// @Tags Menu
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Menu ID"
// @Param image formData file true "Menu image"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=map[string]string} "Image uploaded successfully"
// @Failure 400 {object} utils.Response "Invalid menu ID, missing image, or upload failed"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Menu not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /menu/{id}/image [post]
func (s *Server) UploadMenuImageHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Menu ID", err)
		return
	}
	menuID := uint(id)

	file, err := ctx.FormFile("image")
	if err != nil {
		utils.BadRequest(ctx, "No file uploaded", err)
		return
	}

	url, err := s.uploadServices.UploadMenuImage(menuID, file)
	if err != nil {
		utils.BadRequest(ctx, "Failed to upload file", err)
		return
	}

	err = s.menuServices.AddMenuImageService(ctx.Request.Context(), menuID, url, file.Filename)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "Menu not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to associate image with menu", err)
		}
		return
	}
	utils.SuccessResponse(ctx, "Upload successful", map[string]string{"url": url})
}

// @Summary Delete a menu item image
// @Description Delete the image associated with a menu item. Admin access required.
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response "Image deleted successfully"
// @Failure 400 {object} utils.Response "Invalid menu ID or delete failed"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Menu not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /menu/{id}/image [delete]
func (s *Server) DeleteMenuImageHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid menu ID", err)
		return
	}
	menuID := uint(id)

	err = s.uploadServices.DeleteFile(menuID)
	if err != nil {
		utils.BadRequest(ctx, "Failed to delete file", err)
		return
	}

	if err := s.menuServices.RemoveMenuImageService(ctx.Request.Context(), menuID); err != nil {
		utils.InternalServerError(ctx, "Failed to remove image record", err)
		return
	}

	utils.SuccessResponse(ctx, "Delete successful", nil)
}

// @Summary Toggle menu item availability
// @Description Toggle the availability status of a menu item. Admin access required.
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response "Menu availability toggled successfully"
// @Failure 400 {object} utils.Response "Invalid menu ID or missing is_available field"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Menu not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /menu/{id}/toggle [patch]
func (s *Server) ToggleMenuAvailabilityHandler(ctx *gin.Context) {
	var req struct {
		IsAvailable *bool `json:"is_available" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "is_available is required", err)
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid menu ID", err)
		return
	}
	menuID := uint(id)

	err = s.menuServices.ToggleMenuAvailabilityService(ctx.Request.Context(), menuID, req.IsAvailable)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "Menu not found", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Menu availability toggled successfully", nil)
}

// @Summary Search menu items
// @Description Search menu items by name, category, dietary tags, or ingredients with pagination.
// @Tags Menu
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} utils.Response{data=[]dto.MenuResponse, meta=utils.PaginatedMeta} "Menus retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid search parameters"
// @Failure 404 {object} utils.Response "No results found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /menu/search [get]
func (s *Server) SearchMenuHandler(ctx *gin.Context) {
	var req dto.MenuSearchRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(ctx, "Invalid search parameters", err)
		return
	}

	response, meta, err := s.menuServices.SearchProduct(ctx.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrMenuSearchNotFound:
			utils.NotFound(ctx, "Search returned no result", err)
		case domain.ErrInternalServerError:
			s.log.Error().Err(err).Msg("Menu search failed")
			utils.InternalServerError(ctx, "Something went wrong", errors.New("unable to complete search at this time"))
		default:
			s.log.Error().Err(err).Msg("Menu search failed")
			utils.InternalServerError(ctx, "Something went wrong", errors.New("unable to complete search at this time"))
		}
		return

	}

	utils.PaginatedSuccessResponse(ctx, "Menus retrieved successfully", response, *meta)
}
