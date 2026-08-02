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
		Preload("OrderItems.Menu.Images").
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
		Preload("OrderItems.Menu.Images").Preload("User").Preload("Payment").
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
		Preload("OrderItems.Menu.Images").Preload("User").Preload("Payment").
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
	query.Session(&gorm.Session{}).Count(&total)

	err := query.
		Preload("OrderItems.Menu.Images").Preload("User").Preload("Payment").
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

// TsVector
func (r *OrderRepository) TsvectorSearchOrders(ctx context.Context, req *dto.OrderSearchRequest) ([]models.OrderWithRank, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	offset := utils.Pagination(req.Page, req.Limit)

	// build query
	query := r.db.Model(&models.Order{}).WithContext(ctx).Preload("User").
		Select("orders.*, ts_rank(search_vector, plainto_tsquery('english', ?)) AS rank", req.Query).
		Where("search_vector @@ to_tsquery('english', ? || ':*')", req.Query).
		Offset(offset).Limit(req.Limit)

	var count int64
	query.Count(&count)

	// Execute query with ranking
	// Crreate rank struct

	var rows []models.OrderWithRank
	err :=
		query.Order("rank DESC, created_at DESC").
			Offset(offset).Limit(req.Limit).
			Find(&rows).Error

	if err != nil {
		return nil, 0, err
	}

	return rows, count, nil
}
