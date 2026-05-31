package services

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/mapper"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
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

	err := tx.Preload("OrderItems.Product.Category").
				Preload("User").Where("id = ? ", orderID).
					First(&order).Error
	if err != nil {
		return nil, err
	}

	return mapper.OrderResponse(&order), nil
}

func (s *OrderService) CreateOrder(userID uint, req *dto.CreateOrderItemRequest) (*dto.OrderResponse, error) {
	var orderResponse *dto.OrderResponse

	s.db.Transaction(func(tx *gorm.DB) error {
		var cart models.Cart
	})
}