package services

import (
	"context"
	"mime/multipart"

	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
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
type CategoryServiceInterface interface {
	SearchCategory(ctx context.Context, req *dto.CategorySearchRequest) ([]*dto.CategorySearchResponse, *utils.PaginatedMeta, error)
	GetCategoryService(ctx context.Context, categoryID uint) (*dto.MenuCategoryResponse, error)
	GetCategoriesService(ctx context.Context, page, pageSize int) ([]*dto.MenuCategoryResponse, *utils.PaginatedMeta, error)
	DeleteCategoryService(ctx context.Context, categoryID uint) error
	UpdateCategoryService(ctx context.Context, categoryID uint, req *dto.UpdateCategoryRequest) (*dto.MenuCategoryResponse, error)
	CreateCategoryService(ctx context.Context, req *dto.CreateCategoryRequest, imageURL string) (*dto.MenuCategoryResponse, error)
}

// Dietary tags
type DietaryTagsServiceInterface interface {
	CreateDietaryTagService(ctx context.Context, req *dto.CreateDietaryTagRequest) (*dto.DietaryTagResponse, error)
	GetAllDietaryTagService(ctx context.Context, page, pageSize int) ([]*dto.DietaryTagResponse, *utils.PaginatedMeta, error)
	UpdateDietaryTagService(ctx context.Context, tagID uint, req *dto.UpdateDietaryTagRequest) (*dto.DietaryTagResponse, error)
	DeleteDietaryTagService(ctx context.Context, tagID uint) error
}

// ingredient
type IngredientServiceInterface interface {
	CreateIngredientService(ctx context.Context, req *dto.CreateIngredientRequest) (*dto.IngredientResponse, error)
	GetAllIngredientService(ctx context.Context, page, pageSize int) ([]*dto.IngredientResponse, *utils.PaginatedMeta, error)
	UpdateIngredientService(ctx context.Context, ingredientID uint, req *dto.UpdateIngredientRequest) (*dto.IngredientResponse, error)
	DeleteIngredientService(ctx context.Context, ingredientID uint) error
	GetIngredientService(ctx context.Context, ingredientID uint) (*dto.IngredientResponse, error)
	GetLowStockIngredientsService(ctx context.Context) ([]*dto.IngredientResponse, error)
	SetThresholdLimit(ctx context.Context, ingredientID uint, req *dto.ThresholdRequest) error
	CheckLowStock(ctx context.Context, userID, ingredientID uint) error
	sendLowStockAlert(ctx context.Context, user *models.User, ingredient *models.Ingredient) error
	SearchIngredients(ctx context.Context, req *dto.IngredientSearchRequest) ([]*dto.IngredientSearchResponse, *utils.PaginatedMeta, error)
}

// Inventory interface
type InventoryServiceInterface interface {
	DeductStock(ctx context.Context, tx *gorm.DB, orderItems []models.OrderItem, orderID uint, createdBy uint) error
	CheckAndAlertThreshold(ctx context.Context, tx *gorm.DB) error
}

// Recipes
type ReceipesServiceInterface interface {
	GetAllRecipes(ctx context.Context, menuItemID uint) ([]*dto.MenuItemIngredientResponse, error)
	DeleteRecipe(ctx context.Context, menuItemID uint, ingredientID uint) error
	UpdateMenuRecipe(ctx context.Context, menuItemID uint, ingredientID uint, req *dto.UpdateLinkItemRequest) (*dto.MenuItemIngredientResponse, error)
	AddMenuRecipe(ctx context.Context, menuItemID uint, req *dto.LinkIngredientRequest) (*dto.MenuItemIngredientResponse, error)
}

// Order Service
type OrderServiceInterface interface {
	CreateTakeoutOrder(ctx context.Context, userID uint, req *dto.CreateTakeoutOrderRequest) (*dto.OrderResponse, error)
	GetAllUserTakeoutOrders(ctx context.Context, userID uint, page, pageSize int) ([]*dto.OrderResponse, *utils.PaginatedMeta, error)
	GetTakeoutOrder(ctx context.Context, userID, orderID uint) (*dto.OrderResponse, error)
	CancelTakeoutOrder(ctx context.Context, userID, orderID uint) error
	VerifyUserOrder(ctx context.Context, userID, orderID uint) error
	GetAllOrders(ctx context.Context, page, pageSize int) ([]*dto.OrderResponse, *utils.PaginatedMeta, error)
}

// Payment Srvice
type PaymentServiceInterface interface {
	callPaystackInitialize(email string, amount int64, reference string, orderID uint) (*PaystackInitialiseResponse, error)
	callPaystackVerify(reference string) (*PaystackVerifyResponse, error)
	callPaystackRefund(reference string, amount int64) (*PaystackRefundResponse, error)
	verifySignature(body []byte, signature string) bool
	InitialisePayment(ctx context.Context, userID uint, req *dto.InitializePaymentRequest) (*dto.InitializePaymentResponse, error)
	HandlePaystackWebhook(ctx context.Context, body []byte, signature string) error
	ProcessTakeoutRefund(ctx context.Context, orderID uint, notes string) error
	GetPaymentByReference(ctx context.Context, reference string) (*dto.PaymentResponse, error)
	GetPaymentDetails(ctx context.Context, paymentID uint) (*dto.PaymentResponse, error)
	GetRefundDetails(ctx context.Context, refundID uint) (*dto.RefundResponse, error)
	GetAllPaymentHistory(ctx context.Context, userID uint, page, pageSize int) ([]*dto.PaymentResponse, *utils.PaginatedMeta, error)
	VerifyPayment(ctx context.Context, req *dto.VerifyPaymentRequest) (*dto.PaymentResponse, error)
}

// Reservation Service
type ReservationServiceInterface interface {
	CancelReservation(ctx context.Context, userID uint, reservationID uint) (*dto.ReservationResponse, error)
	CheckInReservation(ctx context.Context, reservationID uint, userID uint) (*dto.ReservationResponse, error)
	GetTodayReservations(ctx context.Context, req *dto.ReservationFilterRequest) (*dto.ReservationListResponse, error)
	GetAllReservations(ctx context.Context, req *dto.ReservationFilterRequest) (*dto.ReservationListResponse, error)
	GetUserReservation(ctx context.Context, userID uint, reservationID uint) (*dto.ReservationResponse, error)
	GetAllUserReservations(ctx context.Context, userID uint, req *dto.ReservationFilterRequest) (*dto.ReservationListResponse, error)
	CreateReservation(ctx context.Context, req *dto.CreateReservationRequest, userID uint) (*dto.ReservationResponse, error)
	CheckTableAvailability(ctx context.Context, req *dto.CheckAvailabilityRequest) (*dto.AvailabilityResponse, error)
	buildReservationListResponse(reservations []models.Reservation, count int64, req *dto.ReservationFilterRequest) *dto.ReservationListResponse
}

// Table service
type TableServiceInterface interface {
	GetAllTablesService(ctx context.Context, page, pageSize int) ([]*dto.TableResponse, *utils.PaginatedMeta, error)
	DeleteTableService(ctx context.Context, tableID uint) error
	GetTableService(ctx context.Context, tableID uint) (*dto.TableDetailResponse, error)
	UpdateTableStatusService(ctx context.Context, tableID uint, req *dto.UpdateTableStatusRequest) (*dto.TableResponse, error)
	UpdateTableService(ctx context.Context, tableID uint, req *dto.UpdateTableRequest) (*dto.TableResponse, error)
	CreateTableService(ctx context.Context, req *dto.CreateTableRequest) (*dto.TableResponse, error)
}

// Uplaod Service interface
type UplaodServiceInterface interface {
	UploadMenuImage(menuID uint, file *multipart.FileHeader) (string, error)
	UploadCategoryImage(file *multipart.FileHeader) (string, error)
	DeleteFile(menuID uint) error
}

// User service
type UserServiceInterface interface {
	GetUserProfileService(ctx context.Context, userID uint) (*dto.UserResponse, error)
	GetStaffProfileService(ctx context.Context, userID uint) (*dto.UserResponse, error)
	GetAllUsersService(ctx context.Context, page, pageSize int) ([]*dto.UserResponse, *utils.PaginatedMeta, error)
	DeactivateUserService(ctx context.Context, userID uint) error
	DeactivateStaffService(ctx context.Context, userID uint) error
	ActivateUserService(ctx context.Context, userID uint) error
	ActivateStaffService(ctx context.Context, userID uint) error
	GetAllStaffService(ctx context.Context, page int, pageSize int) ([]*dto.UserResponse, *utils.PaginatedMeta, error)
	UpdateProfileService(ctx context.Context, userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
	UpdateStaffService(ctx context.Context, userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
	updateUser(ctx context.Context, user *models.User, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
	buildUserListResponse(users []models.User, total int64, page, pageSize int) ([]*dto.UserResponse, *utils.PaginatedMeta, error)
	SearchUser(ctx context.Context, req *dto.UserSearchRequest) ([]*dto.UserSearchResponse, *utils.PaginatedMeta, error)
}

// Waitlist service
type WaitlistServiceInterface interface {
	GetWaitlistPosition(ctx context.Context, userID uint, date, timeSlot string) (int, error)
	JoinWaitlist(ctx context.Context, userID uint, req *dto.JoinWaitlistRequest) (*dto.WaitlistResponse, error)
}
