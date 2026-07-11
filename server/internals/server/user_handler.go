package server

import (
	"errors"
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetUserProfileHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	response, err := s.userServices.GetUserProfileService(ctx.Request.Context(), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.BadRequest(ctx, "User not found", err)
		default:
			utils.UnAuthorized(ctx, "Unauthorized", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "User fetched successfully", response)
}

func (s *Server) GetStaffProfileHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	response, err := s.userServices.GetStaffProfileService(ctx.Request.Context(), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.BadRequest(ctx, "Staff not found", err)
		default:
			utils.UnAuthorized(ctx, "Unauthorized", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Staff fetched successfully", response)
}

func (s *Server) GetAllUserHandler(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	response, meta, err := s.userServices.GetAllUsersService(ctx.Request.Context(), page, pageSize)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to fetch users", err)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Users fetched successfully", response, *meta)
}

func (s *Server) GetAllStaffsHandler(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	response, meta, err := s.userServices.GetAllStaffService(ctx.Request.Context(), page, pageSize)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to fetch users", err)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Users fetched successfully", response, *meta)
}

func (s *Server) DeactivateUserHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	userID := uint(id)

	err = s.userServices.DeactivateUserService(ctx.Request.Context(), userID)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to deactivate users", err)
		return
	}

	utils.SuccessResponse(ctx, "User deactivated successfully", nil)
}

func (s *Server) DeactivateStaffHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	userID := uint(id)

	err = s.userServices.DeactivateStaffService(ctx.Request.Context(), userID)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to deactivate staff", err)
		return
	}

	utils.SuccessResponse(ctx, "Staff deactivated successfully", nil)
}

func (s *Server) ActivateStaffHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	userID := uint(id)

	err = s.userServices.ActivateUserService(ctx.Request.Context(), userID)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to activate staff", err)
		return
	}

	utils.SuccessResponse(ctx, "staff activate successfully", nil)
}

func (s *Server) ActivateUserHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	userID := uint(id)

	err = s.userServices.ActivateStaffService(ctx.Request.Context(), userID)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to activate users", err)
		return
	}

	utils.SuccessResponse(ctx, "User activated successfully", nil)
}

func (s *Server) UpdateUserProfileHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var req dto.UpdateProfileRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.userServices.UpdateProfileService(ctx.Request.Context(), userID, &req)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "User not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to update profile", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "User updated successfully", response)
}

func (s *Server) UpdateStaffProfileHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var req dto.UpdateProfileRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.userServices.UpdateProfileService(ctx.Request.Context(), userID, &req)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "User not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to update profile", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Staff updated successfully", response)
}

func (s *Server) SearchUserHandler(ctx *gin.Context) {
	var req dto.UserSearchRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(ctx, "Invalid search parameters", err)
		return
	}

	response, meta, err := s.userServices.SearchUser(ctx.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrUserSearchNotFound:

			utils.NotFound(ctx, "Search returned no result", err)
		case domain.ErrInternalServerError:
			s.log.Error().Err(err).Msg("User search failed")
			utils.InternalServerError(ctx, "Something went wrong", errors.New("unable to complete search at this time"))
			return
		}

	}

	utils.PaginatedSuccessResponse(ctx, "user retrieved successfully", response, *meta)
}
