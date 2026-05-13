package server

import (
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateMenuHandler(ctx *gin.Context) {
	var req dto.CreateMenuRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	response, err := s.menuServices.CreateMenuService(&req)
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

	response, err := s.menuServices.UpdateMenuService(menuID, &req)
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

func (s *Server) DeleteMenuHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid menu ID", err)
		return
	}
	menuID := uint(id)

	err = s.menuServices.DeleteMenu(menuID)
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

func (s *Server) GetMenuHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid menu ID", err)
		return
	}
	menuID := uint(id)

	response, err := s.menuServices.GetMenu(menuID)
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

func (s *Server) GetAllMenuHandler(ctx *gin.Context) {
	var filter dto.MenuFilterRequest

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		utils.BadRequest(ctx, "Invalid filter params", err)
		return
	}

	response, meta, err := s.menuServices.GetAllMenuService(filter)
	if err != nil {
		utils.InternalServerError(ctx, "Something went wrong", err)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Menu retrieved successfully", response, *meta)
}

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

	err = s.menuServices.AddMenuImageService(menuID, url, file.Filename)
	if err != nil {
		utils.BadRequest(ctx, "Error Adding", err)
		return
	}

	utils.SuccessResponse(ctx, "Upload successful", map[string]string{"url": url})
}

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

	if err := s.menuServices.RemoveMenuImageService(menuID); err != nil {
		utils.InternalServerError(ctx, "Failed to remove image record", err)
		return
	}

	utils.SuccessResponse(ctx, "Delete successful", nil)
}

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

	err = s.menuServices.ToggleMenuAvailabilityService(menuID, req.IsAvailable)
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
