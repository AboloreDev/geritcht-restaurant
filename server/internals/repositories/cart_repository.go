package repositories

import (
	"context"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetCartByUserID(ctx context.Context, userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).Preload("CartItems.Menu").
		Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, domain.ErrCartNotFound
	}
	return &cart, nil
}

func (r *CartRepository) GetCartByID(ctx context.Context, cartID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).First(&cart, cartID).Error
	if err != nil {
		return nil, domain.ErrCartNotFound
	}
	return &cart, nil
}

func (r *CartRepository) CreateCart(ctx context.Context, cart *models.Cart) error {
	return r.db.WithContext(ctx).Create(cart).Error
}

func (r *CartRepository) DeleteCart(ctx context.Context, cartID uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Cart{}, cartID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrCartNotFound
	}
	return nil
}

func (r *CartRepository) ClearCart(ctx context.Context, cartID uint) error {
	return r.db.WithContext(ctx).
		Where("cart_id = ?", cartID).
		Delete(&models.CartItem{}).Error
}

func (r *CartRepository) GetCartItemsByCartID(ctx context.Context, cartID uint) ([]*models.CartItem, error) {
	var items []*models.CartItem
	err := r.db.WithContext(ctx).Preload("Menu").
		Where("cart_id = ?", cartID).Find(&items).Error
	return items, err
}

func (r *CartRepository) GetCartItemByID(ctx context.Context, id uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err != nil {
		return nil, domain.ErrCartItemNotFound
	}
	return &item, nil
}

func (r *CartRepository) GetCartItemByMenuAndCartID(ctx context.Context, menuID, cartID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND menu_id = ?", cartID, menuID).
		First(&item).Error
	if err != nil {
		return nil, domain.ErrCartItemNotFound
	}
	return &item, nil
}

func (r *CartRepository) GetCartItemByIDAndUser(ctx context.Context, itemID, userID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.WithContext(ctx).
		Joins("JOIN carts ON cart_items.cart_id = carts.id").
		Where("cart_items.id = ? AND carts.user_id = ?", itemID, userID).
		First(&item).Error
	if err != nil {
		return nil, domain.ErrCartItemNotFound
	}
	return &item, nil
}

func (r *CartRepository) AddCartItem(ctx context.Context, cartItem *models.CartItem) error {
	return r.db.WithContext(ctx).Create(cartItem).Error
}

func (r *CartRepository) UpdateCartItem(ctx context.Context, cartItem *models.CartItem) error {
	return r.db.WithContext(ctx).Save(cartItem).Error
}

func (r *CartRepository) DeleteCartItem(ctx context.Context, id uint, userID uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND cart_id IN (?)", id,
			r.db.Select("id").Table("carts").Where("user_id = ?", userID)).
		Delete(&models.CartItem{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrCartItemNotFound
	}
	return nil
}

func (r *CartRepository) GetCartByUserIDForTx(ctx context.Context, tx *gorm.DB, userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := tx.WithContext(ctx).Preload("CartItems.Menu").
		Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, domain.ErrCartNotFound
	}
	return &cart, nil
}
