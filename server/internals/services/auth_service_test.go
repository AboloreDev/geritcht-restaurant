package services

import (
	"context"
	"testing"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/stretchr/testify/assert"
)

var testCtx = context.Background()

// ─── MockUserRepository
type MockUserRepository struct {
	user      *models.User
	getErr    error
	createErr error
	updateErr error
	users     []models.User
	userSearchRank []models.UserWithRank
	total     int64
}

type MockAuthRepository struct {
	token *models.Token
	err   error
}

type MockPublisher struct{ err error }

// ─── Helper
func newAuthService(
	userRepo *MockUserRepository,
	authRepo *MockAuthRepository,
	cartRepo *MockCartRepository,
	publisher *MockPublisher,
) *AuthService {
	return NewAuthService(testAuthConfig, publisher, userRepo, authRepo, cartRepo)
}

func (m *MockUserRepository) GetByEmail(_ context.Context, email string) (*models.User, error) {
	return m.user, m.getErr
}
func (m *MockUserRepository) GetByID(_ context.Context, id uint) (*models.User, error) {
	return m.user, m.getErr
}
func (m *MockUserRepository) Create(_ context.Context, user *models.User) error {
	user.ID = 1
	return m.createErr
}
func (m *MockUserRepository) Update(_ context.Context, user *models.User) error { return m.updateErr }
func (m *MockUserRepository) Delete(_ context.Context, id uint) error           { return m.getErr }
func (m *MockUserRepository) GetAll(_ context.Context) ([]*models.User, error) {
	return nil, m.getErr
}
func (m *MockUserRepository) GetByIdAndActive(_ context.Context, id uint, active bool) (*models.User, error) {
	return m.user, m.getErr
}
func (m *MockUserRepository) UpdateActiveByRole(ctx context.Context, id uint, role models.UserRole, active bool) error {
	return m.updateErr
}
func (m *MockUserRepository) GetByIDAndRole(ctx context.Context, id uint, role models.UserRole) (*models.User, error) {
	return m.user, m.getErr
}
func (m *MockUserRepository) GetAllByRole(ctx context.Context, role models.UserRole, page, pageSize int) ([]models.User, int64, error) {
	return m.users, m.total, m.getErr
}

func (m *MockAuthRepository) CreateRefreshToken(userID uint, token string) error { return m.err }
func (m *MockAuthRepository) GetValidRefreshToken(_ context.Context, token string) (*models.RefreshToken, error) {
	return &models.RefreshToken{ID: 1, UserID: 1}, m.err
}
func (m *MockAuthRepository) GetRefreshToken(hash string) (*models.RefreshToken, error) {
	return &models.RefreshToken{ID: 1}, m.err
}
func (m *MockAuthRepository) DeleteRefreshToken(token *models.RefreshToken) error { return m.err }
func (m *MockAuthRepository) DeleteExpiredRefreshTokens(userID uint) error        { return m.err }
func (m *MockAuthRepository) CreateEmailToken(token *models.Token) error          { return m.err }
func (m *MockAuthRepository) GetValidEmailToken(hash string) (*models.Token, error) {
	return m.token, m.err
}
func (m *MockAuthRepository) VerifyUserEmail(user *models.User, token *models.Token) error {
	return m.err
}
func (m *MockAuthRepository) ResetPassword(user *models.User, token *models.Token, password string) error {
	return m.err
}
func (m *MockAuthRepository) ChangePassword(user *models.User, password string) error { return m.err }
func (m *MockUserRepository) TsvectorSearchUsers(_ context.Context, req *dto.UserSearchRequest) ([]models.UserWithRank, int64, error) {
	return m.userSearchRank, m.total, m.getErr
}

// ─── MockPublisher

func (m *MockPublisher) PublishMessage(eventType string, payload interface{}, metadata map[string]string) error {
	return m.err
}

func (m *MockPublisher) CloseMessage() error { return nil }

// ─── Config

var testAuthConfig = &config.Config{
	JWT: config.JWTConfig{
		JWTSecret:                 "test-secret-key",
		JWTTokenExpiration:        15 * time.Minute,
		JWTRefreshTokenExpiration: 7 * 24 * time.Hour,
	},
}

func TestRegisterUser_Success(t *testing.T) {
	service := newAuthService(
		&MockUserRepository{
			getErr:    domain.ErrUserNotFound,
			createErr: nil},
		&MockAuthRepository{},
		&MockCartRepository{},
		&MockPublisher{},
	)

	req := &dto.RegisterRequest{
		Email:       "test@test.com",
		Password:    "password123",
		FirstName:   "Test",
		LastName:    "User",
		PhoneNumber: "08012345678",
	}

	response, err := service.RegisterUserService(testCtx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "test@test.com", response.User.Email)
	assert.False(t, response.User.EmailVerified)
}

func TestRegisterUser_DuplicateEmail(t *testing.T) {
	service := newAuthService(
		&MockUserRepository{
			user:   &models.User{ID: 1, Email: "existing@test.com"},
			getErr: nil,
		},
		&MockAuthRepository{},
		&MockCartRepository{},
		&MockPublisher{},
	)

	req := &dto.RegisterRequest{
		Email:    "existing@test.com",
		Password: "password123",
	}

	response, err := service.RegisterUserService(testCtx, req)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrConflict, err)
}

func TestRegisterUser_CartCreationFails(t *testing.T) {
	service := newAuthService(
		&MockUserRepository{createErr: nil},
		&MockAuthRepository{},
		&MockCartRepository{getCartErr: domain.ErrCartNotFound},
		&MockPublisher{},
	)

	req := &dto.RegisterRequest{
		Email:     "test@test.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	response, err := service.RegisterUserService(testCtx, req)

	assert.Nil(t, response)
	assert.Error(t, err)
}

func TestRegisterUser_PublisherFails(t *testing.T) {
	service := newAuthService(
		&MockUserRepository{createErr: nil},
		&MockAuthRepository{},
		&MockCartRepository{},
		&MockPublisher{err: domain.ErrPublishing},
	)

	req := &dto.RegisterRequest{
		Email:     "test@test.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	response, err := service.RegisterUserService(testCtx, req)

	assert.Nil(t, response)
	assert.Error(t, err)
}

// ─── LoginUser

func TestLoginUser(t *testing.T) {
	hashedPassword, _ := utils.HashPassword("password123")

	tests := []struct {
		name        string
		user        *models.User
		getErr      error
		password    string
		expectedErr error
	}{
		{
			name: "success",
			user: &models.User{
				ID:            1,
				Email:         "test@test.com",
				Password:      hashedPassword,
				Role:          models.RoleCustomer,
				IsActive:      true,
				EmailVerified: true,
			},
			getErr:      nil,
			password:    "password123",
			expectedErr: nil,
		},
		{
			name:        "user not found",
			user:        nil,
			getErr:      domain.ErrUserNotFound,
			password:    "password123",
			expectedErr: domain.ErrInvalidCredentials,
		},
		{
			name: "wrong password",
			user: &models.User{
				ID:            1,
				Password:      hashedPassword,
				IsActive:      true,
				EmailVerified: true,
			},
			getErr:      nil,
			password:    "wrongpassword",
			expectedErr: domain.ErrInvalidCredentials,
		},
		{
			name: "email not verified",
			user: &models.User{
				ID:            1,
				Password:      hashedPassword,
				IsActive:      true,
				EmailVerified: false,
			},
			getErr:      nil,
			password:    "password123",
			expectedErr: domain.ErrNotVerified,
		},
		{
			name: "account deactivated",
			user: &models.User{
				ID:            1,
				Password:      hashedPassword,
				IsActive:      false,
				EmailVerified: true,
			},
			getErr:      nil,
			password:    "password123",
			expectedErr: domain.ErrAccountDeactivated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newAuthService(
				&MockUserRepository{user: tt.user, getErr: tt.getErr},
				&MockAuthRepository{},
				&MockCartRepository{},
				&MockPublisher{},
			)

			req := &dto.LoginRequest{
				Email:    "test@test.com",
				Password: tt.password,
			}

			response, err := service.LoginUserService(testCtx, req)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, response)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.NotEmpty(t, response.AccessToken)
			assert.NotEmpty(t, response.RefreshToken)
		})
	}
}

// ─── VerifyEmail

func TestVerifyEmail(t *testing.T) {
	tests := []struct {
		name        string
		token       *models.Token
		tokenErr    error
		user        *models.User
		getErr      error
		expectedErr error
		expectedVal bool
	}{
		{
			name: "success",
			token: &models.Token{
				ID:     1,
				UserID: 1,
				Type:   models.TokenTypeEmailVerification,
			},
			user: &models.User{
				ID:            1,
				EmailVerified: false,
			},
			getErr:      nil,
			expectedErr: nil,
			expectedVal: true,
		},
		{
			name:        "invalid token",
			token:       nil,
			tokenErr:    domain.ErrTokeNotFoundOrExpired,
			expectedErr: domain.ErrTokeNotFoundOrExpired,
			expectedVal: false,
		},
		{
			name: "already verified",
			token: &models.Token{
				ID:     1,
				UserID: 1,
			},
			user: &models.User{
				ID:            1,
				EmailVerified: true,
			},
			getErr:      nil,
			expectedErr: domain.ErrAlreadyVerified,
			expectedVal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newAuthService(
				&MockUserRepository{user: tt.user, getErr: tt.getErr},
				&MockAuthRepository{token: tt.token, err: tt.tokenErr},
				&MockCartRepository{},
				&MockPublisher{},
			)

			req := &dto.VerifyEmailRequest{Token: "sometoken"}
			result, err := service.VerifyEmailService(req)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedVal, result)
		})
	}
}

// ─── ChangePassword

func TestChangePassword(t *testing.T) {
	hashedPassword, _ := utils.HashPassword("currentpassword")

	tests := []struct {
		name        string
		user        *models.User
		getErr      error
		authErr     error
		req         *dto.ChangePasswordRequest
		expectedErr error
	}{
		{
			name: "success",
			user: &models.User{
				ID:       1,
				Password: hashedPassword,
			},
			getErr: nil,
			req: &dto.ChangePasswordRequest{
				CurrentPassword: "currentpassword",
				NewPassword:     "newpassword123",
			},
			expectedErr: nil,
		},
		{
			name:   "user not found",
			user:   nil,
			getErr: domain.ErrUserNotFound,
			req: &dto.ChangePasswordRequest{
				CurrentPassword: "currentpassword",
				NewPassword:     "newpassword123",
			},
			expectedErr: domain.ErrUserNotFound,
		},
		{
			name: "wrong current password",
			user: &models.User{
				ID:       1,
				Password: hashedPassword,
			},
			getErr: nil,
			req: &dto.ChangePasswordRequest{
				CurrentPassword: "wrongpassword",
				NewPassword:     "newpassword123",
			},
			expectedErr: domain.ErrInvalidCredentials,
		},
		{
			name: "weak new password",
			user: &models.User{
				ID:       1,
				Password: hashedPassword,
			},
			getErr: nil,
			req: &dto.ChangePasswordRequest{
				CurrentPassword: "currentpassword",
				NewPassword:     "short",
			},
			expectedErr: domain.ErrWeakPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newAuthService(
				&MockUserRepository{user: tt.user, getErr: tt.getErr},
				&MockAuthRepository{err: tt.authErr},
				&MockCartRepository{},
				&MockPublisher{},
			)

			err := service.ChangePasswordService(testCtx, 1, tt.req)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

// ─── ForgotPassword

func TestForgotPassword(t *testing.T) {
	tests := []struct {
		name        string
		user        *models.User
		getErr      error
		expectedErr error
	}{
		{
			name: "success",
			user: &models.User{
				ID:        1,
				Email:     "test@test.com",
				FirstName: "Test",
			},
			getErr:      nil,
			expectedErr: nil,
		},
		{
			name:        "email not found returns nil", // security — don't reveal
			user:        nil,
			getErr:      domain.ErrUserNotFound,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newAuthService(
				&MockUserRepository{user: tt.user, getErr: tt.getErr},
				&MockAuthRepository{},
				&MockCartRepository{},
				&MockPublisher{},
			)

			req := &dto.ForgotPasswordRequest{Email: "test@test.com"}
			err := service.ForgotPasswordService(testCtx, req)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
