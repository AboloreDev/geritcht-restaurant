package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/mapper"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	db          *gorm.DB // still needed for transaction control
	orderRepo   repositories.OrderRepositoryInterface
	paymentRepo repositories.PaymentRepositoryInterface
	cartRepo    repositories.CartRepositoryInterface
	redisStore  interfaces.Cacher
}

func NewOrderService(
	db *gorm.DB,
	orderRepo repositories.OrderRepositoryInterface,
	paymentRepo repositories.PaymentRepositoryInterface,
	cartRepo repositories.CartRepositoryInterface,
	redisStore interfaces.Cacher,
) *OrderService {
	return &OrderService{
		db:          db,
		orderRepo:   orderRepo,
		paymentRepo: paymentRepo,
		cartRepo:    cartRepo,
		redisStore:  redisStore,
	}
}

func (s *OrderService) GetOrderResponse(ctx context.Context, tx *gorm.DB, orderID uint) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, tx, orderID)
	if err != nil {
		return nil, err
	}
	return mapper.OrderResponse(order), nil
}

func (s *OrderService) CreateTakeoutOrder(ctx context.Context, userID uint, req *dto.CreateTakeoutOrderRequest) (*dto.OrderResponse, error) {
	var orderResponse *dto.OrderResponse

	err := s.db.Transaction(func(tx *gorm.DB) error {
		cart, err := s.cartRepo.GetCartByUserIDForTx(ctx, tx, userID)
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
			totalAmount += itemTotal

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
			PaymentStatus: models.PaymentStatusUnpaid,
		}

		if err := s.orderRepo.Create(ctx, tx, &order); err != nil {
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

		if err := s.paymentRepo.Create(ctx, tx, &payment); err != nil {
			return err
		}

		response, err := s.GetOrderResponse(ctx, tx, order.ID)
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

func (s *OrderService) GetAllUserTakeoutOrders(ctx context.Context, userID uint, page, pageSize int) ([]*dto.OrderResponse, *utils.PaginatedMeta, error) {
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

	orders, total, err := s.orderRepo.GetAllByUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, nil, domain.ErrOrderNotFound
	}

	response := make([]*dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		response = append(response, mapper.OrderResponse(&order))
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &utils.PaginatedMeta{Page: page, Limit: pageSize, Total: total, TotalPages: totalPages}

	cacheData := struct {
		Data []*dto.OrderResponse `json:"data"`
		Meta *utils.PaginatedMeta `json:"meta"`
	}{Data: response, Meta: meta}
	data, _ := json.Marshal(&cacheData)
	s.redisStore.Set(ctx, cacheKey, string(data), 1*time.Hour)

	return response, meta, nil
}

func (s *OrderService) GetTakeoutOrder(ctx context.Context, userID, orderID uint) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.GetByIDAndUser(ctx, orderID, userID)
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}
	return mapper.OrderResponse(order), nil
}

func (s *OrderService) CancelTakeoutOrder(ctx context.Context, userID, orderID uint) error {
	order, err := s.orderRepo.GetByIDAndUser(ctx, orderID, userID)
	if err != nil {
		return domain.ErrOrderNotFound
	}

	if order.UserID != nil && *order.UserID != userID {
		return domain.ErrForbidden
	}

	switch {
	case order.Status == models.OrderStatusCancelled:
		return domain.ErrAlreadyCancelled
	case order.Status == models.OrderStatusPreparing ||
		order.Status == models.OrderStatusReady ||
		order.Status == models.OrderStatusCompleted:
		return domain.ErrCannotCancel
	case order.Status == models.OrderStatusConfirmed:
		return domain.ErrRefundIsProcessing
	default:
		return s.orderRepo.UpdateStatus(ctx, orderID, models.OrderStatusCancelled)
	}
}

func (s *OrderService) VerifyUserOrder(ctx context.Context, userID, orderID uint) error {
	count, err := s.orderRepo.CountByUserAndID(ctx, orderID, userID)
	if err != nil {
		return domain.ErrOrderNotFound
	}
	if count == 0 {
		return domain.ErrOrderNotFound
	}
	return nil
}

func (s *OrderService) GetAllOrders(ctx context.Context, page, pageSize int) ([]*dto.OrderResponse, *utils.PaginatedMeta, error) {
	cacheKey := fmt.Sprintf("staff:orders:p:%d:s:%d", page, pageSize)

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

	orders, total, err := s.orderRepo.GetAll(ctx, page, pageSize)
	if err != nil {
		return nil, nil, domain.ErrOrderNotFound
	}

	response := make([]*dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		response = append(response, mapper.OrderResponse(&order))
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &utils.PaginatedMeta{Page: page, Limit: pageSize, Total: total, TotalPages: totalPages}

	cacheData := struct {
		Data []*dto.OrderResponse `json:"data"`
		Meta *utils.PaginatedMeta `json:"meta"`
	}{Data: response, Meta: meta}
	data, _ := json.Marshal(&cacheData)
	s.redisStore.Set(ctx, cacheKey, string(data), 30*time.Second)

	return response, meta, nil
}
