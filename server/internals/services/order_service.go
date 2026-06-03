package services

import (
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/mapper"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) GetOrderResponse(tx *gorm.DB, orderID uint) (*dto.OrderResponse, error) {
	var order models.Order

	err := tx.Preload("OrderItems.Menu.MenuCategory").
		Preload("User").Where("id = ? ", orderID).
		First(&order).Error
	if err != nil {
		return nil, err
	}

	return mapper.OrderResponse(&order), nil
}

func (s *OrderService) CreateOrder(userID uint, notes string) (*dto.OrderResponse, error) {
	var orderResponse *dto.OrderResponse

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var cart models.Cart
		err := tx.Preload("CartItems.Menu").Where("user_id = ?", userID).First(&cart).Error
		if err != nil {
			return domain.ErrCartNotFound
		}

		if len(cart.CartItems) == 0 {
			return domain.ErrCartIsEmpty
		}

		var orderItems []models.OrderItem
		var totalAmount float64

		for i := range cart.CartItems {
			cartItem := &cart.CartItems[i]

			if !cartItem.Menu.IsAvailable {
				return domain.ErrMenuNotAvailable
			}

			itemTotal := cartItem.Menu.Price * float64(cartItem.Quantity)
			totalAmount = totalAmount + itemTotal

			orderItems = append(orderItems, models.OrderItem{
				MenuID:              cartItem.MenuID,
				Quantity:            cartItem.Quantity,
				Price:               cartItem.Menu.Price,
				SpecialInstructions: cartItem.SpecialInstructions,
			})
		}

		order := models.Order{
			UserID:        &userID,
			TotalAmount:   totalAmount,
			OrderItems:    orderItems,
			Type:          models.OrderTypeTakeout,
			Status:        models.OrderStatusPending,
			Notes:         notes,
			CreatedAt:     time.Now(),
			PaymentStatus: models.PaymentStatusUnpaid,
		}

		err = tx.Create(&order).Error
		if err != nil {
			return err
		}

		payment := models.Payment{
			OrderID:        order.ID,
			UserID:         userID,
			Reference:      uuid.New().String(),
			IdempotencyKey: uuid.New().String(),
			Amount:         totalAmount,
			Currency:       "NGN",
			Status:         models.PaymentStatusUnpaid,
			Provider:       "paystack",
		}

		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		response, err := s.GetOrderResponse(tx, order.ID)
		if err != nil {
			return err
		}

		orderResponse = response

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Publish to queue

	return orderResponse, nil
}
