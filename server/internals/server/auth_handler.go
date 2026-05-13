package server

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterUserHandler(ctx *gin.Context) {
	var req dto.RegisterRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	response, err := s.authServices.RegisterUserService(&req)
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

func (s *Server) LoginUserHandler(ctx *gin.Context) {
	var req dto.LoginRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	response, err := s.authServices.LoginUserService(&req)
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

func (s *Server) RefreshTokenHandler(ctx *gin.Context) {
	token, err := ctx.Cookie("refresh_token")
	if err != nil {
		utils.UnAuthorized(ctx, "No refresh token provided", err)
		return
	}

	response, err := s.authServices.GenerateRefreshTokenService(token)
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

	err = s.authServices.LogoutService(token)
	if err != nil {
		utils.BadRequest(ctx, "Something went wrong", err)
		return
	}

	ctx.SetCookie("refresh_token", "", -1, "/", "", false, true)

	utils.SuccessResponse(ctx, "User LoggedOut successfully", nil)
}

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

func (s *Server) ForgotPasswordHandler(ctx *gin.Context) {
	var req dto.ForgotPasswordRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	err = s.authServices.ForgotPasswordService(&req)
	if err != nil {
		utils.BadRequest(ctx, "Something went wrong", err)
		return
	}

	utils.SuccessResponse(ctx, "Password reset code send to Mail, Check Your Inbox", nil)
}

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

func (s *Server) ChangePasswordHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var req dto.ChangePasswordRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Request Data", err)
		return
	}

	err = s.authServices.ChangePasswordService(userID, &req)
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
