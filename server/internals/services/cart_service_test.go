package services

import (
	"context"
	"testing"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testCartCtx = context.Background()

// ─── MockCartRepository

type MockCartRepository struct {
	cart       *models.Cart
	item       *models.CartItem
	getCartErr error
	getItemErr error
	createErr  error
	updateErr  error
	deleteErr  error
	clearErr   error
}

func (m *MockCartRepository) GetCartByUserID(_ context.Context, userID uint) (*models.Cart, error) {
	return m.cart, m.getCartErr
}
func (m *MockCartRepository) GetCartByID(_ context.Context, cartID uint) (*models.Cart, error) {
	return m.cart, m.getCartErr
}
func (m *MockCartRepository) CreateCart(_ context.Context, cart *models.Cart) error {
	return m.createErr
}
func (m *MockCartRepository) DeleteCart(_ context.Context, cartID uint) error { return m.deleteErr }
func (m *MockCartRepository) ClearCart(_ context.Context, cartID uint) error  { return m.clearErr }
func (m *MockCartRepository) GetCartItemsByCartID(_ context.Context, cartID uint) ([]*models.CartItem, error) {
	return nil, nil
}
func (m *MockCartRepository) GetCartItemByID(_ context.Context, id uint) (*models.CartItem, error) {
	return m.item, m.getItemErr
}
func (m *MockCartRepository) GetCartItemByMenuAndCartID(_ context.Context, menuID, cartID uint) (*models.CartItem, error) {
	return m.item, m.getItemErr
}
func (m *MockCartRepository) GetCartItemByIDAndUser(_ context.Context, itemID, userID uint) (*models.CartItem, error) {
	return m.item, m.getItemErr
}
func (m *MockCartRepository) AddCartItem(_ context.Context, item *models.CartItem) error {
	return m.createErr
}
func (m *MockCartRepository) UpdateCartItem(_ context.Context, item *models.CartItem) error {
	return m.updateErr
}
func (m *MockCartRepository) DeleteCartItem(_ context.Context, id uint, userID uint) error {
	return m.deleteErr
}
func (m *MockCartRepository) GetCartByUserIDForTx(_ context.Context, _ *gorm.DB, userID uint) (*models.Cart, error) {
	return m.cart, m.getCartErr
}

// ─── MockMenuRepository (reuse from menu tests, just need GetByID)

func newCartService(cartRepo *MockCartRepository, menuRepo *MockMenuRepository) *CartService {
	return NewCartService(cartRepo, menuRepo)
}

// ─── AddItemToCart Tests

func TestAddItemToCart_Success_NewItem(t *testing.T) {
	service := newCartService(
		&MockCartRepository{
			cart:       &models.Cart{ID: 1, UserID: 1},
			getItemErr: domain.ErrCartItemNotFound, // not in cart yet
		},
		&MockMenuRepository{menu: &models.Menu{ID: 1}},
	)

	req := &dto.AddToCartRequest{MenuItemID: 1, Quantity: 2}
	response, err := service.AddItemToCart(testCartCtx, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestAddItemToCart_MenuNotFound(t *testing.T) {
	service := newCartService(
		&MockCartRepository{},
		&MockMenuRepository{getErr: domain.ErrMenuNotFound},
	)

	req := &dto.AddToCartRequest{MenuItemID: 999, Quantity: 1}
	response, err := service.AddItemToCart(testCartCtx, 1, req)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrMenuNotFound, err)
}

func TestAddItemToCart_CartNotFound(t *testing.T) {
	service := newCartService(
		&MockCartRepository{getCartErr: domain.ErrCartNotFound},
		&MockMenuRepository{menu: &models.Menu{ID: 1}},
	)

	req := &dto.AddToCartRequest{MenuItemID: 1, Quantity: 1}
	response, err := service.AddItemToCart(testCartCtx, 999, req)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrCartNotFound, err)
}

func TestAddItemToCart_ExistingItem_IncrementsQuantity(t *testing.T) {
	service := newCartService(
		&MockCartRepository{
			cart: &models.Cart{ID: 1, UserID: 1},
			item: &models.CartItem{ID: 1, CartID: 1, MenuID: 1, Quantity: 2},
		},
		&MockMenuRepository{menu: &models.Menu{ID: 1}},
	)

	req := &dto.AddToCartRequest{MenuItemID: 1, Quantity: 3}
	response, err := service.AddItemToCart(testCartCtx, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
}

// ─── RemoveCartItem Tests

func TestRemoveCartItem_Success(t *testing.T) {
	service := newCartService(&MockCartRepository{}, &MockMenuRepository{})

	err := service.RemoveCartItem(testCartCtx, 1, 1)

	assert.NoError(t, err)
}

func TestRemoveCartItem_NotFound(t *testing.T) {
	service := newCartService(
		&MockCartRepository{deleteErr: domain.ErrCartItemNotFound},
		&MockMenuRepository{},
	)

	err := service.RemoveCartItem(testCartCtx, 1, 999)

	assert.Equal(t, domain.ErrCartItemNotFound, err)
}

// ─── ClearCart Tests

func TestClearCart_Success(t *testing.T) {
	service := newCartService(
		&MockCartRepository{cart: &models.Cart{ID: 1, UserID: 1}},
		&MockMenuRepository{},
	)

	err := service.ClearCart(testCartCtx, 1)

	assert.NoError(t, err)
}

func TestClearCart_CartNotFound(t *testing.T) {
	service := newCartService(
		&MockCartRepository{getCartErr: domain.ErrCartNotFound},
		&MockMenuRepository{},
	)

	err := service.ClearCart(testCartCtx, 999)

	assert.Equal(t, domain.ErrCartNotFound, err)
}
