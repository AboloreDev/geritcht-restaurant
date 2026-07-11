package repositories

import (
	"context"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User Repo
type UserRepositoryInterface interface {
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	GetAll(ctx context.Context) ([]*models.User, error)
	GetByIdAndActive(ctx context.Context, id uint, active bool) (*models.User, error)

	GetAllByRole(ctx context.Context, role models.UserRole, page, pageSize int) ([]models.User, int64, error)
	GetByIDAndRole(ctx context.Context, id uint, role models.UserRole) (*models.User, error)
	UpdateActiveByRole(ctx context.Context, id uint, role models.UserRole, active bool) error

	TsvectorSearchUsers(ctx context.Context, req *dto.UserSearchRequest) ([]models.UserWithRank, int64, error)
}

// Auth Repository
type AuthRepositoryInterface interface {
	// refresh tokens
	CreateRefreshToken(userID uint, refreshToken string) error
	GetValidRefreshToken(ctx context.Context, refreshToken string) (*models.RefreshToken, error)
	GetRefreshToken(tokenHash string) (*models.RefreshToken, error)
	DeleteRefreshToken(token *models.RefreshToken) error
	DeleteExpiredRefreshTokens(userID uint) error

	// email tokens
	CreateEmailToken(token *models.Token) error
	GetValidEmailToken(tokenHash string) (*models.Token, error)
	VerifyUserEmail(user *models.User, token *models.Token) error

	// password
	ResetPassword(user *models.User, token *models.Token, hashedPassword string) error
	ChangePassword(user *models.User, hashedPassword string) error
}

// Cart Repository
type CartRepositoryInterface interface {
	GetCartByUserID(ctx context.Context, userID uint) (*models.Cart, error)
	GetCartByID(ctx context.Context, cartID uint) (*models.Cart, error)
	CreateCart(ctx context.Context, cart *models.Cart) error
	DeleteCart(ctx context.Context, cartID uint) error
	ClearCart(ctx context.Context, cartID uint) error

	GetCartItemsByCartID(ctx context.Context, cartID uint) ([]*models.CartItem, error)
	GetCartItemByID(ctx context.Context, id uint) (*models.CartItem, error)
	GetCartItemByMenuAndCartID(ctx context.Context, menuID, cartID uint) (*models.CartItem, error)
	GetCartItemByIDAndUser(ctx context.Context, itemID, userID uint) (*models.CartItem, error)
	AddCartItem(ctx context.Context, cartItem *models.CartItem) error
	UpdateCartItem(ctx context.Context, cartItem *models.CartItem) error
	DeleteCartItem(ctx context.Context, id uint, userID uint) error
	GetCartByUserIDForTx(ctx context.Context, tx *gorm.DB, userID uint) (*models.Cart, error)
}

// Category Repository
type CategoryRepositoryInterface interface {
	Create(ctx context.Context, category *models.MenuCategory) error
	GetByID(ctx context.Context, categoryID uint) (*models.MenuCategory, error)
	GetByName(ctx context.Context, name string) (*models.MenuCategory, error)
	GetAll(ctx context.Context, page, pageSize int) ([]models.MenuCategory, int64, error)
	Update(ctx context.Context, category *models.MenuCategory) error
	Delete(ctx context.Context, categoryID uint) error
	CountMenuItems(ctx context.Context, categoryID uint) (int64, error)

	// Search tsvector
	TsvectorSearchCategories(ctx context.Context, req *dto.CategorySearchRequest) ([]models.MenuCategoryWithRank, int64, error)
}

// Dietary tags interface
type DietaryTagRepositoryInterface interface {
	Create(ctx context.Context, tag *models.DietaryTag) error
	GetByName(ctx context.Context, name string) (*models.DietaryTag, error)
	GetByID(ctx context.Context, tagID uint) (*models.DietaryTag, error)
	GetAll(ctx context.Context, page, pageSize int) ([]models.DietaryTag, int64, error)
	Update(ctx context.Context, tag *models.DietaryTag) error
	Delete(ctx context.Context, tagID uint) error
	CountMenuItemsUsingTag(ctx context.Context, tagID uint) (int64, error)
}

// Allergen repository
type AllergenRepositoryInterface interface {
	Create(ctx context.Context, allergen *models.Allergen) error
	GetByName(ctx context.Context, name string) (*models.Allergen, error)
	GetByID(ctx context.Context, allergenID uint) (*models.Allergen, error)
	GetAll(ctx context.Context, page, pageSize int) ([]models.Allergen, int64, error)
	Update(ctx context.Context, allergen *models.Allergen) error
	Delete(ctx context.Context, allergenID uint) error
	CountMenuItemsUsingAllergen(ctx context.Context, allergenID uint) (int64, error)
}

// Menu Repository Interface
type MenuRepositoryInterface interface {
	GetCategoryByID(ctx context.Context, categoryID uint) (*models.MenuCategory, error)
	GetAllergensByIDs(ctx context.Context, ids []uint) ([]models.Allergen, error)
	GetDietaryTagsByIDs(ctx context.Context, ids []uint) ([]models.DietaryTag, error)
	CountByNameAndCategory(ctx context.Context, name string, categoryID uint) (int64, error)
	Create(ctx context.Context, menu *models.Menu) error
	GetByID(ctx context.Context, menuID uint) (*models.Menu, error)
	GetByIDAvailable(ctx context.Context, menuID uint) (*models.Menu, error)
	Update(ctx context.Context, menu *models.Menu) error
	ReplaceAllergens(ctx context.Context, menu *models.Menu, allergens []models.Allergen) error
	ReplaceDietaryTags(ctx context.Context, menu *models.Menu, tags []models.DietaryTag) error
	Delete(ctx context.Context, menuID uint) error
	GetAll(ctx context.Context, filter dto.MenuFilterRequest) ([]models.Menu, int64, error)

	// images
	CountImages(ctx context.Context, menuID uint) (int64, error)
	CreateImage(ctx context.Context, image *models.MenuImage) error
	GetImageByID(ctx context.Context, imageID uint) (*models.MenuImage, error)
	DeleteImage(ctx context.Context, image *models.MenuImage) error
	GetNextPrimaryImage(ctx context.Context, menuID uint, excludeID uint) (*models.MenuImage, error)
	SetImagePrimary(ctx context.Context, image *models.MenuImage) error

	// Search tsvector
	TsvectorSearchMenuItems(ctx context.Context, req *dto.MenuSearchRequest) ([]models.MenuWithRank, int64, error)
}

// Order repository interface
type OrderRepositoryInterface interface {
	Create(ctx context.Context, tx *gorm.DB, order *models.Order) error
	GetByID(ctx context.Context, tx *gorm.DB, orderID uint) (*models.Order, error)
	GetByIDAndUser(ctx context.Context, orderID, userID uint) (*models.Order, error)
	GetAllByUser(ctx context.Context, userID uint, page, pageSize int) ([]models.Order, int64, error)
	GetAll(ctx context.Context, page, pageSize int) ([]models.Order, int64, error)
	UpdateStatus(ctx context.Context, orderID uint, status models.OrderStatus) error
	CountByUserAndID(ctx context.Context, orderID, userID uint) (int64, error)
}

// Payment repository interface
type PaymentRepositoryInterface interface {
	// Order lookups
	GetOrderByIDAndUser(ctx context.Context, tx *gorm.DB, orderID, userID uint) (*models.Order, error)
	GetOrderByID(ctx context.Context, tx *gorm.DB, orderID uint) (*models.Order, error)
	UpdateOrderStatus(ctx context.Context, tx *gorm.DB, orderID uint, updates map[string]interface{}) error

	// Payment operations
	GetPaymentByOrderID(ctx context.Context, orderID uint) (*models.Payment, error)
	GetPaymentByReference(ctx context.Context, reference string) (*models.Payment, error)
	GetPaymentByID(ctx context.Context, paymentID uint) (*models.Payment, error)
	UpdatePayment(ctx context.Context, tx *gorm.DB, payment *models.Payment, updates map[string]interface{}) error
	GetAllByUserID(ctx context.Context, userID uint, page, pageSize int) ([]models.Payment, int64, error)
	Create(ctx context.Context, tx *gorm.DB, payment *models.Payment) error
	RecheckPaymentWithReference(ctx context.Context, reference string) error

	// Cart
	ClearCartByUserID(ctx context.Context, tx *gorm.DB, userID uint) error

	// Refund
	GetRefundByOrderID(ctx context.Context, orderID uint) (*models.Refund, error)
	GetRefundByID(ctx context.Context, refundID uint) (*models.Refund, error)
	CreateRefund(ctx context.Context, tx *gorm.DB, refund *models.Refund) error

	// Outbox
	CreateOutboxEvent(ctx context.Context, tx *gorm.DB, event *models.OutboxEvent) error
	MarkOutboxPublished(ctx context.Context, eventType string) error
}

// Waitlist repository interface
type WaitlistRepositoryInterface interface {
	CountAvailableTables(ctx context.Context, date, timeSlot string, partySize int) (int64, error)
	GetByUserDateSlot(ctx context.Context, userID uint, date, timeSlot string) (*models.Waitlist, error)
	Create(ctx context.Context, waitlist *models.Waitlist) error
	GetPosition(ctx context.Context, date string, timeSlot datatypes.Time, createdAt time.Time) (int64, error)
}

// Table repository interface
type TableRepositoryInterface interface {
	GetByName(ctx context.Context, name string) (*models.Table, error)
	GetByID(ctx context.Context, tableID uint) (*models.Table, error)
	GetByIDWithRelations(ctx context.Context, tableID uint) (*models.Table, error)
	Create(ctx context.Context, table *models.Table) error
	Update(ctx context.Context, table *models.Table) error
	Delete(ctx context.Context, tableID uint) error
	GetAll(ctx context.Context, page, pageSize int) ([]models.Table, int64, error)
}

// Reservation repository interface
type ReservationRepositoryInterface interface {
	// Table queries
	GetTableByIDAndCapacity(ctx context.Context, tableID uint, partySize int) (*models.Table, error)
	GetTablesByCapacity(ctx context.Context, partySize int) ([]models.Table, error)
	UpdateTableStatus(ctx context.Context, tx *gorm.DB, tableID uint, status models.TableStatus) error

	// Reservation queries
	GetReservationsByDateAndSlot(ctx context.Context, date string, timeSlot datatypes.Time) ([]models.Reservation, error)
	CountByTableDateSlot(ctx context.Context, tableID uint, date string, timeSlot datatypes.Time) (int64, error)
	Create(ctx context.Context, tx *gorm.DB, reservation *models.Reservation) error
	GetByIDAndUser(ctx context.Context, reservationID, userID uint) (*models.Reservation, error)
	GetByIDWithRelations(ctx context.Context, reservationID uint) (*models.Reservation, error)
	GetByIDAndStatus(ctx context.Context, reservationID uint, status models.ReservationStatus) (*models.Reservation, error)
	UpdateStatus(ctx context.Context, tx *gorm.DB, reservationID uint, updates map[string]interface{}) error
	GetAllByUser(ctx context.Context, userID uint, req *dto.ReservationFilterRequest) ([]models.Reservation, int64, error)
	GetAll(ctx context.Context, req *dto.ReservationFilterRequest) ([]models.Reservation, int64, error)
	GetTodayReservations(ctx context.Context, req *dto.ReservationFilterRequest) ([]models.Reservation, int64, error)

	// Waitlist
	GetFirstWaitlistByDateSlot(ctx context.Context, tx *gorm.DB, date interface{}, timeSlot datatypes.Time, partySize int) (*models.Waitlist, error)
	UpdateWaitlistStatus(ctx context.Context, tx *gorm.DB, waitlist *models.Waitlist, updates map[string]interface{}) error
	LockTableForUpdate(ctx context.Context, tx *gorm.DB, tableID uint) (*models.Table, error)
}

// Inventory repository interface
type InventoryRepositoryInterface interface {
	GetRecipesByMenuItemID(ctx context.Context, tx *gorm.DB, menuItemID uint) ([]models.MenuItemIngredient, error)
	GetIngredientByID(ctx context.Context, tx *gorm.DB, ingredientID uint) (*models.Ingredient, error)
	DeductIngredientStock(ctx context.Context, tx *gorm.DB, ingredientID uint, required float64) (int64, error)
	CreateStockMovement(ctx context.Context, tx *gorm.DB, movement *models.StockMovement) error
	GetAdminUser(ctx context.Context, tx *gorm.DB) (*models.User, error)
	GetLowStockIngredients(ctx context.Context, tx *gorm.DB) ([]models.Ingredient, error)
	GetOutOfStockIngredients(ctx context.Context, tx *gorm.DB) ([]models.Ingredient, error)
	CreateOutboxEvent(ctx context.Context, tx *gorm.DB, event *models.OutboxEvent) error
	MarkOutboxPublished(ctx context.Context, tx *gorm.DB, outboxID uint) error
	GetMenuItemIDsByIngredient(ctx context.Context, tx *gorm.DB, ingredientID uint) ([]uint, error)
	DisableMenuItems(ctx context.Context, tx *gorm.DB, menuItemIDs []uint) error
}

// Ingredient repository interface
type IngredientRepositoryInterface interface {
	GetIngredientByName(ctx context.Context, name string) (*models.Ingredient, error)
	CreateIngredient(ctx context.Context, ingredient *models.Ingredient) error
	GetIngredientByID(ctx context.Context, ingredientID uint) (*models.Ingredient, error)
	UpdateIngredient(ctx context.Context, ingredient *models.Ingredient) error
	IngredientCount(ctx context.Context, ingredientID uint) (int64, error)
	DeleteIngredient(ctx context.Context, ingredientID uint) error
	GetAllIngredients(ctx context.Context, page, pageSize int) ([]models.Ingredient, int64, error)
	CompareCurrentStockAgainstMinTheshold(ctx context.Context) ([]models.Ingredient, error)
	UpdateThreshHoldLimit(ctx context.Context, ingredientID uint, threshHold float64) error

	// Search tsvector
	TsvectorSearchIngredients(ctx context.Context, req *dto.IngredientSearchRequest) ([]models.IngredientWithRank, int64, error)
}

// Outbox repository interface
type OutboxRepositoryInterface interface {
	GetPendingEvents(ctx context.Context) ([]models.OutboxEvent, error)
	UpdateRetryCount(ctx context.Context, event *models.OutboxEvent) error
	MarkAsPublished(ctx context.Context, event *models.OutboxEvent) error
}

// Reservation NoShow interface
type ReservationNoShowInterface interface {
	GetAllReservations(ctx context.Context) ([]models.Reservation, error)

	// Business Logic
	MarkReservationNoShow(ctx context.Context, reservation *models.Reservation) error
}

// Reservation reminder interface
type ReservationReminderInterface interface {
	GetAllUpcomingReservations(ctx context.Context, now time.Time, windowStart, windowEnd string) ([]models.Reservation, error)
	UpdateReminderValue(ctx context.Context, reservation *models.Reservation) error
}

// Reservation Checkout Interface
type ReservationCheckoutInterface interface {
	GetAllRservations(ctx context.Context, now time.Time) ([]models.Reservation, error)

	Checkout(ctx context.Context, reservation models.Reservation) error
}

// Menu Item Ingredient (Recipes) Interface
type RecipesRepositoryInterface interface {
	CheckForLinkedIngredient(ctx context.Context, menuItemID uint, ingredientID uint) (*models.MenuItemIngredient, error)
	CreateRecipe(ctx context.Context, recipe *models.MenuItemIngredient) error
	GetLinkedIngredient(ctx context.Context, menuItemID uint, ingredientID uint) (*models.MenuItemIngredient, error)
	UpdateLinkedIngredients(ctx context.Context, recipe *models.MenuItemIngredient) error
	DeleteLinkedIngredient(ctx context.Context, menuItemID uint, ingredientID uint) error
	GetRecipesByMenuItemID(ctx context.Context, menuItemID uint) ([]models.MenuItemIngredient, error)
}
