package server

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Register a new user
// @Description Create a new user account with email,firstname, lastname and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 201 {object} utils.Response{data=dto.AuthResponse} "User registered successfully"
// @Failure 400 {object} utils.Response "Invalid request data or user already exists"
// @Failure 429 {object} utils.Response "Too many requests"
// @Router /auth/register [post]
func (s *Server) RegisterUserHandler(ctx *gin.Context) {
	var req dto.RegisterRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	response, err := s.authServices.RegisterUserService(ctx.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrConflict:
			utils.ConflictResponse(ctx, "Email already in use", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
		}
		return
	}

	utils.CreatedResponse(ctx, "User created successfully", response)
}

// @Summary User login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login credentials"
// @Success 200 {object} utils.Response{data=dto.AuthResponse} "Login successful"
// @Failure 401 {object} utils.Response "Invalid credentials"
// @Failure 429 {object} utils.Response "Too many requests"
// @Router /auth/login [post]
func (s *Server) LoginUserHandler(ctx *gin.Context) {
	var req dto.LoginRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	response, err := s.authServices.LoginUserService(ctx.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			utils.UnAuthorized(ctx, "Invalid email or password", err)
		case domain.ErrAccountDeactivated:
			utils.Forbidden(ctx, "Account has been deactivated", err)
		case domain.ErrNotVerified:
			utils.Forbidden(ctx, "Please verify your email first", err)
		default:
			utils.UnAuthorized(ctx, "Login failed", err)
		}
		return
	}

	ctx.SetCookie(
		"refresh_token",
		response.RefreshToken,
		int(s.cfg.JWT.JWTRefreshTokenExpiration),
		"/",
		"",
		false,
		true,
	)

	response.RefreshToken = ""

	utils.SuccessResponse(ctx, "User loggedIn successfully", response)
}

// @Summary Refresh access token
// @Description Generate a new access token using a valid refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=dto.AuthResponse} "New access token generated"
// @Failure 401 {object} utils.Response "Invalid or expired refresh token"
// @Router /auth/refresh [post]
func (s *Server) RefreshTokenHandler(ctx *gin.Context) {
	token, err := ctx.Cookie("refresh_token")
	if err != nil {
		utils.UnAuthorized(ctx, "No refresh token provided", err)
		return
	}

	response, err := s.authServices.GenerateRefreshTokenService(ctx.Request.Context(), token)
	if err != nil {
		switch err {
		case domain.ErrInvalidRefreshToken, domain.ErrTokeNotFoundOrExpired:
			utils.UnAuthorized(ctx, "Invalid or expired refresh token", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
		}
		return
	}

	ctx.SetCookie(
		"refresh_token",
		response.RefreshToken,
		int(s.cfg.JWT.JWTRefreshTokenExpiration),
		"/",
		"",
		false,
		true,
	)

	response.RefreshToken = ""

	utils.SuccessResponse(ctx, "Success", response)
}

// @Summary User logout
// @Description Invalidate refresh token and logout user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token to invalidate"
// @Success 200 {object} utils.Response "Logout successful"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Router /auth/logout [post]
func (s *Server) LogoutHandler(ctx *gin.Context) {
	token, err := ctx.Cookie("refresh_token")
	if err != nil {
		utils.UnAuthorized(ctx, "No refresh token provided", err)
		return
	}

	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	err = s.authServices.LogoutService(ctx.Request.Context(), token)
	if err != nil {
		utils.BadRequest(ctx, "Something went wrong", err)
		return
	}

	ctx.SetCookie("refresh_token", "", -1, "/", "", false, true)

	utils.SuccessResponse(ctx, "User LoggedOut successfully", nil)
}

// @Summary Verify email
// @Description Verify a user's email address using the verification token or OTP.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dto.VerifyEmailRequest true "Email verification request"
// @Success 200 {object} utils.Response "Email verified successfully"
// @Failure 400 {object} utils.Response "Invalid or expired verification token"
// @Failure 429 {object} utils.Response "Too many requests"
// @Router /auth/verify-email [post]
func (s *Server) VerifyEmailHandler(ctx *gin.Context) {
	var req dto.VerifyEmailRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	response, err := s.authServices.VerifyEmailService(&req)
	if err != nil {
		switch err {
		case domain.ErrTokeNotFoundOrExpired:
			utils.BadRequest(ctx, "Invalid or expired verification token", err)
		case domain.ErrAlreadyVerified:
			utils.BadRequest(ctx, "Email already verified", err)
		case domain.ErrUserNotFound:
			utils.NotFound(ctx, "User not found", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Verification Success", response)
}

// @Summary Forgot password
// @Description Send a password reset code to the user's registered email address.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dto.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} utils.Response "Password reset code sent successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 429 {object} utils.Response "Too many requests"
// @Router /auth/forgot-password [post]
func (s *Server) ForgotPasswordHandler(ctx *gin.Context) {
	var req dto.ForgotPasswordRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	err = s.authServices.ForgotPasswordService(ctx.Request.Context(), &req)
	if err != nil {
		utils.BadRequest(ctx, "Something went wrong", err)
		return
	}

	utils.SuccessResponse(ctx, "Password reset code send to Mail, Check Your Inbox", nil)
}

// @Summary Verify password reset token
// @Description Verify that a password reset token or OTP is valid before allowing the user to reset their password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dto.VerifyResetToken true "Password reset token verification request"
// @Success 200 {object} utils.Response "Password reset token verified successfully"
// @Failure 400 {object} utils.Response "Invalid or expired password reset token"
// @Failure 429 {object} utils.Response "Too many requests"
// @Router /auth/verify-reset-token [post]
func (s *Server) VerifyResetOTPHandler(ctx *gin.Context) {
	var req dto.VerifyResetToken

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	err = s.authServices.VerifyResetOTP(&req)
	if err != nil {
		switch err {
		case domain.ErrTokeNotFoundOrExpired:
			utils.BadRequest(ctx, "Invalid or expired verification token", err)
		case domain.ErrUserNotFound:
			utils.NotFound(ctx, "User not found", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Verification Success", nil)
}

// @Summary User reset password
// @Description Reset user password with reset token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset password data"
// @Success 200 {object} utils.Response "Password reset successful"
// @Failure 400 {object} utils.Response "Invalid request data or token expired"
// @Failure 429 {object} utils.Response "Too many requests"
// @Router /auth/reset-password [post]
func (s *Server) ResetPasswordHandler(ctx *gin.Context) {
	var req dto.ResetPasswordRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	err = s.authServices.ResetPasswordService(&req)
	if err != nil {
		switch err {
		case domain.ErrTokeNotFoundOrExpired:
			utils.BadRequest(ctx, "Invalid or expired verification token", err)
		case domain.ErrAlreadyVerified:
			utils.BadRequest(ctx, "Email already verified", err)
		case domain.ErrUserNotFound:
			utils.NotFound(ctx, "User not found", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Success", nil)
}

// @Summary User change password
// @Description Change user password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordRequest true "Change password data"
// @Success 200 {object} utils.Response "Password changed successfully"
// @Failure 400 {object} utils.Response "Invalid request data or current password is incorrect"
// @Security BearerAuth
// @Router /auth/password-change [post]
func (s *Server) ChangePasswordHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var req dto.ChangePasswordRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	err = s.authServices.ChangePasswordService(ctx.Request.Context(), userID, &req)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			utils.UnAuthorized(ctx, "Current password is incorrect", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Password changed successfully", nil)
}
