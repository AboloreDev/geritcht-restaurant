package services

import (
	"context"
	"testing"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	redisStore "github.com/AboloreDev/geritcht-restaurant/internals/redis"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testInventoryCtx = context.Background()

//  MockInventoryRepository

type MockInventoryRepository struct {
	recipes        []models.MenuItemIngredient
	ingredient     *models.Ingredient
	admin          *models.User
	lowIngredients []models.Ingredient
	outOfStock     []models.Ingredient
	menuItemIDs    []uint
	rowsAffected   int64

	recipesErr    error
	ingredientErr error
	adminErr      error
	deductErr     error
	movementErr   error
	outboxErr     error
	lowStockErr   error
	outOfStockErr error
	menuIDsErr    error
	disableErr    error
}

func (m *MockInventoryRepository) GetRecipesByMenuItemID(_ context.Context, _ *gorm.DB, menuItemID uint) ([]models.MenuItemIngredient, error) {
	return m.recipes, m.recipesErr
}
func (m *MockInventoryRepository) GetIngredientByID(_ context.Context, _ *gorm.DB, ingredientID uint) (*models.Ingredient, error) {
	return m.ingredient, m.ingredientErr
}
func (m *MockInventoryRepository) DeductIngredientStock(_ context.Context, _ *gorm.DB, ingredientID uint, required float64) (int64, error) {
	return m.rowsAffected, m.deductErr
}
func (m *MockInventoryRepository) CreateStockMovement(_ context.Context, _ *gorm.DB, movement *models.StockMovement) error {
	return m.movementErr
}
func (m *MockInventoryRepository) GetAdminUser(_ context.Context, _ *gorm.DB) (*models.User, error) {
	return m.admin, m.adminErr
}
func (m *MockInventoryRepository) GetLowStockIngredients(_ context.Context, _ *gorm.DB) ([]models.Ingredient, error) {
	return m.lowIngredients, m.lowStockErr
}
func (m *MockInventoryRepository) GetOutOfStockIngredients(_ context.Context, _ *gorm.DB) ([]models.Ingredient, error) {
	return m.outOfStock, m.outOfStockErr
}
func (m *MockInventoryRepository) CreateOutboxEvent(_ context.Context, _ *gorm.DB, event *models.OutboxEvent) error {
	event.ID = 1
	return m.outboxErr
}
func (m *MockInventoryRepository) MarkOutboxPublished(_ context.Context, _ *gorm.DB, outboxID uint) error {
	return nil
}
func (m *MockInventoryRepository) GetMenuItemIDsByIngredient(_ context.Context, _ *gorm.DB, ingredientID uint) ([]uint, error) {
	return m.menuItemIDs, m.menuIDsErr
}
func (m *MockInventoryRepository) DisableMenuItems(_ context.Context, _ *gorm.DB, menuItemIDs []uint) error {
	return m.disableErr
}

func newInventoryService(repo *MockInventoryRepository) *InventoryService {
	return NewInventoryService(
		nil,
		&MockPublisher{},
		redisStore.NewNopCache(),
		repo)
}

//  DeductStock Tests

func TestDeductStock_NoRecipes_Skips(t *testing.T) {
	service := newInventoryService(&MockInventoryRepository{
		recipes:        []models.MenuItemIngredient{}, // no recipes
		admin:          &models.User{ID: 1, Email: "admin@test.com"},
		lowIngredients: []models.Ingredient{},
		outOfStock:     []models.Ingredient{},
	})

	orderItems := []models.OrderItem{
		{MenuID: 1, Quantity: 2},
	}

	err := service.DeductStock(testInventoryCtx, nil, orderItems, 1, 1)

	assert.NoError(t, err) // skips item with no recipe
}

func TestDeductStock_InsufficientStock(t *testing.T) {
	service := newInventoryService(&MockInventoryRepository{
		recipes: []models.MenuItemIngredient{
			{IngredientID: 1, Quantity: 300}, // needs 300g
		},
		ingredient: &models.Ingredient{
			ID:           1,
			CurrentStock: 100, // only 100g available
		},
	})

	orderItems := []models.OrderItem{
		{MenuID: 1, Quantity: 1},
	}

	err := service.DeductStock(testInventoryCtx, nil, orderItems, 1, 1)

	assert.Equal(t, domain.ErrInsufficientStock, err)
}

func TestDeductStock_AtomicCheckFails(t *testing.T) {
	// enough stock at check time but atomic update shows 0 rows
	// simulates race condition
	service := newInventoryService(&MockInventoryRepository{
		recipes: []models.MenuItemIngredient{
			{IngredientID: 1, Quantity: 100},
		},
		ingredient: &models.Ingredient{
			ID:           1,
			CurrentStock: 500, // passes pre-check
		},
		rowsAffected: 0, // atomic update fails → race condition
	})

	orderItems := []models.OrderItem{
		{MenuID: 1, Quantity: 1},
	}

	err := service.DeductStock(testInventoryCtx, nil, orderItems, 1, 1)

	assert.Equal(t, domain.ErrInsufficientStock, err)
}

func TestDeductStock_Success(t *testing.T) {
	service := newInventoryService(&MockInventoryRepository{
		recipes: []models.MenuItemIngredient{
			{IngredientID: 1, Quantity: 300},
		},
		ingredient: &models.Ingredient{
			ID:           1,
			CurrentStock: 5000, // plenty of stock
		},
		rowsAffected:   1, // deduction succeeded
		admin:          &models.User{ID: 1, Email: "admin@test.com", FirstName: "Admin"},
		lowIngredients: []models.Ingredient{},
		outOfStock:     []models.Ingredient{},
	})

	orderItems := []models.OrderItem{
		{MenuID: 1, Quantity: 2}, // needs 600g
	}

	err := service.DeductStock(testInventoryCtx, nil, orderItems, 1, 1)

	assert.NoError(t, err)
}

func TestDeductStock_MultipleItems(t *testing.T) {
	service := newInventoryService(&MockInventoryRepository{
		recipes: []models.MenuItemIngredient{
			{IngredientID: 1, Quantity: 300},
		},
		ingredient: &models.Ingredient{
			ID:           1,
			CurrentStock: 5000,
		},
		rowsAffected:   1,
		admin:          &models.User{ID: 1, Email: "admin@test.com"},
		lowIngredients: []models.Ingredient{},
		outOfStock:     []models.Ingredient{},
	})

	orderItems := []models.OrderItem{
		{MenuID: 1, Quantity: 2},
		{MenuID: 2, Quantity: 1},
	}

	err := service.DeductStock(testInventoryCtx, nil, orderItems, 1, 1)

	assert.NoError(t, err)
}

//  CheckAndAlertThreshold Tests

func TestCheckAndAlertThreshold_NoLowStock(t *testing.T) {
	service := newInventoryService(&MockInventoryRepository{
		admin:          &models.User{ID: 1, Email: "admin@test.com"},
		lowIngredients: []models.Ingredient{}, // none low
		outOfStock:     []models.Ingredient{},
	})

	err := service.CheckAndAlertThreshold(testInventoryCtx, nil)

	assert.NoError(t, err)
}

func TestCheckAndAlertThreshold_LowStockTriggersAlert(t *testing.T) {
	service := newInventoryService(&MockInventoryRepository{
		admin: &models.User{ID: 1, Email: "admin@test.com", FirstName: "Admin"},
		lowIngredients: []models.Ingredient{
			{ID: 1, Name: "Rice", CurrentStock: 50, MinThreshold: 100},
		},
		outOfStock: []models.Ingredient{},
	})

	err := service.CheckAndAlertThreshold(testInventoryCtx, nil)

	assert.NoError(t, err)
}

func TestCheckAndAlertThreshold_OutOfStockDisablesMenu(t *testing.T) {
	service := newInventoryService(&MockInventoryRepository{
		admin:          &models.User{ID: 1, Email: "admin@test.com"},
		lowIngredients: []models.Ingredient{},
		outOfStock: []models.Ingredient{
			{ID: 1, Name: "Rice", CurrentStock: 0},
		},
		menuItemIDs: []uint{3, 5}, // menu items using this ingredient
	})

	err := service.CheckAndAlertThreshold(testInventoryCtx, nil)

	assert.NoError(t, err)
}

func TestCheckAndAlertThreshold_AdminNotFound(t *testing.T) {
	service := newInventoryService(&MockInventoryRepository{
		adminErr: domain.ErrUserNotFound,
	})

	err := service.CheckAndAlertThreshold(testInventoryCtx, nil)

	assert.Equal(t, domain.ErrUserNotFound, err)
}
