package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/redis"
	redisStore "github.com/AboloreDev/geritcht-restaurant/internals/redis"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MockHTTPClient struct {
	Err        error
	Body       string
	StatusCode int
}

var testPaymentCtx = context.Background()

// MockPaymentRepository
type MockPaymentRepository struct {
	order    *models.Order
	payment  *models.Payment
	payments []models.Payment
	refund   *models.Refund
	total    int64

	orderErr   error
	paymentErr error
	refundErr  error
	updateErr  error
	createErr  error
}

func (m *MockPaymentRepository) GetOrderByIDAndUser(_ context.Context, _ *gorm.DB, orderID, userID uint) (*models.Order, error) {
	return m.order, m.orderErr
}
func (m *MockPaymentRepository) GetOrderByID(_ context.Context, _ *gorm.DB, orderID uint) (*models.Order, error) {
	return m.order, m.orderErr
}
func (m *MockPaymentRepository) UpdateOrderStatus(_ context.Context, _ *gorm.DB, orderID uint, updates map[string]interface{}) error {
	return m.updateErr
}
func (m *MockPaymentRepository) GetPaymentByOrderID(_ context.Context, orderID uint) (*models.Payment, error) {
	return m.payment, m.paymentErr
}
func (m *MockPaymentRepository) GetPaymentByReference(_ context.Context, reference string) (*models.Payment, error) {
	return m.payment, m.paymentErr
}
func (m *MockPaymentRepository) GetPaymentByID(_ context.Context, paymentID uint) (*models.Payment, error) {
	return m.payment, m.paymentErr
}
func (m *MockPaymentRepository) UpdatePayment(_ context.Context, _ *gorm.DB, payment *models.Payment, updates map[string]interface{}) error {
	return m.updateErr
}
func (m *MockPaymentRepository) GetAllByUserID(_ context.Context, userID uint, page, pageSize int) ([]models.Payment, int64, error) {
	return m.payments, m.total, m.paymentErr
}
func (m *MockPaymentRepository) Create(_ context.Context, _ *gorm.DB, payment *models.Payment) error {
	return m.createErr
}
func (m *MockPaymentRepository) ClearCartByUserID(_ context.Context, _ *gorm.DB, userID uint) error {
	return nil
}
func (m *MockPaymentRepository) GetRefundByOrderID(_ context.Context, orderID uint) (*models.Refund, error) {
	return m.refund, m.refundErr
}
func (m *MockPaymentRepository) GetRefundByID(_ context.Context, refundID uint) (*models.Refund, error) {
	return m.refund, m.refundErr
}
func (m *MockPaymentRepository) CreateRefund(_ context.Context, _ *gorm.DB, refund *models.Refund) error {
	return m.createErr
}
func (m *MockPaymentRepository) CreateOutboxEvent(_ context.Context, _ *gorm.DB, event *models.OutboxEvent) error {
	return nil
}
func (m *MockPaymentRepository) MarkOutboxPublished(_ context.Context, eventType string) error {
	return nil
}
func (m *MockPaymentRepository) RecheckPaymentWithReference(_ context.Context, reference string) error {
	return nil
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	return &http.Response{
		StatusCode: m.StatusCode,
		Body:       io.NopCloser(strings.NewReader(m.Body)),
	}, nil
}

func StatusCode(status int, body string) *MockHTTPClient {
	return &MockHTTPClient{
		StatusCode: status,
		Body:       body,
	}
}

func newPaymentService(repo *MockPaymentRepository) *PaymentService {
	return &PaymentService{
		db:             nil,
		paymentRepo:    repo,
		redisStore:     redisStore.NewNopCache(),
		eventPublisher: &MockPublisher{},
		config:         testAuthConfig,
		HTTPClient:     nil,
		inventoryRepo:  &MockInventoryRepository{},
		inventory:      InventoryService{},
	}
}

func Test_CallPaystackInitialize(t *testing.T) {
	mockHTTPClient := StatusCode(
		200,
		`{
            "status": true,
            "message": "Authorization URL created",
            "data": {
                "authorization_url": "https://checkout.paystack.com/test",
                "reference": "test_ref_123",
                "access_code": "test_code"
            }
        }`,
	)
	fakeService := NewPaymentService(
		nil, redis.NewNopCache(), nil, &config.Config{
			Paystack: config.PaystackConfig{
				PaystackSecretKey: "sk_test_fake",
				PaystackURL:       "https://api.paystack.co/transaction/initialize",
			},
		}, nil, mockHTTPClient, nil, InventoryService{},
	)

	resp, err := fakeService.callPaystackInitialize(
		"test@gmail.com",
		5000,
		"test_ref_123",
		1,
	)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp.Data.Reference != "test_ref_123" {
		t.Errorf("Expected reference test_ref_123, got %s", resp.Data.Reference)
	}
	if resp.Data.AuthorizationURL != "https://checkout.paystack.com/test" {
		t.Errorf("Expected authorization url https://checkout.paystack.com/test, got %s", resp.Data.AuthorizationURL)
	}
	if resp.Status != true {
		t.Errorf("Expected status true, got %t", resp.Status)
	}
	if resp.Message != "Authorization URL created" {
		t.Errorf("Expected message 'Authorization URL created', got %s", resp.Message)
	}
}

func Test_CallPaystackVerify(t *testing.T) {
	mockHTTPClient := StatusCode(
		200,
		`{
            "status": true,
            "message": "Payment Verification Success",
            "data": {
               "amount": 5000,
                "reference": "test_ref_123",
                "paid_at": "2023-01-01 00:00:00"
            }
        }`,
	)
	fakeService := NewPaymentService(
		nil, redis.NewNopCache(), nil, &config.Config{
			Paystack: config.PaystackConfig{
				PaystackSecretKey: "sk_test_fake",
				PaystackURL:       "https://api.paystack.co/transaction/verify/test_ref_123",
			},
		}, nil, mockHTTPClient, &MockPaymentRepository{}, InventoryService{},
	)

	response, err := fakeService.callPaystackVerify("test_ref_123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Data.Reference != "test_ref_123" {
		t.Errorf("Expected reference test_ref_123, got %s", response.Data.Reference)
	}
	if response.Status != true {
		t.Errorf("Expected status true, got %t", response.Status)
	}
	if response.Message != "Payment Verification Success" {
		t.Errorf("Expected message 'Payment Verification Success', got %s", response.Message)
	}
}

func Test_CallPaystackRefund(t *testing.T) {
	mockHTTPClient := StatusCode(
		200,
		`{
            "status": true,
            "message": "Refund request initiated",
            "data": {
               "amount": 5000,
                "reference": "test_ref_123",
                "paid_at": "2023-01-01 00:00:00"
            }
        }`,
	)

	fakeService := NewPaymentService(
		nil, redis.NewNopCache(), nil, &config.Config{
			Paystack: config.PaystackConfig{
				PaystackSecretKey: "sk_test_fake",
				PaystackURL:       "https://api.paystack.co/refund",
			},
		}, nil, mockHTTPClient, &MockPaymentRepository{}, InventoryService{},
	)

	response, err := fakeService.callPaystackRefund("test_ref_123", 5000)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Data.Reference != "test_ref_123" {
		t.Errorf("Expected reference test_ref_123, got %s", response.Data.Reference)
	}
	if response.Status != true {
		t.Errorf("Expected status true, got %t", response.Status)
	}
	if response.Message != "Refund request initiated" {
		t.Errorf("Expected message 'Refund request initiated', got %s", response.Message)
	}
	if response.Data.Amount != 5000 {
		t.Errorf("Expected amount 5000, got %d", response.Data.Amount)
	}
}

func Test_VerifySignature(t *testing.T) {
	secretKey := "sk_test_fake"
	body := "test_body"
	signature := "x-paystack-signature"

	mac := hmac.New(sha512.New, []byte(secretKey))
	mac.Write([]byte(body))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	hmac.Equal([]byte(expectedSignature), []byte(signature))

	tests := []struct {
		name      string
		signature string
		body      string
		expected  bool
	}{
		{
			name:      "valid signature",
			signature: expectedSignature,
			body:      body,
			expected:  true,
		},
		{
			name:      "invalid signature",
			signature: "invalid_signature",
			body:      body,
			expected:  false,
		},
		{
			name:      "empty signature",
			signature: "",
			body:      body,
			expected:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockHTTPClient := StatusCode(
				200,
				`{
					"status": true,
					"message": "Payment Verification Success",
					"data": {
					   "amount": 5000,
						"reference": "test_ref_123",
						"paid_at": "2023-01-01 00:00:00"
					}
				}`,
			)

			fakeService := NewPaymentService(
				nil, redis.NewNopCache(), nil, &config.Config{
					Paystack: config.PaystackConfig{
						PaystackSecretKey: secretKey,
					},
				}, nil, mockHTTPClient, nil, InventoryService{},
			)

			result := fakeService.verifySignature([]byte(tc.body), tc.signature)
			if result != tc.expected {
				t.Errorf("Expected %t, got %t", tc.expected, result)
			}
		})
	}
}

func TestInitialisePayment(t *testing.T) {
	tests := []struct {
		name        string
		order       *models.Order
		payment     *models.Payment
		orderErr    error
		paymentErr  error
		expectedErr error
	}{
		{
			name:        "order not found",
			orderErr:    domain.ErrOrderNotFound,
			expectedErr: domain.ErrOrderNotFound,
		},
		{
			name: "order not pending",
			order: &models.Order{
				ID:     1,
				Status: models.OrderStatusConfirmed,
				User: &models.User{
					ID:          1,
					FirstName:   "John",
					LastName:    "Doe",
					Email:       "johndoe@example.com",
					PhoneNumber: "1234567890",
				},
			},
			expectedErr: domain.ErrInvalidOrderStatus,
		},
		{
			name: "order already paid",
			order: &models.Order{
				ID:            1,
				Status:        models.OrderStatusPending,
				PaymentStatus: models.PaymentStatusPaid,
				User: &models.User{
					ID:          1,
					FirstName:   "John",
					LastName:    "Doe",
					Email:       "johndoe@example.com",
					PhoneNumber: "1234567890",
				},
			},
			expectedErr: domain.ErrOrderAlreadyPaid,
		},
		{
			name: "payment not found",
			order: &models.Order{
				ID:            1,
				Status:        models.OrderStatusPending,
				PaymentStatus: models.PaymentStatusUnpaid,
				User: &models.User{
					ID:          1,
					FirstName:   "John",
					LastName:    "Doe",
					Email:       "johndoe@example.com",
					PhoneNumber: "1234567890",
				},
			},
			paymentErr:  domain.ErrPaymentNotFound,
			expectedErr: domain.ErrPaymentNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newPaymentService(&MockPaymentRepository{
				order:      tt.order,
				payment:    tt.payment,
				orderErr:   tt.orderErr,
				paymentErr: tt.paymentErr,
			})

			req := &dto.InitializePaymentRequest{OrderID: 1}
			response, err := service.InitialisePayment(testPaymentCtx, 1, req)

			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedErr != nil {
				assert.Nil(t, response)
			}
		})
	}
}

func TestHandleWebhook_InvalidSignature(t *testing.T) {
	service := newPaymentService(&MockPaymentRepository{})

	err := service.HandlePaystackWebhook(
		testPaymentCtx,
		[]byte(`{"event":"charge.success","data":{"reference":"ref123"}}`),
		"invalidsignature",
	)

	assert.Equal(t, domain.ErrInvalidSignature, err)
}

func TestHandleWebhook_AlreadyPaid_Idempotency(t *testing.T) {
	service := newPaymentService(&MockPaymentRepository{
		payment: &models.Payment{
			ID:        1,
			Reference: "ref123",
			Status:    models.PaymentStatusPaid,
		},
	})

	body := buildWebhookBody("ref123")
	sig := buildValidSignature(body, testAuthConfig.Paystack.PaystackSecretKey)

	err := service.HandlePaystackWebhook(testPaymentCtx, body, sig)

	assert.NoError(t, err)
}

func TestHandleWebhook_AmountMismatch(t *testing.T) {
	service := newPaymentService(&MockPaymentRepository{
		payment: &models.Payment{
			ID:        1,
			Reference: "ref123",
			Amount:    5000, // expects 500000 kobo
			Status:    models.PaymentStatusUnpaid,
		},
	})

	// webhook sends wrong amount
	body := buildWebhookBodyWithAmount("ref123", 100) // sends 100 kobo instead of 500000
	sig := buildValidSignature(body, testAuthConfig.Paystack.PaystackSecretKey)

	err := service.HandlePaystackWebhook(testPaymentCtx, body, sig)

	assert.Equal(t, domain.ErrPaymentAmountMismatch, err)
}

func TestHandleWebhook_PaymentNotFound(t *testing.T) {
	service := newPaymentService(&MockPaymentRepository{
		paymentErr: domain.ErrPaymentNotFound,
	})

	body := buildWebhookBody("ref_unknown")
	sig := buildValidSignature(body, testAuthConfig.Paystack.PaystackSecretKey)

	err := service.HandlePaystackWebhook(testPaymentCtx, body, sig)

	assert.Equal(t, domain.ErrPaymentNotFound, err)
}

// ProcessTakeoutRefund Tests

func TestProcessTakeoutRefund(t *testing.T) {
	tests := []struct {
		name        string
		order       *models.Order
		payment     *models.Payment
		refund      *models.Refund
		orderErr    error
		paymentErr  error
		refundErr   error
		expectedErr error
	}{
		{
			name:        "order not found",
			orderErr:    domain.ErrOrderNotFound,
			expectedErr: domain.ErrOrderNotFound,
		},
		{
			name:        "payment not found",
			order:       &models.Order{ID: 1},
			paymentErr:  domain.ErrPaymentNotFound,
			expectedErr: domain.ErrPaymentNotFound,
		},
		{
			name:  "order not paid",
			order: &models.Order{ID: 1, TotalAmount: 5000},
			payment: &models.Payment{
				ID:     1,
				Status: models.PaymentStatusUnpaid, // not paid
			},
			refundErr:   domain.ErrRefundNotFound,
			expectedErr: domain.ErrOrderNotPaid,
		},
		{
			name:  "already refunded",
			order: &models.Order{ID: 1, TotalAmount: 5000},
			payment: &models.Payment{
				ID:     1,
				Status: models.PaymentStatusPaid,
			},
			refund:      &models.Refund{ID: 1}, // refund exists
			refundErr:   nil,
			expectedErr: domain.ErrRefundAlreadyProcessed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newPaymentService(&MockPaymentRepository{
				order:      tt.order,
				payment:    tt.payment,
				refund:     tt.refund,
				orderErr:   tt.orderErr,
				paymentErr: tt.paymentErr,
				refundErr:  tt.refundErr,
			})

			err := service.ProcessTakeoutRefund(testPaymentCtx, 1, "cancelled by customer")

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

// GetPaymentByReference Tests

func TestGetPaymentByReference_Success(t *testing.T) {
	service := newPaymentService(&MockPaymentRepository{
		payment: &models.Payment{
			ID:        1,
			Reference: "ref_123",
			Amount:    5000,
			Status:    models.PaymentStatusPaid,
		},
	})

	response, err := service.GetPaymentByReference(testPaymentCtx, "ref_123")

	assert.NoError(t, err)
	assert.Equal(t, "ref_123", response.Reference)
}

func TestGetPaymentByReference_NotFound(t *testing.T) {
	service := newPaymentService(&MockPaymentRepository{
		paymentErr: domain.ErrPaymentNotFound,
	})

	response, err := service.GetPaymentByReference(testPaymentCtx, "bad_ref")

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrPaymentNotFound, err)
}

// GetAllPaymentHistory Tests

func TestGetAllPaymentHistory_Success(t *testing.T) {
	service := newPaymentService(&MockPaymentRepository{
		payments: []models.Payment{
			{ID: 1, Reference: "ref1", Amount: 5000},
			{ID: 2, Reference: "ref2", Amount: 3500},
		},
		total: 2,
	})

	response, meta, err := service.GetAllPaymentHistory(testPaymentCtx, 1, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, int64(2), meta.Total)
}

func TestHandleWebhook_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectCommit()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	userID := uint(1)
	amount := float64(5000)
	reference := "ref_success_123"

	repo := &MockPaymentRepository{
		payment: &models.Payment{
			ID:        1,
			Reference: reference,
			Amount:    amount,
			UserID:    userID,
			OrderID:   1,
			Status:    models.PaymentStatusUnpaid, // not yet paid
		},
		order: &models.Order{
			ID:     1,
			UserID: &userID,
			User: &models.User{
				ID:        userID,
				Email:     "test@test.com",
				FirstName: "Test",
			},
			OrderItems: []models.OrderItem{
				{
					ID:       1,
					MenuID:   1,
					Quantity: 2,
					Price:    2500,
					Menu:     models.Menu{ID: 1, Name: "Jollof Rice"},
				},
			},
		},
	}

	service := &PaymentService{
		db:             gormDB,
		paymentRepo:    repo,
		redisStore:     redisStore.NewNopCache(),
		eventPublisher: &MockPublisher{},
		config:         testAuthConfig,
		HTTPClient: &MockHTTPClient{
			StatusCode: 200,
			Body: `{
                "status": true,
            }`,
		},
		inventoryRepo: &MockInventoryRepository{}, // empty inventory — DeductStock skips if no recipes
	}

	// build valid body and signature
	body := buildWebhookBodyWithAmount(reference, int64(amount*100))
	sig := buildValidSignature(body, testAuthConfig.Paystack.PaystackSecretKey)

	err = service.HandlePaystackWebhook(testPaymentCtx, body, sig)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Helpers for webhook tests
func buildWebhookBody(reference string) []byte {
	return []byte(fmt.Sprintf(`{
        "event": "charge.success",
        "data": {
            "reference": "%s",
            "amount": 500000,
            "status": "success"
        }
    }`, reference))
}

func buildWebhookBodyWithAmount(reference string, amount int64) []byte {
	return []byte(fmt.Sprintf(`{
        "event": "charge.success",
        "data": {
            "reference": "%s",
            "amount": %d,
            "status": "success"
        }
    }`, reference, amount))
}

func buildValidSignature(body []byte, secret string) string {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}
