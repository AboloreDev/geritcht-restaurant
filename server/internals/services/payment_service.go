package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/events"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/mapper"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type PaystackInitialiseResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

type PaystackVerifyResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Status    string `json:"status"`
		Amount    int64  `json:"amount"`
		Reference string `json:"reference"`
		PaidAt    string `json:"paid_at"`
	} `json:"data"`
}
type PaystackRefundResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Status    string `json:"status"`
		Amount    int64  `json:"amount"`
		Reference string `json:"reference"`
		PaidAt    string `json:"paid_at"`
	} `json:"data"`
}

type PaystackWebhookPayload struct {
	Status bool   `json:"status"`
	Event  string `json:"event"`
	Data   struct {
		Reference string `json:"reference"`
		Status    string `json:"status"`
		Amount    int64  `json:"amount"`
		PaidAt    string `json:"paid_at"`
	} `json:"data"`
}

type PaymentService struct {
	db             *gorm.DB
	redisStore     interfaces.Cacher
	eventPublisher interfaces.Publisher
	config         *config.Config
	inventoryRepo  repositories.InventoryRepositoryInterface
	HTTPClient     HTTPClient
	paymentRepo    repositories.PaymentRepositoryInterface
	inventory      InventoryService
	log            zerolog.Logger
}

func NewPaymentService(
	db *gorm.DB,
	redisStore interfaces.Cacher,
	eventPublisher interfaces.Publisher,
	config *config.Config,
	inventoryRepo repositories.InventoryRepositoryInterface,
	HTTPClient HTTPClient,
	paymentRepo repositories.PaymentRepositoryInterface,
	inventory InventoryService,
	log zerolog.Logger,
) *PaymentService {
	return &PaymentService{
		db:             db,
		redisStore:     redisStore,
		eventPublisher: eventPublisher,
		config:         config,
		inventoryRepo:  inventoryRepo,
		HTTPClient:     HTTPClient,
		paymentRepo:    paymentRepo,
		inventory:      inventory,
	}
}

func (s *PaymentService) callPaystackInitialize(email string, amount int64, reference string, orderID uint) (*PaystackInitialiseResponse, error) {
	url := fmt.Sprintf("%s/transaction/initialize", s.config.Paystack.PaystackURL)
	// Creating the payload
	payload := map[string]interface{}{
		"email":  email,
		"amount": amount,
		"metadata": map[string]interface{}{
			"order_id":  orderID,
			"reference": reference,
		},
	}
	// Convert the payload from Go struct to json
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Call Paystck endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {

		s.log.Error().Err(err).Msg("operation failed")
		return nil, err
	}

	// Set the Req headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.Paystack.PaystackSecretKey)

	response, err := s.HTTPClient.Do(req)
	if err != nil {
		s.log.Error().Err(err).Msg("operation failed")
		return nil, err
	}
	defer response.Body.Close()

	// Parse response body to struct
	// Reading using NewDecoder and Encode
	var result PaystackInitialiseResponse

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		s.log.Error().Err(err).Msg("operation failed")
		return nil, err
	}

	if !result.Status {
		return nil, err
	}

	return &result, nil
}

func (s *PaymentService) callPaystackVerify(reference string) (*PaystackVerifyResponse, error) {
	url := fmt.Sprintf("%s/transaction/verify/%s", s.config.Paystack.PaystackURL, reference)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.config.Paystack.PaystackSecretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PaystackVerifyResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Status {
		return nil, fmt.Errorf("paystack verify failed: %s", result.Message)
	}

	return &result, nil
}

func (s *PaymentService) callPaystackRefund(reference string, amount int64) (*PaystackRefundResponse, error) {
	payload := map[string]interface{}{
		"reference": reference,
		"amount":    amount,
		"currency":  "NGN",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/refund", s.config.Paystack.PaystackURL)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.config.Paystack.PaystackSecretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PaystackRefundResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Status {
		return nil, fmt.Errorf("paystack refund failed: %s", result.Message)
	}

	return &result, nil
}

func (s *PaymentService) verifySignature(body []byte, signature string) bool {
	mac := hmac.New(sha512.New, []byte(s.config.Paystack.PaystackSecretKey))
	mac.Write(body)
	expectedMac := hex.EncodeToString(mac.Sum(nil))

	fmt.Println("Received :", signature)
	fmt.Println("Expected :", expectedMac)
	fmt.Println("Equal    :", hmac.Equal([]byte(expectedMac), []byte(signature)))

	return hmac.Equal([]byte(expectedMac), []byte(signature))
}

func (s *PaymentService) InitialisePayment(ctx context.Context, userID uint, req *dto.InitializePaymentRequest) (*dto.InitializePaymentResponse, error) {
	order, err := s.paymentRepo.GetOrderByIDAndUser(ctx, nil, req.OrderID, userID)
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}

	if order.Status != models.OrderStatusPending {
		return nil, domain.ErrInvalidOrderStatus
	}
	if order.PaymentStatus != models.PaymentStatusUnpaid {
		return nil, domain.ErrOrderAlreadyPaid
	}

	payment, err := s.paymentRepo.GetPaymentByOrderID(ctx, order.ID)
	if err != nil {
		return nil, domain.ErrPaymentNotFound
	}

	response, err := s.callPaystackInitialize(
		order.User.Email,
		int64(order.TotalAmount*100),
		payment.Reference,
		order.ID,
	)
	if err != nil {
		s.log.Error().Err(err).Msg("operation failed")
		return nil, fmt.Errorf("paystack initialization failed: %w", err)
	}

	if err := s.paymentRepo.UpdatePayment(ctx, nil, payment, map[string]interface{}{
		"status": models.PaymentStatusPending,
	}); err != nil {
		s.log.Error().Err(err).Msg("operation failed")

		return nil, err
	}

	return &dto.InitializePaymentResponse{
		AuthorizationURL: response.Data.AuthorizationURL,
		Reference:        response.Data.Reference,
		Payment:          *mapper.PaymentResponse(payment),
	}, nil
}

func (s *PaymentService) VerifyPayment(ctx context.Context, req *dto.VerifyPaymentRequest) (*dto.PaymentResponse, error) {
	payment, err := s.paymentRepo.GetPaymentByReference(ctx, req.Reference)
	if err != nil {
		return nil, domain.ErrPaymentNotFound
	}

	response, err := s.callPaystackVerify(req.Reference)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		ID:        payment.ID,
		OrderID:   payment.OrderID,
		Reference: payment.Reference,
		Amount:    payment.Amount,
		Status:    response.Data.Status,
		Currency:  payment.Currency,
	}, nil
}

func (s *PaymentService) HandlePaystackWebhook(ctx context.Context, body []byte, signature string) error {
	
	if !s.verifySignature(body, signature) {
		return domain.ErrInvalidSignature
	}

	var payload PaystackWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return err
	}

	if payload.Event != "charge.success" {
		return nil
	}

	reference := payload.Data.Reference

	payment, err := s.paymentRepo.GetPaymentByReference(ctx, reference)
	if err != nil {
		return domain.ErrPaymentNotFound
	}

	// idempotency check
	if payment.Status == models.PaymentStatusPaid {
		return nil
	}

	// distributed lock
	lockKey := fmt.Sprintf("lock:webhook:payment:%s", reference)
	lockValue, _ := json.Marshal(reference)
	if err := s.redisStore.Hold(ctx, lockKey, lockValue, redis.SetArgs{
		Mode: "NX",
		TTL:  30 * time.Second,
	}); err != nil {
		return nil
	}
	defer s.redisStore.Delete(ctx, lockKey)

	// Check again, to avoid double order confirmation, double cart deletion
	// if the payment is paid, return nil
	err = s.paymentRepo.RecheckPaymentWithReference(ctx, reference)
	if err != nil {
		return nil
	}

	// amount verification
	expectedAmount := int64(payment.Amount * 100)
	if payload.Data.Amount != expectedAmount {
		return domain.ErrPaymentAmountMismatch
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// update payment
		if err := s.paymentRepo.UpdatePayment(ctx, tx, payment, map[string]interface{}{
			"status":    models.PaymentStatusPaid,
			"reference": payload.Data.Reference,
			"paid_at":   &now,
		}); err != nil {
			return err
		}

		// confirm order
		if err := s.paymentRepo.UpdateOrderStatus(ctx, tx, payment.OrderID, map[string]interface{}{
			"payment_status": models.PaymentStatusPaid,
			"status":         models.OrderStatusConfirmed,
		}); err != nil {
			return err
		}

		// deduct stock
		fullOrder, err := s.paymentRepo.GetOrderByID(ctx, tx, payment.OrderID)
		if err != nil {
			return err
		}

		err = s.inventory.DeductStock(ctx, tx, fullOrder.OrderItems, payment.OrderID, payment.UserID)
		if err != nil {
			return err
		}

		// clear cart
		if err := s.paymentRepo.ClearCartByUserID(ctx, tx, payment.UserID); err != nil {
			return err
		}

		var itemsOrdered []events.OrderItemPayload
		for _, item := range fullOrder.OrderItems {
			itemsOrdered = append(itemsOrdered, events.OrderItemPayload{
				Name:     item.Menu.Name,
				Quantity: item.Quantity,
				Price:    item.Price,
			})
		}

		// write outbox
		orderConfirmation, _ := json.Marshal(events.OrderConfirmationPayload{
			UserID:    payment.UserID,
			Email:     fullOrder.User.Email,
			FirstName: fullOrder.User.FirstName,
			OrderID:   payment.OrderID,
			Amount:    int64(payment.Amount),
			Reference: payment.Reference,
			Items:     itemsOrdered,
		})

		return s.paymentRepo.CreateOutboxEvent(ctx, tx, &models.OutboxEvent{
			EventType: events.ChannelOrderConfirmation,
			Payload:   string(orderConfirmation),
			Status:    "pending",
		})
	})
	if err != nil {
		return err
	}

	// The list of orders is in transactions
	// Fetch the orders lits for go routine to use in publishing
	orderForEmail, err := s.paymentRepo.GetOrderByID(ctx, nil, payment.OrderID)
	if err != nil {
		log.Printf("failed to fetch order for email: %v", err)
	}

	go func() {
		// Build the order items
		var itemsOrdered []events.OrderItemPayload
		if orderForEmail != nil {
			for _, item := range orderForEmail.OrderItems {
				itemsOrdered = append(itemsOrdered, events.OrderItemPayload{
					Name:     item.Menu.Name,
					Quantity: item.Quantity,
					Price:    item.Price,
				})
			}
		}
		err := s.eventPublisher.PublishMessage(
			events.ChannelOrderConfirmation,
			&events.OrderConfirmationPayload{
				OrderID:   payment.OrderID,
				UserID:    payment.UserID,
				Amount:    int64(payment.Amount),
				Reference: reference,
			},
			map[string]string{"Priority": "Important Mail"},
		)
		if err != nil {
			return
		}
		s.paymentRepo.MarkOutboxPublished(ctx, events.ChannelOrderConfirmation)
	}()

	return nil
}

func (s *PaymentService) ProcessTakeoutRefund(ctx context.Context, orderID uint, notes string) error {
	order, err := s.paymentRepo.GetOrderByID(ctx, nil, orderID)
	if err != nil {
		return domain.ErrOrderNotFound
	}

	payment, err := s.paymentRepo.GetPaymentByOrderID(ctx, orderID)
	if err != nil {
		return domain.ErrPaymentNotFound
	}

	if payment.Status != models.PaymentStatusPaid {
		return domain.ErrOrderNotPaid
	}

	// idempotency check
	_, err = s.paymentRepo.GetRefundByOrderID(ctx, orderID)
	if err == nil {
		return domain.ErrRefundAlreadyProcessed
	}

	response, err := s.callPaystackRefund(payment.Reference, int64(order.TotalAmount*100))
	if err != nil {
		return err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		refund := &models.Refund{
			OrderID:        orderID,
			PaymentID:      payment.ID,
			Amount:         order.TotalAmount,
			Reason:         notes,
			Reference:      payment.Reference,
			Currency:       "NGN",
			IdempotencyKey: uuid.New().String(),
			Status:         response.Data.Status,
			ProcessedAt:    &now,
		}

		if err := s.paymentRepo.CreateRefund(ctx, tx, refund); err != nil {
			return err
		}

		if err := s.paymentRepo.UpdateOrderStatus(ctx, tx, orderID, map[string]interface{}{
			"status":         models.OrderStatusCancelled,
			"payment_status": models.PaymentStatusRefunded,
		}); err != nil {
			return err
		}

		orderRefund, _ := json.Marshal(events.OrderRefundedPayload{
			UserID:    order.User.ID,
			Email:     order.User.Email,
			FirstName: order.User.FirstName,
			OrderID:   orderID,
			Amount:    int64(payment.Amount),
			Reference: payment.Reference,
			Reason:    notes,
		})

		return s.paymentRepo.CreateOutboxEvent(ctx, tx, &models.OutboxEvent{
			EventType: events.ChannelOrderRefunded,
			Payload:   string(orderRefund),
			Status:    "pending",
		})
	})
	if err != nil {
		return err
	}

	go func() {
		err := s.eventPublisher.PublishMessage(
			events.ChannelOrderRefunded,
			&events.OrderRefundedPayload{
				UserID:    order.User.ID,
				OrderID:   orderID,
				FirstName: order.User.FirstName,
				Amount:    int64(payment.Amount),
				Reference: payment.Reference,
				Reason:    notes,
			},
			map[string]string{"Priority": "Important Mail"},
		)
		if err != nil {
			return
		}
		s.paymentRepo.MarkOutboxPublished(ctx, events.ChannelOrderRefunded)
	}()

	return nil
}

func (s *PaymentService) GetPaymentByReference(ctx context.Context, reference string) (*dto.PaymentResponse, error) {
	payment, err := s.paymentRepo.GetPaymentByReference(ctx, reference)
	if err != nil {
		return nil, domain.ErrPaymentNotFound
	}
	return mapper.PaymentResponse(payment), nil
}

func (s *PaymentService) GetPaymentDetails(ctx context.Context, paymentID uint) (*dto.PaymentResponse, error) {
	payment, err := s.paymentRepo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return nil, domain.ErrPaymentNotFound
	}
	return mapper.PaymentResponse(payment), nil
}

func (s *PaymentService) GetRefundDetails(ctx context.Context, refundID uint) (*dto.RefundResponse, error) {
	refund, err := s.paymentRepo.GetRefundByID(ctx, refundID)
	if err != nil {
		return nil, domain.ErrRefundNotFound
	}
	return &dto.RefundResponse{
		ID:        refund.ID,
		OrderID:   refund.OrderID,
		Amount:    refund.Amount,
		Reason:    refund.Reason,
		Reference: refund.Reference,
		Status:    refund.Status,
		Currency:  refund.Currency,
	}, nil
}

func (s *PaymentService) GetAllPaymentHistory(ctx context.Context, userID uint, page, pageSize int) ([]*dto.PaymentResponse, *utils.PaginatedMeta, error) {
	cacheKey := fmt.Sprintf("user:payments:%d:p:%d:s:%d", userID, page, pageSize)

	cached, err := s.redisStore.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cachedResponse struct {
			Data []*dto.PaymentResponse `json:"data"`
			Meta *utils.PaginatedMeta   `json:"meta"`
		}
		if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
			return cachedResponse.Data, cachedResponse.Meta, nil
		}
	}

	payments, total, err := s.paymentRepo.GetAllByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, nil, domain.ErrPaymentNotFound
	}

	response := make([]*dto.PaymentResponse, 0, len(payments))
	for _, payment := range payments {
		response = append(response, mapper.PaymentResponse(&payment))
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &utils.PaginatedMeta{Page: page, Limit: pageSize, Total: total, TotalPages: totalPages}

	cacheData := struct {
		Data []*dto.PaymentResponse `json:"data"`
		Meta *utils.PaginatedMeta   `json:"meta"`
	}{Data: response, Meta: meta}
	data, _ := json.Marshal(&cacheData)
	s.redisStore.Set(ctx, cacheKey, string(data), 1*time.Hour)

	return response, meta, nil
}
