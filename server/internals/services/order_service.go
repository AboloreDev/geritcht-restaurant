package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
			CreatedAt:     time.Now(),
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

	s.redisStore.FlushByPattern(ctx, "orders:user:*")
	s.redisStore.FlushByPattern(ctx, "orders:all:*")

	return orderResponse, nil
}

func (s *OrderService) GetAllUserTakeoutOrders(ctx context.Context, userID uint, filter *dto.OrderFilterRequest) (*dto.OrderListResponse, error) {
	cacheKey := utils.BuildUserOrderCacheKey(userID, filter)
	cached, err := s.redisStore.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cachedResponse struct {
			Data *dto.OrderListResponse `json:"data"`
		}
		if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
			return cachedResponse.Data, nil
		}
	}

	orders, total, err := s.orderRepo.GetAllByUser(ctx, userID, filter)
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}

	response := s.buildOrderListResponse(orders, total, filter)

	cacheData := struct {
		Data *dto.OrderListResponse `json:"data"`
	}{Data: response}
	data, _ := json.Marshal(&cacheData)
	s.redisStore.Set(ctx, cacheKey, string(data), utils.GetOrderCacheTTL(filter))

	return response, nil
}

func (s *OrderService) GetTakeoutOrder(ctx context.Context, userID, orderID uint) (*dto.OrderResponse, error) {
	cachedKey := fmt.Sprintf("orders:user:%d:order:%d", userID, orderID)

	exists, _ := s.redisStore.Exists(ctx, cachedKey)
	if exists {
		cache, err := s.redisStore.Get(ctx, cachedKey)
		if err == nil && cache != "" {
			var order models.Order
			err := json.Unmarshal([]byte(cache), &order)
			if err != nil {
				return nil, err
			}
			return mapper.OrderResponse(&order), nil
		}
	}
	order, err := s.orderRepo.GetByIDAndUser(ctx, orderID, userID)
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}

	data, err := json.Marshal(&order)
	if err != nil {
		return nil, fmt.Errorf("Failed to set data: %w", err)
	}
	s.redisStore.Set(ctx, cachedKey, string(data), 20*time.Minute)

	return mapper.OrderResponse(order), nil
}

func (s *OrderService) CancelTakeoutOrder(ctx context.Context, userID, orderID uint) error {
	order, err := s.orderRepo.GetByIDAndUser(ctx, orderID, userID)
	if err != nil {
		return domain.ErrOrderNotFound
	}

	switch {
	case order.Status == models.OrderStatusCancelled:
		return domain.ErrAlreadyCancelled
	case order.Status == models.OrderStatusPreparing,
		order.Status == models.OrderStatusReady,
		order.Status == models.OrderStatusCompleted:
		return domain.ErrCannotCancel
	case order.Status == models.OrderStatusConfirmed:
		return domain.ErrRefundIsProcessing
	}

	if err := s.orderRepo.UpdateStatus(ctx, orderID, models.OrderStatusCancelled); err != nil {
		return err
	}

	s.redisStore.FlushByPattern(ctx, "orders:user:*")
	s.redisStore.FlushByPattern(ctx, "orders:all:*")

	return nil
}

func (s *OrderService) AdminCancelOrder(ctx context.Context, orderID uint) error {
	order, err := s.orderRepo.GetByID(ctx, nil, orderID)
	if err != nil {
		return domain.ErrOrderNotFound
	}

	switch {
	case order.Status == models.OrderStatusCancelled:
		return domain.ErrAlreadyCancelled
	case order.Status == models.OrderStatusCompleted:
		return domain.ErrCannotCancel
	case order.Status == models.OrderStatusConfirmed:
		return domain.ErrRefundIsProcessing
	}

	if err := s.orderRepo.UpdateStatus(ctx, orderID, models.OrderStatusCancelled); err != nil {
		return err
	}

	s.redisStore.FlushByPattern(ctx, "orders:user:*")
	s.redisStore.FlushByPattern(ctx, "orders:all:*")

	return nil
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

func (s *OrderService) GetAllOrders(ctx context.Context, filter *dto.OrderFilterRequest) (*dto.OrderListResponse, error) {
	cacheKey := utils.BuildOrderCacheKey(filter)
	cached, err := s.redisStore.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cachedResponse struct {
			Data *dto.OrderListResponse `json:"data"`
		}
		if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
			return cachedResponse.Data, nil
		}
	}

	log.Printf("filter.Page: %d, offset would be: %d", filter.Page, utils.Pagination(filter.Page, filter.PageSize))

	orders, total, err := s.orderRepo.GetAll(ctx, filter)
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}

	response := s.buildOrderListResponse(orders, total, filter)

	cacheData := struct {
		Data *dto.OrderListResponse `json:"data"`
	}{Data: response}
	data, _ := json.Marshal(&cacheData)
	s.redisStore.Set(ctx, cacheKey, string(data), utils.GetOrderCacheTTL(filter))

	return response, nil
}

func (s *OrderService) SearchOrders(ctx context.Context, req *dto.OrderSearchRequest) ([]*dto.OrderSearchResponse, *utils.PaginatedMeta, error) {
	cacheKey := fmt.Sprintf("orders:%s:p%d:s%d", req.Query, req.Page, req.Limit)
	cached, err := s.redisStore.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cachedResponse struct {
			Data []*dto.OrderSearchResponse `json:"data"`
			Meta *utils.PaginatedMeta       `json:"meta"`
		}
		if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
			return cachedResponse.Data, cachedResponse.Meta, nil
		}
	}

	rows, count, err := s.orderRepo.TsvectorSearchOrders(ctx, req)
	if err != nil {
		return nil, nil, domain.ErrIngredientSearchNotFound
	}

	response := make([]*dto.OrderSearchResponse, len(rows))

	for i := range rows {
		var userResponse *dto.UserResponse

		if rows[i].User != nil {
			userResponse = &dto.UserResponse{
				ID:        rows[i].User.ID,
				FirstName: rows[i].User.FirstName,
				LastName:  rows[i].User.LastName,
				Email:     rows[i].User.Email,
			}
		}

		response[i] = &dto.OrderSearchResponse{
			OrderResponse: dto.OrderResponse{
				ID:          rows[i].ID,
				UserID:      rows[i].UserID,
				User:        userResponse,
				Type:        string(rows[i].Type),
				Status:      string(rows[i].Status),
				TotalAmount: rows[i].TotalAmount,
				CreatedAt:   rows[i].CreatedAt,
			},
			Rank: rows[i].Rank,
		}
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	totalPages := int((count + int64(req.Limit) - 1) / int64(req.Limit))
	meta := &utils.PaginatedMeta{
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      count,
		TotalPages: totalPages,
	}

	cacheData := struct {
		Data []*dto.OrderSearchResponse `json:"data"`
		Meta *utils.PaginatedMeta       `json:"meta"`
	}{Data: response, Meta: meta}

	data, _ := json.Marshal(&cacheData)
	s.redisStore.Set(ctx, cacheKey, string(data), 5*time.Minute)

	return response, meta, nil
}

func (s *OrderService) buildOrderListResponse(orders []models.Order, count int64, req *dto.OrderFilterRequest) *dto.OrderListResponse {
	response := make([]dto.OrderResponse, 0, len(orders))

	for _, ord := range orders {
		response = append(response, *mapper.OrderResponse(&ord))
	}

	totalPages := int((count + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &dto.OrderListResponse{
		Orders:     response,
		Total:      count,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}
}
