package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/mapper"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	db         *gorm.DB
	redisStore interfaces.Cacher
}

func NewOrderService(db *gorm.DB, redisStore interfaces.Cacher) *OrderService {
	return &OrderService{
		db:         db,
		redisStore: redisStore,
	}
}

func (s *OrderService) GetOrderResponse(tx *gorm.DB, orderID uint) (*dto.OrderResponse, error) {
	var order models.Order

	err := tx.Preload("OrderItems.Menu.MenuCategory").
		Preload("User").Preload("Payment").Where("id = ? ", orderID).
		First(&order).Error
	if err != nil {
		return nil, err
	}

	return mapper.OrderResponse(&order), nil
}

func (s *OrderService) CreateTakeoutOrder(userID uint, req *dto.CreateTakeoutOrderRequest) (*dto.OrderResponse, error) {
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
			Notes:         req.Notes,
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

	return orderResponse, nil
}

func (s *OrderService) GetAllTakeoutOrders(userID uint, page, pageSize int) ([]*dto.OrderResponse, *utils.PaginatedMeta, error) {
	cacheKey := fmt.Sprintf("user:orders:%d:p:%d:s:%d", userID, page, pageSize)

	cached, err := s.redisStore.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cachedResponse struct {
			Data []*dto.OrderResponse `json:"data"`
			Meta *utils.PaginatedMeta `json:"meta"`
		}
		if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
			return cachedResponse.Data, cachedResponse.Meta, nil
		}
	}
	var orders []models.Order
	var total int64
	offset := utils.Pagination(page, pageSize)

	s.db.Model(models.Order{}).Count(&total)

	err = s.db.Preload("OrderItems.Menu").
		Preload("User").
		Preload("Payment").
		Where("user_id = ? AND type = ?", userID, models.OrderTypeTakeout).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&orders).Error
	if err != nil {
		return nil, nil, domain.ErrOrderNotFound
	}

	response := make([]*dto.OrderResponse, 0, len(orders))

	for _, order := range orders {
		response = append(response, mapper.OrderResponse(&order))
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	meta := &utils.PaginatedMeta{
		Page:       page,
		Limit:      pageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	cacheData := struct {
		Data []*dto.OrderResponse `json:"data"`
		Meta *utils.PaginatedMeta `json:"meta"`
	}{Data: response, Meta: meta}

	// store in cache
	data, err := json.Marshal(&cacheData)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to set data: %d", err)
	}
	s.redisStore.Set(ctx, cacheKey, string(data), 1*time.Hour)

	return response, meta, nil
}

func (s *OrderService) GetTakeoutOrder(userID, orderID uint) (*dto.OrderResponse, error) {
	var order models.Order

	err := s.db.Preload("OrderItems.Menu").
		Preload("User").
		Preload("Payment").
		Where("id = ? AND user_id = ? AND type = ?", orderID, userID, models.OrderTypeTakeout).
		First(&order).Error
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}

	return mapper.OrderResponse(&order), nil
}

func (s *OrderService) VerifyUserOrder(userID, orderID uint) error {
	var count int64

	err := s.db.Model(&models.Order{}).
				Where("id = ? AND user_id = ?", orderID, userID).
				Count(&count).Error
	if err != nil {
		return domain.ErrOrderNotFound
	}

	return nil
}
