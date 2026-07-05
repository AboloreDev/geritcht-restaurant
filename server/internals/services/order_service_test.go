package services

import (
	"context"
	"testing"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	redisStore "github.com/AboloreDev/geritcht-restaurant/internals/redis"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testOrderCtx = context.Background()

// ─── MockOrderRepository

type MockOrderRepository struct {
	order     *models.Order
	orders    []models.Order
	total     int64
	count     int64
	getErr    error
	createErr error
	updateErr error
	countErr  error
}

func (m *MockOrderRepository) Create(_ context.Context, tx *gorm.DB, order *models.Order) error {
	order.ID = 1
	return m.createErr
}
func (m *MockOrderRepository) GetByID(_ context.Context, tx *gorm.DB, orderID uint) (*models.Order, error) {
	return m.order, m.getErr
}
func (m *MockOrderRepository) GetByIDAndUser(_ context.Context, orderID, userID uint) (*models.Order, error) {
	return m.order, m.getErr
}
func (m *MockOrderRepository) GetAllByUser(_ context.Context, userID uint, page, pageSize int) ([]models.Order, int64, error) {
	return m.orders, m.total, m.getErr
}
func (m *MockOrderRepository) GetAll(_ context.Context, page, pageSize int) ([]models.Order, int64, error) {
	return m.orders, m.total, m.getErr
}
func (m *MockOrderRepository) UpdateStatus(_ context.Context, orderID uint, status models.OrderStatus) error {
	return m.updateErr
}
func (m *MockOrderRepository) CountByUserAndID(_ context.Context, orderID, userID uint) (int64, error) {
	return m.count, m.countErr
}

func newOrderService(orderRepo *MockOrderRepository) *OrderService {
	return NewOrderService(nil, orderRepo, nil, nil, redisStore.NewNopCache())
}

// ─── CancelTakeoutOrder Tests (table test — most important)
func TestCancelTakeoutOrder(t *testing.T) {
	userID := uint(1)

	tests := []struct {
		name        string
		order       *models.Order
		expectedErr error
	}{
		{
			name: "success - pending order cancels",
			order: &models.Order{
				ID: 1, UserID: &userID,
				Status: models.OrderStatusPending,
			},
			expectedErr: nil,
		},
		{
			name: "already cancelled",
			order: &models.Order{
				ID: 1, UserID: &userID,
				Status: models.OrderStatusCancelled,
			},
			expectedErr: domain.ErrAlreadyCancelled,
		},
		{
			name: "preparing - cannot cancel",
			order: &models.Order{
				ID: 1, UserID: &userID,
				Status: models.OrderStatusPreparing,
			},
			expectedErr: domain.ErrCannotCancel,
		},
		{
			name: "ready - cannot cancel",
			order: &models.Order{
				ID: 1, UserID: &userID,
				Status: models.OrderStatusReady,
			},
			expectedErr: domain.ErrCannotCancel,
		},
		{
			name: "completed - cannot cancel",
			order: &models.Order{
				ID: 1, UserID: &userID,
				Status: models.OrderStatusCompleted,
			},
			expectedErr: domain.ErrCannotCancel,
		},
		{
			name: "confirmed - refund processing",
			order: &models.Order{
				ID: 1, UserID: &userID,
				Status: models.OrderStatusConfirmed,
			},
			expectedErr: domain.ErrRefundIsProcessing,
		},
		{
			name:        "order not found",
			order:       nil,
			expectedErr: domain.ErrOrderNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getErr := error(nil)
			if tt.order == nil {
				getErr = domain.ErrOrderNotFound
			}

			service := newOrderService(&MockOrderRepository{
				order:  tt.order,
				getErr: getErr,
			})

			err := service.CancelTakeoutOrder(testOrderCtx, userID, 1)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestCancelTakeoutOrder_WrongUser(t *testing.T) {
	correctUserID := uint(1)

	service := newOrderService(&MockOrderRepository{
		order: &models.Order{
			ID:     1,
			UserID: &correctUserID,
			Status: models.OrderStatusPending,
		},
	})

	// attacker tries to cancel with different userID
	err := service.CancelTakeoutOrder(testOrderCtx, 999, 1)

	// repo would actually return not found since query filters by user_id
	// but if it somehow returned the order, forbidden check catches it
	assert.Error(t, err)
}

// ─── VerifyUserOrder Tests

func TestVerifyUserOrder(t *testing.T) {
	tests := []struct {
		name        string
		count       int64
		countErr    error
		expectedErr error
	}{
		{name: "order belongs to user", count: 1, expectedErr: nil},
		{name: "order not found", count: 0, expectedErr: domain.ErrOrderNotFound},
		{name: "db error", countErr: domain.ErrOrderNotFound, expectedErr: domain.ErrOrderNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newOrderService(&MockOrderRepository{
				count:    tt.count,
				countErr: tt.countErr,
			})

			err := service.VerifyUserOrder(testOrderCtx, 1, 1)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

// ─── GetTakeoutOrder Tests
func TestGetTakeoutOrder_Success(t *testing.T) {
	userID := uint(1)
	service := newOrderService(&MockOrderRepository{
		order: &models.Order{
			ID:            1,
			UserID:        &userID,
			Type:          models.OrderTypeTakeout,
			Status:        models.OrderStatusPending,
			TotalAmount:   3500,
			PaymentStatus: models.PaymentStatusUnpaid,
			Payment: &models.Payment{
				ID:        1,
				Reference: "ref_123",
				Amount:    3500,
				Status:    models.PaymentStatusUnpaid,
			},
			User: &models.User{
				ID:          userID,
				FirstName:   "John",
				LastName:    "Doe",
				Email:       "johndoe@example.com",
				PhoneNumber: "1234567890",
			},
			OrderItems: []models.OrderItem{},
		},
	})

	response, err := service.GetTakeoutOrder(testOrderCtx, 1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestGetTakeoutOrder_NotFound(t *testing.T) {
	service := newOrderService(&MockOrderRepository{
		getErr: domain.ErrOrderNotFound,
	})

	response, err := service.GetTakeoutOrder(testOrderCtx, 1, 999)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrOrderNotFound, err)
}

// ─── GetAllTakeoutOrders Tests

func TestGetAllTakeoutOrders_Success(t *testing.T) {
	userID := uint(1)
	service := newOrderService(&MockOrderRepository{
		orders: []models.Order{ // ← slice, each item is models.Order{}
			{
				ID:            1,
				UserID:        &userID,
				Type:          models.OrderTypeTakeout,
				Status:        models.OrderStatusPending,
				TotalAmount:   3500,
				PaymentStatus: models.PaymentStatusUnpaid,
				User: &models.User{
					ID:          1,
					Email:       "test@test.com",
					FirstName:   "John",
					LastName:    "Doe",
					PhoneNumber: "1234567890",
				},
				Payment: &models.Payment{
					ID:        1,
					Reference: "ref_123",
					Amount:    3500,
					Status:    models.PaymentStatusUnpaid,
				},
				OrderItems: []models.OrderItem{},
			},
			{
				ID:            2,
				UserID:        &userID,
				Type:          models.OrderTypeTakeout,
				Status:        models.OrderStatusConfirmed,
				TotalAmount:   5000,
				PaymentStatus: models.PaymentStatusPaid,
				User: &models.User{
					ID:          1,
					Email:       "test@test.com",
					FirstName:   "John",
					LastName:    "Doe",
					PhoneNumber: "1234567890",
				},
				Payment: &models.Payment{
					ID:        2,
					Reference: "ref_456",
					Amount:    5000,
					Status:    models.PaymentStatusPaid,
				},
				OrderItems: []models.OrderItem{},
			},
		},
		total: 2,
	})

	response, meta, err := service.GetAllTakeoutOrders(testOrderCtx, userID, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, int64(2), meta.Total)
}

// Create takeout order
func TestCreateTakeoutOrder_Success(t *testing.T) {
	userID := uint(1)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectBegin()
	mock.ExpectCommit()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

	cartRepo := &MockCartRepository{
		cart: &models.Cart{
			ID:     1,
			UserID: userID,
			CartItems: []models.CartItem{
				{
					ID:       1,
					CartID:   1,
					MenuID:   1,
					Quantity: 2,
					Menu: models.Menu{
						ID:          1,
						Name:        "Jollof Rice",
						Price:       3500,
						IsAvailable: true,
					},
				},
			},
		},
	}

	// order returned by GetByID after creation
	orderRepo := &MockOrderRepository{
		order: &models.Order{
			ID:            1,
			UserID:        &userID,
			Type:          models.OrderTypeTakeout,
			Status:        models.OrderStatusPending,
			TotalAmount:   7000,
			PaymentStatus: models.PaymentStatusUnpaid,
			Notes:         "no pepper",
			User: &models.User{
				ID:    userID,
				Email: "test@test.com",
			},
			Payment: &models.Payment{
				ID:        1,
				Reference: "ref_123",
				Amount:    7000,
				Status:    models.PaymentStatusUnpaid,
			},
			OrderItems: []models.OrderItem{
				{
					ID:       1,
					MenuID:   1,
					Quantity: 2,
					Price:    3500,
					Menu:     models.Menu{ID: 1, Name: "Jollof Rice"},
				},
			},
		},
	}

	service := &OrderService{
		db:          gormDB,
		orderRepo:   orderRepo,
		paymentRepo: &MockPaymentRepository{},
		cartRepo:    cartRepo,
		redisStore:  redisStore.NewNopCache(),
	}

	req := &dto.CreateTakeoutOrderRequest{Notes: "no pepper"}
	response, err := service.CreateTakeoutOrder(testOrderCtx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, float64(7000), response.TotalAmount)
	assert.Equal(t, string(models.OrderStatusPending), response.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTakeoutOrder_EmptyCart(t *testing.T) {
	userID := uint(1)
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectBegin()
	mock.ExpectCommit()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

	cartRepo := &MockCartRepository{
		cart: &models.Cart{
			ID:        1,
			UserID:    userID,
			CartItems: []models.CartItem{}, // empty
		},
	}

	service := &OrderService{
		db:          gormDB,
		orderRepo:   &MockOrderRepository{},
		paymentRepo: &MockPaymentRepository{},
		cartRepo:    cartRepo,
		redisStore:  redisStore.NewNopCache(),
	}

	req := &dto.CreateTakeoutOrderRequest{}
	response, err := service.CreateTakeoutOrder(testOrderCtx, userID, req)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrCartIsEmpty, err)
}

func TestCreateTakeoutOrder_CartNotFound(t *testing.T) {
	userID := uint(1)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectBegin()
	mock.ExpectCommit()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

	cartRepo := &MockCartRepository{
		getCartErr: domain.ErrCartNotFound,
	}

	service := &OrderService{
		db:          gormDB,
		orderRepo:   &MockOrderRepository{},
		paymentRepo: &MockPaymentRepository{},
		cartRepo:    cartRepo,
		redisStore:  redisStore.NewNopCache(),
	}

	req := &dto.CreateTakeoutOrderRequest{}
	response, err := service.CreateTakeoutOrder(testOrderCtx, userID, req)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrCartNotFound, err)
}

func TestCreateTakeoutOrder_UnavailableItem(t *testing.T) {
	userID := uint(1)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectBegin()
	mock.ExpectCommit()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

	cartRepo := &MockCartRepository{
		cart: &models.Cart{
			ID:     1,
			UserID: userID,
			CartItems: []models.CartItem{
				{
					MenuID:   1,
					Quantity: 1,
					Menu: models.Menu{
						ID:          1,
						IsAvailable: false, // unavailable
					},
				},
			},
		},
	}

	service := &OrderService{
		db:          gormDB,
		orderRepo:   &MockOrderRepository{},
		paymentRepo: &MockPaymentRepository{},
		cartRepo:    cartRepo,
		redisStore:  redisStore.NewNopCache(),
	}

	req := &dto.CreateTakeoutOrderRequest{}
	response, err := service.CreateTakeoutOrder(testOrderCtx, userID, req)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrMenuNotAvailable, err)
}

func TestProcessTakeoutRefund_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectCommit()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	userID := uint(1)

	repo := &MockPaymentRepository{
		order: &models.Order{
			ID:          1,
			UserID:      &userID,
			TotalAmount: 5000,
			User:        &models.User{ID: userID, Email: "test@test.com", FirstName: "Test"},
		},
		payment: &models.Payment{
			ID:        1,
			OrderID:   1,
			Amount:    5000,
			Reference: "ref_123",
			Status:    models.PaymentStatusPaid,
		},
		refundErr: gorm.ErrRecordNotFound, // no existing refund
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
                "data": {
                    "status": "pending",
                    "amount": 500000,
                    "reference": "ref_123"
                }
            }`,
		},
	}

	err = service.ProcessTakeoutRefund(testPaymentCtx, 1, "customer requested cancellation")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
