package services

import (
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthService(
	db *gorm.DB,
	cfg *config.Config) *AuthService {
	return &AuthService{
		db:  db,
		cfg: cfg,
	}
}

func (s *AuthService) RegisterUserService(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	var user models.User

	err := s.db.Where("email = ? ", req.Email).First(&user).Error
	if err == nil {
		return nil, domain.ErrConflict
	}

	hashedPassword, _ := utils.HashPassword(req.Password)

	user = models.User{
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Email:         req.Email,
		Password:      hashedPassword,
		PhoneNumber:   req.PhoneNumber,
		Role:          models.RoleCustomer,
		EmailVerified: false,
	}

	err = s.db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	cart := models.Cart{
		UserID: user.ID,
	}

	err = s.db.Create(&cart).Error
	if err != nil {
		return nil, err
	}

	verificationToken, _ := utils.GenerateVerificationToken()

	hashedToken := utils.HashToken(verificationToken)

	token := models.Token{
		UserID:    user.ID,
		TokenHash: hashedToken,
		Type:      models.TokenTypeEmailVerification,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = s.db.Create(&token).Error
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			PhoneNumber:   user.PhoneNumber,
			IsActive:      user.IsActive,
			Role:          string(user.Role),
			EmailVerified: user.EmailVerified,
		},
	}, nil
}

func (s *AuthService) VerifyEmailService(req *dto.VerifyEmailRequest) (*dto.AuthResponse, error) {
	var token models.Token

	hashedToken := utils.HashToken(req.Token)

	err := s.db.
		Where("token_hash = ? AND expires_at > ?", hashedToken, time.Now()).
		First(&token).Error
	if err != nil {
		return nil, domain.ErrTokeNotFoundOrExpired
	}

	var user models.User
	if err := s.db.First(&user, token.UserID).Error; err != nil {
		return nil, domain.ErrUserNotFound
	}

	if user.EmailVerified {
		return nil, domain.ErrAlreadyVerified
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).
			Update("email_verified", true).Error; err != nil {
			return err
		}

		if err := tx.Delete(&token).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GenerateAuthResponse(&user)
}

func (s *AuthService) LoginUserService(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	var user models.User

	err := s.db.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, domain.ErrAccountDeactivated
	}

	if !user.EmailVerified {
		return nil, domain.ErrNotVerified
	}

	ok := utils.CheckPassword(user.Password, req.Password)
	if !ok {
		return nil, domain.ErrInvalidCredentials
	}

	return s.GenerateAuthResponse(&user)

}

func (s *AuthService) GenerateRefreshTokenService(refresh string) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(refresh, s.cfg.JWT.JWTSecret)
	if err != nil {
		return nil, domain.ErrInvalidRefreshToken
	}

	var refreshToken models.RefreshToken
	hashedToken := utils.HashToken(refresh)
	err = s.db.Where("token_hash = ? AND expires_at > ?", hashedToken, time.Now()).First(&refreshToken).Error
	if err != nil {
		return nil, domain.ErrTokeNotFoundOrExpired
	}

	var user models.User
	err = s.db.First(&user, claims.UserID).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	s.db.Delete(&refreshToken)

	return s.GenerateAuthResponse(&user)
}

func (s *AuthService) LogoutService(refresh string) error {
	var refreshToken models.RefreshToken

	err := s.db.Where("token_hash = ?", refresh).Delete(&refreshToken).Error
	if err != nil {
		return nil
	}

	return nil
}

func (s *AuthService) ForgotPasswordService(req *dto.ForgotPasswordRequest) error {
	var user models.User
	var token models.Token

	err := s.db.Where("email = ? ", req.Email).First(&user).Error
	if err != nil {
		return nil
	}

	passwordResetToken, _ := utils.GeneratePasswordResetToken()

	hashedToken := utils.HashToken(passwordResetToken)

	token = models.Token{
		UserID:    user.ID,
		TokenHash: hashedToken,
		Type:      models.TokenTypePasswordReset,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	err = s.db.Create(&token).Error
	if err != nil {
		return nil
	}
	// TODO: Publisher to publish to email microservice

	return nil
}

func (s *AuthService) ResetPasswordService(req *dto.ResetPasswordRequest) error {
	var token models.Token
	var user models.User

	hashedToken := utils.HashToken(req.Token)

	err := s.db.
		Where("token_hash = ? AND expires_at > ?", hashedToken, time.Now()).
		First(&token).Error
	if err != nil {
		return domain.ErrTokeNotFoundOrExpired
	}

	if err := s.db.First(&user, token.UserID).Error; err != nil {
		return domain.ErrUserNotFound
	}

	if len(req.NewPassword) < 8 {
		return domain.ErrWeakPassword
	}

	newHashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).
			Update("password", newHashedPassword).Error; err != nil {
			return err
		}

		if err := tx.Delete(&token).Error; err != nil {
			return err
		}

		if err := tx.Where("user_id = ?", user.ID).
			Delete(&models.RefreshToken{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *AuthService) ChangePasswordService(userID uint, req *dto.ChangePasswordRequest) error {
	var user models.User
	var refreshToken models.RefreshToken

	err := s.db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return err
	}

	ok := utils.CheckPassword(user.Password, req.CurrentPassword)
	if !ok {
		return ErrInvalidCredentials
	}

	newHashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&user).Update("password", newHashedPassword).Error
		if err != nil {
			return err
		}

		if err := tx.Where("user_id = ?", user.ID).
			Delete(&refreshToken).Error; err != nil {
			return err
		}

		return nil

	})

}

func (s *AuthService) GenerateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(&s.cfg.JWT, user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshToken,
		ExpiresAt: time.Now().Add(s.cfg.JWT.JWTRefreshTokenExpiration),
	}

	s.db.Create(&refreshTokenModel)

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			PhoneNumber:   user.PhoneNumber,
			IsActive:      user.IsActive,
			Role:          string(user.Role),
			EmailVerified: user.EmailVerified,
		},
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}
