package repositories

import (
	"context"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

// ─── Order

func (r *PaymentRepository) GetOrderByIDAndUser(ctx context.Context, tx *gorm.DB, orderID, userID uint) (*models.Order, error) {
	var order models.Order
	err := r.getDB(tx).WithContext(ctx).Preload("User").
		Where("user_id = ? AND id = ?", userID, orderID).
		First(&order).Error
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}
	return &order, nil
}

func (r *PaymentRepository) GetOrderByID(ctx context.Context, tx *gorm.DB, orderID uint) (*models.Order, error) {
	var order models.Order
	err := r.getDB(tx).WithContext(ctx).
		Preload("User").Preload("OrderItems.Menu").
		Where("id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}
	return &order, nil
}

func (r *PaymentRepository) UpdateOrderStatus(ctx context.Context, tx *gorm.DB, orderID uint, updates map[string]interface{}) error {
	return r.getDB(tx).WithContext(ctx).Model(&models.Order{}).
		Where("id = ?", orderID).Updates(updates).Error
}

// ─── Payment

func (r *PaymentRepository) GetPaymentByOrderID(ctx context.Context, orderID uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		return nil, domain.ErrPaymentNotFound
	}
	return &payment, nil
}

func (r *PaymentRepository) GetPaymentByReference(ctx context.Context, reference string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).Preload("User").Where("reference = ?", reference).First(&payment).Error
	if err != nil {
		return nil, domain.ErrPaymentNotFound
	}
	return &payment, nil
}

func (r *PaymentRepository) GetPaymentByID(ctx context.Context, paymentID uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).Preload("Order").Where("id = ?", paymentID).First(&payment).Error
	if err != nil {
		return nil, domain.ErrPaymentNotFound
	}
	return &payment, nil
}

func (r *PaymentRepository) UpdatePayment(ctx context.Context, tx *gorm.DB, payment *models.Payment, updates map[string]interface{}) error {
	return r.getDB(tx).WithContext(ctx).Model(payment).Updates(updates).Error
}

func (r *PaymentRepository) GetAllByUserID(ctx context.Context, userID uint, page, pageSize int) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64
	offset := utils.Pagination(page, pageSize)

	r.db.WithContext(ctx).Model(&models.Payment{}).
		Where("user_id = ?", userID).Count(&total)

	err := r.db.WithContext(ctx).
		Preload("Order").Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&payments).Error
	if err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

func (r *PaymentRepository) Create(ctx context.Context, tx *gorm.DB, payment *models.Payment) error {
	return r.getDB(tx).WithContext(ctx).Create(payment).Error
}

// ─── Cart

func (r *PaymentRepository) ClearCartByUserID(ctx context.Context, tx *gorm.DB, userID uint) error {
	var cart models.Cart
	if err := r.getDB(tx).WithContext(ctx).
		Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return nil // cart not found → fine
	}
	return r.getDB(tx).WithContext(ctx).Unscoped().
		Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error
}

// ─── Refund

func (r *PaymentRepository) GetRefundByOrderID(ctx context.Context, orderID uint) (*models.Refund, error) {
	var refund models.Refund
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&refund).Error
	if err != nil {
		return nil, err
	}
	return &refund, nil
}

func (r *PaymentRepository) GetRefundByID(ctx context.Context, refundID uint) (*models.Refund, error) {
	var refund models.Refund
	err := r.db.WithContext(ctx).
		Preload("Order").Preload("Payment").
		Where("id = ?", refundID).First(&refund).Error
	if err != nil {
		return nil, domain.ErrRefundNotFound
	}
	return &refund, nil
}

func (r *PaymentRepository) CreateRefund(ctx context.Context, tx *gorm.DB, refund *models.Refund) error {
	return r.getDB(tx).WithContext(ctx).Create(refund).Error
}

// ─── Outbox

func (r *PaymentRepository) CreateOutboxEvent(ctx context.Context, tx *gorm.DB, event *models.OutboxEvent) error {
	return r.getDB(tx).WithContext(ctx).Create(event).Error
}

func (r *PaymentRepository) MarkOutboxPublished(ctx context.Context, eventType string) error {
	return r.db.WithContext(ctx).Model(&models.OutboxEvent{}).
		Where("id = ? AND status = ? AND event_type = ?", "pending", eventType).
		Updates(map[string]interface{}{
			"status":       "published",
			"processed_at": time.Now(),
		}).Error
}

func (r *PaymentRepository) RecheckPaymentWithReference(ctx context.Context, reference string) error {
	var fresh models.Payment
	if err := r.db.WithContext(ctx).Where("reference = ?", reference).First(&fresh).Error; err != nil {
		return domain.ErrPaymentNotFound
	}

	if fresh.Status == models.PaymentStatusPaid {
		return nil
	}

	return nil
}
