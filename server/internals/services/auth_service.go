package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/events"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
)

type AuthService struct {
	cfg       *config.Config
	publisher interfaces.Publisher
	userRepo  repositories.UserRepositoryInterface
	authRepo  repositories.AuthRepositoryInterface
	cartRepo  repositories.CartRepositoryInterface
}

func NewAuthService(
	cfg *config.Config,
	publisher interfaces.Publisher,
	userRepo repositories.UserRepositoryInterface,
	authRepo repositories.AuthRepositoryInterface,
	cartRepo repositories.CartRepositoryInterface) *AuthService {
	return &AuthService{
		cfg:       cfg,
		publisher: publisher,
		userRepo:  userRepo,
		authRepo:  authRepo,
		cartRepo:  cartRepo,
	}
}

func (s *AuthService) RegisterUserService(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, domain.ErrConflict
	}

	hashedPassword, _ := utils.HashPassword(req.Password)

	user := models.User{
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Email:         req.Email,
		Password:      hashedPassword,
		PhoneNumber:   req.PhoneNumber,
		Role:          models.RoleCustomer,
		EmailVerified: false,
	}

	err = s.userRepo.Create(ctx, &user)
	if err != nil {
		return nil, err
	}

	cart := models.Cart{
		UserID: user.ID,
	}

	err = s.cartRepo.CreateCart(ctx, &cart)
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

	err = s.authRepo.CreateEmailToken(&token)
	if err != nil {
		return nil, err
	}

	// Publish message to the queue
	err = s.publisher.PublishMessage(
		events.ChannelEmailVerification,
		events.VerificationEmailPayload{
			Email:     user.Email,
			FirstName: user.FirstName,
			Token:     verificationToken,
		},
		map[string]string{"Priority": "Important Mail"},
	)
	if err != nil {
		log.Printf("Failed to put messages in queue: %v", err)
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

func (s *AuthService) VerifyEmailService(req *dto.VerifyEmailRequest) (bool, error) {

	hashedToken := utils.HashToken(req.Token)

	token, err := s.authRepo.GetValidEmailToken(hashedToken)
	if err != nil {
		return false, domain.ErrTokeNotFoundOrExpired
	}

	user, err := s.userRepo.GetByID(context.Background(), token.UserID)
	if err != nil {
		return false, domain.ErrUserNotFound
	}

	if user.EmailVerified {
		return false, domain.ErrAlreadyVerified
	}

	err = s.authRepo.VerifyUserEmail(user, token)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *AuthService) LoginUserService(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
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

	return s.GenerateAuthResponse(user)
}

func (s *AuthService) GenerateRefreshTokenService(ctx context.Context, refresh string) (*dto.AuthResponse, error) {

	claims, err := utils.ValidateToken(refresh, s.cfg.JWT.JWTSecret)
	if err != nil {
		return nil, domain.ErrInvalidRefreshToken
	}

	refreshToken, err := s.authRepo.GetValidRefreshToken(ctx, refresh)
	if err != nil {
		return nil, domain.ErrTokeNotFoundOrExpired
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	err = s.authRepo.DeleteRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to delete refresh token: %w",
			err,
		)
	}

	return s.GenerateAuthResponse(user)
}

func (s *AuthService) LogoutService(ctx context.Context, refresh string) error {
	refreshToken, err := s.authRepo.GetValidRefreshToken(ctx, refresh)
	if err != nil {
		return err
	}
	err = s.authRepo.DeleteRefreshToken(refreshToken)

	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ForgotPasswordService(ctx context.Context, req *dto.ForgotPasswordRequest) error {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil
	}

	passwordResetToken, _ := utils.GenerateOTP()

	hashedToken := utils.HashToken(passwordResetToken)

	token := models.Token{
		UserID:    user.ID,
		TokenHash: hashedToken,
		Type:      models.TokenTypePasswordReset,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	err = s.authRepo.CreateEmailToken(&token)
	if err != nil {
		return err
	}

	// Publish message to queue
	err = s.publisher.PublishMessage(
		events.ChannelEmailPasswordReset,
		events.PasswordResetEmailPayload{
			Email:     user.Email,
			FirstName: user.FirstName,
			Token:     passwordResetToken,
		},
		map[string]string{"Priority": "Important Mail"},
	)

	if err != nil {
		log.Printf("Failed to put messages in queue: %v", err)
		return err
	}

	return nil
}

func (s *AuthService) VerifyResetOTP(req *dto.VerifyResetToken) error {

	hashedToken := utils.HashToken(req.Token)

	_, err := s.authRepo.GetValidEmailToken(hashedToken)
	if err != nil {
		return domain.ErrTokeNotFoundOrExpired
	}

	return nil
}

func (s *AuthService) ResetPasswordService(req *dto.ResetPasswordRequest) error {

	hashedToken := utils.HashToken(req.Token)

	token, err := s.authRepo.GetValidEmailToken(hashedToken)
	if err != nil {
		return domain.ErrTokeNotFoundOrExpired
	}

	user, err := s.userRepo.GetByID(context.Background(), token.UserID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	newHashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	err = s.authRepo.ResetPassword(user, token, newHashedPassword)

	if err != nil {
		return err
	}

	// Publish password changed event
	if err := s.publisher.PublishMessage(
		events.ChannelEmailPasswordChanged,
		events.PasswordChangedEmailPayload{
			Email:     user.Email,
			FirstName: user.FirstName,
		},
		map[string]string{
			"Priority": "Low Priority",
		},
	); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ChangePasswordService(ctx context.Context, userID uint, req *dto.ChangePasswordRequest) error {

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	ok := utils.CheckPassword(user.Password, req.CurrentPassword)
	if !ok {
		return domain.ErrInvalidCredentials
	}

	if len(req.NewPassword) < 8 {
		return domain.ErrWeakPassword
	}

	newHashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	err = s.authRepo.ChangePassword(user, newHashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GenerateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(&s.cfg.JWT, user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	err = s.authRepo.CreateRefreshToken(user.ID, refreshToken)
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
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}
