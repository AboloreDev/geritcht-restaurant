package services

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/mapper"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"gorm.io/gorm"
)

type CartService struct {
	db *gorm.DB
}

func NewCartService(db *gorm.DB) *CartService {
	return &CartService{db: db}
}

func (s *CartService) GetUserCart(userID uint) (*dto.CartResponse, error) {
	var cart models.Cart
	err := s.db.Preload("CartItems.Menu").
		Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	return mapper.ConvertToCartResponse(&cart), nil
}

func (s *CartService) AddItemToCart(userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) {
	var menu models.Menu
	var cart models.Cart
	var item models.CartItem
	err := s.db.Where("id = ?", req.MenuItemID).First(&menu).Error
	if err != nil {
		return nil, err
	}

	err = s.db.Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		err = s.db.Create(&cart).Error
		if err != nil {
			return nil, err
		}
	}

	err = s.db.Where("cart_id = ? AND menu_item_id = ?", cart.ID, req.MenuItemID).First(&item).Error
	if err == nil {
		item.Quantity += req.Quantity
		err = s.db.Save(&item).Error
		if err != nil {
			return nil, err
		}
	} else {
		newItem := models.CartItem{
			CartID:              cart.ID,
			MenuID:              req.MenuItemID,
			Quantity:            req.Quantity,
			SpecialInstructions: req.SpecialInstructions,
		}
		err = s.db.Create(&newItem).Error
		if err != nil {
			return nil, err
		}
	}

	return mapper.ConvertToCartResponse(&cart), nil
}

func (s *CartService) UpdateCartItem(userID uint, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	var cartItem models.CartItem
	err := s.db.Joins("JOIN carts ON cart_items.cart_id = carts.id ").
		Where("cart_items.id = ? AND carts.user_id = ? ", itemID, userID).
		First(&cartItem).Error
	if err != nil {
		return nil, domain.ErrCartItemNotFound
	}

	var menu models.Menu
	err = s.db.First(&menu, cartItem.MenuID).Error
	if err != nil {
		return nil, domain.ErrMenuNotFound
	}

	cartItem.Quantity = req.Quantity
	cartItem.SpecialInstructions = req.SpecialInstructions
	err = s.db.Save(&cartItem).Error
	if err != nil {
		return nil, err
	}

	return s.GetUserCart(userID)
}

func (s *CartService) RemoveCartItem(userID uint, itemID uint) error {
	var cartItem models.CartItem
	result := s.db.Where("id = ? AND cart_id IN (?)", itemID,
		s.db.Select("id").Table("carts").
			Where("user_id = ?", userID)).
		Delete(&cartItem)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrCartItemNotFound
	}

	return nil
}

func (s *CartService) ClearCart(userID uint) error {
	var cart models.Cart
	err := s.db.Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return domain.ErrCartNotFound
	}

	err = s.db.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error
	if err != nil {
		return err
	}

	return nil
}
