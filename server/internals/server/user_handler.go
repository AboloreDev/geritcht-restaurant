package server

import (
	"errors"
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Get user profile
// @Description Retrieve the authenticated user's profile.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.UserResponse} "User profile retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /users/profile [get]
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

// @Summary Get staff profile
// @Description Retrieve the authenticated staff member's profile.
// @Tags Staff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.UserResponse} "Staff profile retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Staff not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /staff/profile [get]
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

// @Summary List users
// @Description Retrieve a paginated list of all registered users. Admin access required.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]dto.UserResponse} "Users retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /users [get]
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

// @Summary List staff
// @Description Retrieve a paginated list of all staff members. Admin access required.
// @Tags Staff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]dto.UserResponse} "Staff retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /staff [get]
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

// @Summary Deactivate user
// @Description Deactivate a user account. Admin access required.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response "User deactivated successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /users/profile/deactivate [patch]
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

// @Summary Deactivate staff
// @Description Deactivate a staff account. Admin access required.
// @Tags Staff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Staff ID"
// @Success 200 {object} utils.Response "Staff deactivated successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Staff not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /staff/profile/deactivate [patch]
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

// @Summary Activate staff
// @Description Activate a staff account. Admin access required.
// @Tags Staff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Staff ID"
// @Success 200 {object} utils.Response "Staff activated successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Staff not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /staff/profile/activate [patch]
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

// @Summary Activate user
// @Description Activate a user account. Admin access required.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response "User activated successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /users/profile/activate [patch]
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

// @Summary Update user profile
// @Description Update the authenticated user's profile information.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.UpdateProfileRequest true "Updated profile details"
// @Success 200 {object} utils.Response{data=dto.UserResponse} "Profile updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /users/profile [patch]
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

// @Summary Update staff profile
// @Description Update the authenticated staff member's profile information.
// @Tags Staff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.UpdateProfileRequest true "Updated profile details"
// @Success 200 {object} utils.Response{data=dto.UserResponse} "Profile updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Staff not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /staff/profile [patch]
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

// @Summary Search users
// @Description Search users using supported filter parameters. Admin access required.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param query query string false "Search keyword"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]dto.UserResponse} "Users retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid search parameters"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "No users found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /users/search [get]
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
