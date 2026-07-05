package services

import (
	"context"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/mapper"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
)

type CartService struct {
	cartRepo repositories.CartRepositoryInterface
	menuRepo repositories.MenuRepositoryInterface
}

func NewCartService(
	cartRepo repositories.CartRepositoryInterface,
	menuRepo repositories.MenuRepositoryInterface,
) *CartService {
	return &CartService{cartRepo: cartRepo, menuRepo: menuRepo}
}

func (s *CartService) GetUserCart(ctx context.Context, userID uint) (*dto.CartResponse, error) {
	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return mapper.ConvertToCartResponse(cart), nil
}

func (s *CartService) AddItemToCart(ctx context.Context, userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) {
	// verify menu item exists
	_, err := s.menuRepo.GetByID(ctx, req.MenuItemID)
	if err != nil {
		return nil, domain.ErrMenuNotFound
	}

	// get cart — must already exist (created on registration)
	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, domain.ErrCartNotFound
	}

	// check if item already in cart
	existingItem, err := s.cartRepo.GetCartItemByMenuAndCartID(ctx, req.MenuItemID, cart.ID)
	if err == nil {
		// item exists → increment quantity
		existingItem.Quantity += req.Quantity
		if err := s.cartRepo.UpdateCartItem(ctx, existingItem); err != nil {
			return nil, err
		}
	} else {
		// item doesn't exist → add new
		newItem := &models.CartItem{
			CartID:              cart.ID,
			MenuID:              req.MenuItemID,
			Quantity:            req.Quantity,
			SpecialInstructions: req.SpecialInstructions,
		}
		if err := s.cartRepo.AddCartItem(ctx, newItem); err != nil {
			return nil, err
		}
	}

	return s.GetUserCart(ctx, userID)
}

func (s *CartService) UpdateCartItem(ctx context.Context, userID uint, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	cartItem, err := s.cartRepo.GetCartItemByIDAndUser(ctx, itemID, userID)
	if err != nil {
		return nil, domain.ErrCartItemNotFound
	}

	_, err = s.menuRepo.GetByID(ctx, cartItem.MenuID)
	if err != nil {
		return nil, domain.ErrMenuNotFound
	}

	cartItem.Quantity = req.Quantity
	cartItem.SpecialInstructions = req.SpecialInstructions

	if err := s.cartRepo.UpdateCartItem(ctx, cartItem); err != nil {
		return nil, err
	}

	return s.GetUserCart(ctx, userID)
}

func (s *CartService) RemoveCartItem(ctx context.Context, userID uint, itemID uint) error {
	return s.cartRepo.DeleteCartItem(ctx, itemID, userID)
}

func (s *CartService) ClearCart(ctx context.Context, userID uint) error {
	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return domain.ErrCartNotFound
	}
	return s.cartRepo.ClearCart(ctx, cart.ID)
}
