package repositories

import (
	"context"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, tx *gorm.DB, order *models.Order) error {
	return tx.WithContext(ctx).Create(order).Error
}

func (r *OrderRepository) GetByID(ctx context.Context, tx *gorm.DB, orderID uint) (*models.Order, error) {
	var order models.Order
	db := tx
	if db == nil {
		db = r.db
	}
	err := db.WithContext(ctx).
		Preload("OrderItems.Menu.MenuCategory").
		Preload("User").Preload("Payment").
		Where("id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetByIDAndUser(ctx context.Context, orderID, userID uint) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		Preload("OrderItems.Menu").Preload("User").Preload("Payment").
		Where("id = ? AND user_id = ? AND type = ?", orderID, userID, models.OrderTypeTakeout).
		First(&order).Error
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}
	return &order, nil
}

func (r *OrderRepository) GetAllByUser(ctx context.Context, userID uint, filter *dto.OrderFilterRequest) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64
	offset := utils.Pagination(filter.Page, filter.PageSize)

	query := r.db.WithContext(ctx).Model(&models.Order{}).
		Where("user_id = ? AND type = ?", userID, models.OrderTypeTakeout)

	query = utils.ApplyOrderFilters(query, filter)

	query.Count(&total)

	err := query.
		Preload("OrderItems.Menu").Preload("User").Preload("Payment").
		Order("created_at DESC").
		Offset(offset).Limit(filter.PageSize).
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepository) GetAll(ctx context.Context, filter *dto.OrderFilterRequest) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64
	offset := utils.Pagination(filter.Page, filter.PageSize)

	query := r.db.WithContext(ctx).Model(&models.Order{})
	query = utils.ApplyOrderFilters(query, filter)
	query.Count(&total)

	err := query.
		Preload("OrderItems.Menu").Preload("User").Preload("Payment").
		Order("created_at DESC").
		Offset(offset).Limit(filter.PageSize).
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID uint, status models.OrderStatus) error {
	return r.db.WithContext(ctx).Model(&models.Order{}).
		Where("id = ?", orderID).
		Update("status", status).Error
}

func (r *OrderRepository) CountByUserAndID(ctx context.Context, orderID, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Order{}).
		Where("id = ? AND user_id = ?", orderID, userID).
		Count(&count).Error
	return count, err
}
