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
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PaymentService struct {
	db             *gorm.DB
	redisStore     interfaces.Cacher
	eventPublisher interfaces.Publisher
	config         *config.Config
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

func NewPaymentService(
	db *gorm.DB,
	redisStore interfaces.Cacher,
	eventPublisher interfaces.Publisher,
	config *config.Config,
) *PaymentService {
	return &PaymentService{
		db:             db,
		redisStore:     redisStore,
		eventPublisher: eventPublisher,
		config:         config,
	}
}

func (s *PaymentService) callPaystackInitialize(email string, amount int64, reference string, orderID uint) (*PaystackInitialiseResponse, error) {
	url := "https://api.paystack.co/transaction/initialize"

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
	req, err := http.NewRequest("POST",url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Set the Req headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.Paystack.PaystackSecretKey)

	client := &http.Client{Timeout: 10 * time.Second}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Parse response body to struct
	// Reading using NewDecoder and Encode
	var result PaystackInitialiseResponse

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if !result.Status {
		return nil, err
	}

	return &result, nil
}

func (s *PaymentService) callPaystackVerify(reference string) (*PaystackVerifyResponse, error) {
	url := fmt.Sprintf("https://api.paystack.co/transaction/verify/%s", reference)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.config.Paystack.PaystackSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
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

	url := "https://api.paystack.co/refund"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.config.Paystack.PaystackSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
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
	return hmac.Equal([]byte(expectedMac), []byte(signature))
}

func (s *PaymentService) InitialisePayment(userID uint, req *dto.InitializePaymentRequest) (*dto.InitializePaymentResponse, error) {
	var order models.Order
	var payment models.Payment

	err := s.db.Preload("User").
		Where("user_id = ? AND id = ?", userID, req.OrderID).
		First(&order).Error
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}

	if order.Status != models.OrderStatusPending {
		return nil, domain.ErrInvalidOrderStatus
	}

	if order.PaymentStatus != models.PaymentStatusUnpaid {
		return nil, domain.ErrOrderAlreadyPaid
	}

	err = s.db.Where("order_id = ?", order.ID).First(&payment).Error
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
		return nil, fmt.Errorf("paystack initialization failed: %w", err)
	}

	err = s.db.Model(&payment).Updates(map[string]interface{}{
		"status": models.PaymentStatusPending,
	}).Error
	if err != nil {
		return nil, err
	}

	return &dto.InitializePaymentResponse{
		AuthorizationURL: response.Data.AuthorizationURL,
		Reference:        response.Data.Reference,
		Payment:          *mapper.PaymentResponse(&payment),
	}, nil
}

func (s *PaymentService) VerifyPayment(req *dto.VerifyPaymentRequest) (*dto.PaymentResponse, error) {
	var payment models.Payment
	err := s.db.Where("reference = ?", req.Reference).First(&payment).Error
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

func (s *PaymentService) HandlePaystackWebhook(body []byte, signature string) error {
	// Verify the signarure
	if !s.verifySignature(body, signature) {
		return domain.ErrInvalidSignature
	}

	var payload PaystackWebhookPayload
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return err
	}

	if payload.Event != "charge.success" {
		return nil
	}

	reference := payload.Data.Reference

	var payment models.Payment
	err = s.db.Where("reference = ?", reference).First(&payment).Error
	if err != nil {
		return domain.ErrPaymentNotFound
	}

	if payment.Status == models.PaymentStatusPaid {
		return nil
	}

	// prevents two webhook deliveries processing simultaneously
	lockKey := fmt.Sprintf("lock:webhook:payment:%s", reference)
	lockValue, _ := json.Marshal(reference)

	err = s.redisStore.Hold(context.Background(), lockKey, lockValue, redis.SetArgs{
		Mode: "NX",
		TTL:  30 * time.Second,
	})
	if err != nil {
		return nil
	}

	defer func() {
		log.Println("Releasing lock")
		err := s.redisStore.Delete(context.Background(), lockKey)
		if err != nil {
			log.Printf("Failed to release lock: %v", err)
		}
	}()

	// Prevent partial payment
	// User pays 100 for 5000 order
	expectedAmount := int64(payment.Amount * 100)
	if payload.Data.Amount != expectedAmount {
		return domain.ErrPaymentAmountMismatch
	}

	// Atomic update
	err = s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		err := tx.Model(&payment).Updates(map[string]interface{}{
			"status":    models.PaymentStatusPaid,
			"reference": payload.Data.Reference,
			"paidAt":    &now,
		}).Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Order{}).
			Where("id = ?", payment.OrderID).
			Updates(map[string]interface{}{
				"payment_status": models.PaymentStatusPaid,
				"status":         models.OrderStatusConfirmed,
			}).Error
		if err != nil {
			return err
		}

		var cart models.Cart
		err = tx.Unscoped().Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error
		if err != nil {
			return err
		}

		orderConfirmation, _ := json.Marshal(events.OrderConfirmationPayload{
			UserID:    payment.UserID,
			Email:     payment.User.Email,
			FirstName: payment.User.FirstName,
			OrderID:   payment.OrderID,
			Amount:    int64(payment.Amount),
			Reference: payment.Reference,
		})

		err = tx.Create(&models.OutboxEvent{
			EventType: events.ChannelOrderConfirmation,
			Payload:   string(orderConfirmation),
			Status:    "pending",
			CreatedAt: time.Now(),
		}).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	// start goRoutine to publsh message
	// paystack doesnt timeout
	// if publsh fails, background outbox workers will retry
	go func() {
		err := s.eventPublisher.PublishMessage(
			events.ChannelOrderConfirmation,
			&events.OrderConfirmationPayload{
				Email:     payment.User.Email,
				FirstName: payment.User.FirstName,
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

		s.db.Model(&models.OutboxEvent{}).
			Where("status = ? AND event_type = ?", "pending", events.ChannelOrderConfirmation).
			Updates(map[string]interface{}{
				"status":      "published",
				"processedAt": time.Now(),
			})
	}()

	return nil
}

func (s *PaymentService) ProcessTakoutRefund(orderID uint, req *dto.ProcessRefundRequest) error {
	var order models.Order

	err := s.db.Preload("Payment").Where("id = ?", orderID).First(&order).Error
	if err != nil {
		return domain.ErrOrderNotFound
	}

	if order.Payment.ID == 0 {
		return domain.ErrPaymentNotFound
	}

	if order.Payment.Status != models.PaymentStatusPaid {
		return domain.ErrOrderNotPaid
	}

	// Idempotency check
	err = s.db.Model(&models.Refund{}).
		Where("order_id = ?", orderID).
		First(&models.Refund{}).Error
	if err == nil {
		return domain.ErrRefundAlreadyProcessed
	}

	response, err := s.callPaystackRefund(
		order.Payment.Reference,
		int64(order.TotalAmount*100),
	)
	if err != nil {
		return err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		refunds := models.Refund{
			OrderID:        orderID,
			PaymentID:      order.Payment.ID,
			Amount:         order.TotalAmount,
			Reason:         req.Notes,
			Reference:      order.Payment.Reference,
			Currency:       "NGN",
			IdempotencyKey: uuid.New().String(),
			Status:         response.Data.Status,
			ProcessedAt:    &now,
			CreatedAt:      now,
		}

		err = tx.Create(&refunds).Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Order{}).
			Where("id = ?", orderID).
			Updates(map[string]string{
				"status":         string(models.OrderStatusCancelled),
				"payment_status": string(models.PaymentStatusRefunded),
			}).Error
		if err != nil {
			return err
		}

		orderRefund, _ := json.Marshal(events.OrderRefundedPayload{
			UserID:    order.User.ID,
			Email:     order.User.Email,
			FirstName: order.User.FirstName,
			OrderID:   orderID,
			Amount:    int64(order.Payment.Amount),
			Reference: order.Payment.Reference,
		})

		err = tx.Create(&models.OutboxEvent{
			EventType: events.ChannelOrderConfirmation,
			Payload:   string(orderRefund),
			Status:    "pending",
			CreatedAt: time.Now(),
		}).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	go func() {
		err := s.eventPublisher.PublishMessage(
			events.ChannelOrderRefunded,
			&events.OrderRefundedPayload{
				UserID:    order.User.ID,
				Email:     order.User.Email,
				FirstName: order.User.FirstName,
				OrderID:   orderID,
				Amount:    int64(order.Payment.Amount),
				Reference: order.Payment.Reference,
			},
			map[string]string{"Priority": "Important Mail"},
		)
		if err != nil {
			return
		}

		s.db.Model(&models.OutboxEvent{}).
			Where("status = ? AND event_type = ?", "pending", events.ChannelOrderRefunded).
			Updates(map[string]interface{}{
				"status":      "published",
				"processedAt": time.Now(),
			})
	}()

	return nil

}

func (s *PaymentService) GetPaymentByReference(reference string) (*dto.PaymentResponse, error) {
	var payment models.Payment

	err := s.db.Preload("User").Where("reference = ?", reference).First(&payment).Error
	if err != nil {
		return nil, domain.ErrPaymentNotFound
	}

	return mapper.PaymentResponse(&payment), nil
}

func (s *PaymentService) GetAllPaymentHistory(userID uint, page, pageSize int) ([]*dto.PaymentResponse, *utils.PaginatedMeta, error) {
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
	var payments []models.Payment
	var total int64
	offset := utils.Pagination(page, pageSize)

	s.db.Model(models.Payment{}).Count(&total)

	err = s.db.Preload("Orders").
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&payments).Error
	if err != nil {
		return nil, nil, domain.ErrPaymentNotFound
	}

	response := make([]*dto.PaymentResponse, 0, len(payments))

	for _, payment := range payments {
		response = append(response, mapper.PaymentResponse(&payment))
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	meta := &utils.PaginatedMeta{
		Page:       page,
		Limit:      pageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	cacheData := struct {
		Data []*dto.PaymentResponse `json:"data"`
		Meta *utils.PaginatedMeta   `json:"meta"`
	}{Data: response, Meta: meta}

	// store in cache
	data, err := json.Marshal(&cacheData)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to set data: %d", err)
	}
	s.redisStore.Set(ctx, cacheKey, string(data), 1*time.Hour)

	return response, meta, nil
}

func (s *PaymentService) GetPaymentDetails(paymentID uint) (*dto.PaymentResponse, error) {
	var payment models.Payment
	err := s.db.Preload("Order").Where("id = ?", paymentID).First(&payment).Error
	if err != nil {
		return nil, domain.ErrPaymentNotFound
	}

	return mapper.PaymentResponse(&payment), nil
}

func (s *PaymentService) GetRefundDetails(refundID uint) (*dto.RefundResponse, error) {
	var refund models.Refund

	err := s.db.Preload("Order").
		Preload("Payment").
		Where("id = ? ", refundID).First(&refund).Error
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
