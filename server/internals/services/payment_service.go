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
	req, err := http.NewRequest("POST", s.config.Paystack.PaystakcinitialiseURL, bytes.NewBuffer(body))
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
