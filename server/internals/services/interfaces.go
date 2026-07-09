package services

import (
	"context"

	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
)

// Allergen Service
type AllergenServiceInterface interface {
	CreateAllergenServices(ctx context.Context, req *dto.CreateAllergenRequest) (*dto.AllergenResponse, error)
	GetAllAllergenService(ctx context.Context, page, pageSize int) ([]*dto.AllergenResponse, *utils.PaginatedMeta, error) 
	UpdateAllergenService(ctx context.Context, allergenID uint, req *dto.UpdateAllergenRequest) (*dto.AllergenResponse, error)
	DeleteAllergenService(ctx context.Context, allergenID uint) error
}

// Auth Service 
type AuthServiceInterface interface {
	RegisterUserService(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	VerifyEmailService(req *dto.VerifyEmailRequest) (bool, error)
	LoginUserService(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	GenerateRefreshTokenService(ctx context.Context, refresh string) (*dto.AuthResponse, error)
	 LogoutService(ctx context.Context, refresh string) error 
	 ForgotPasswordService(ctx context.Context, req *dto.ForgotPasswordRequest) error 
	 VerifyResetOTP(req *dto.VerifyResetToken) error
	 ResetPasswordService(req *dto.ResetPasswordRequest) error 
	 ChangePasswordService(ctx context.Context, userID uint, req *dto.ChangePasswordRequest) error
	 GenerateAuthResponse(user *models.User) (*dto.AuthResponse, error)
}

// Menu
type MenuServiceInterface interface {
	SearchProduct(ctx context.Context, req *dto.MenuSearchRequest) ([]*dto.MenuSearchResponse, *utils.PaginatedMeta, error) 
	 GetAllMenuService(ctx context.Context, filter dto.MenuFilterRequest) ([]*dto.MenuResponse, *utils.PaginatedMeta, error)
	 ToggleMenuAvailabilityService(ctx context.Context, menuID uint, isAvailable *bool) error 
	 RemoveMenuImageService(ctx context.Context, menuImageID uint) error
	 AddMenuImageService(ctx context.Context, menuID uint, altText, url string) error
	 DeleteMenu(ctx context.Context, menuID uint) error
	 GetMenu(ctx context.Context, menuID uint) (*dto.MenuResponse, error) 
	 UpdateMenuService(ctx context.Context, menuID uint, req *dto.UpdateMenuRequest) (*dto.MenuResponse, error)
	  CreateMenuService(ctx context.Context, req *dto.CreateMenuRequest) (*dto.MenuResponse, error)
}

// Cart 
type CartServiceInterface interface {
	GetUserCart(ctx context.Context, userID uint) (*dto.CartResponse, error)
	 AddItemToCart(ctx context.Context, userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) 
	 UpdateCartItem(ctx context.Context, userID uint, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error) 
	 RemoveCartItem(ctx context.Context, userID uint, itemID uint) error
	  ClearCart(ctx context.Context, userID uint) error 
}

// Category Service