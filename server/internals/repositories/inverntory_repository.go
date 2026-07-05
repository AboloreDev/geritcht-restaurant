package repositories

import (
	"context"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *InventoryRepository) GetRecipesByMenuItemID(ctx context.Context, tx *gorm.DB, menuItemID uint) ([]models.MenuItemIngredient, error) {
	var recipes []models.MenuItemIngredient
	err := r.getDB(tx).WithContext(ctx).
		Where("menu_item_id = ?", menuItemID).Find(&recipes).Error
	return recipes, err
}

func (r *InventoryRepository) GetIngredientByID(ctx context.Context, tx *gorm.DB, ingredientID uint) (*models.Ingredient, error) {
	var ingredient models.Ingredient
	err := r.getDB(tx).WithContext(ctx).First(&ingredient, ingredientID).Error
	if err != nil {
		return nil, err
	}
	return &ingredient, nil
}

func (r *InventoryRepository) DeductIngredientStock(ctx context.Context, tx *gorm.DB, ingredientID uint, required float64) (int64, error) {
	result := r.getDB(tx).WithContext(ctx).Exec(
		"UPDATE ingredients SET current_stock = current_stock - ? WHERE id = ? AND current_stock >= ?",
		required, ingredientID, required,
	)
	return result.RowsAffected, result.Error
}

func (r *InventoryRepository) CreateStockMovement(ctx context.Context, tx *gorm.DB, movement *models.StockMovement) error {
	return r.getDB(tx).WithContext(ctx).Create(movement).Error
}

func (r *InventoryRepository) GetAdminUser(ctx context.Context, tx *gorm.DB) (*models.User, error) {
	var admin models.User
	err := r.getDB(tx).WithContext(ctx).Where("role = ?", models.RoleAdmin).First(&admin).Error
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	return &admin, nil
}

func (r *InventoryRepository) GetLowStockIngredients(ctx context.Context, tx *gorm.DB) ([]models.Ingredient, error) {
	var ingredients []models.Ingredient
	err := r.getDB(tx).WithContext(ctx).
		Where("current_stock <= min_threshold AND min_threshold > 0").
		Find(&ingredients).Error
	return ingredients, err
}

func (r *InventoryRepository) GetOutOfStockIngredients(ctx context.Context, tx *gorm.DB) ([]models.Ingredient, error) {
	var ingredients []models.Ingredient
	err := r.getDB(tx).WithContext(ctx).
		Where("current_stock <= 0").Find(&ingredients).Error
	return ingredients, err
}

func (r *InventoryRepository) CreateOutboxEvent(ctx context.Context, tx *gorm.DB, event *models.OutboxEvent) error {
	return r.getDB(tx).WithContext(ctx).Create(event).Error
}

func (r *InventoryRepository) MarkOutboxPublished(ctx context.Context, tx *gorm.DB, outboxID uint) error {
	return r.getDB(tx).WithContext(ctx).Model(&models.OutboxEvent{}).
		Where("id = ?", outboxID).
		Updates(map[string]interface{}{
			"status":       "published",
			"processed_at": time.Now(),
		}).Error
}

func (r *InventoryRepository) GetMenuItemIDsByIngredient(ctx context.Context, tx *gorm.DB, ingredientID uint) ([]uint, error) {
	var menuItemIDs []uint
	err := r.getDB(tx).WithContext(ctx).Model(&models.MenuItemIngredient{}).
		Where("ingredient_id = ?", ingredientID).
		Pluck("menu_item_id", &menuItemIDs).Error
	return menuItemIDs, err
}

func (r *InventoryRepository) DisableMenuItems(ctx context.Context, tx *gorm.DB, menuItemIDs []uint) error {
	return r.getDB(tx).WithContext(ctx).Model(&models.Menu{}).
		Where("id IN ?", menuItemIDs).
		Update("is_available", false).Error
}
