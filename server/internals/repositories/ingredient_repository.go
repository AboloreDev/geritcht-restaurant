package repositories

import (
	"context"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
)

type IngredientRepository struct {
	db *gorm.DB
}

func NewIngredientRepository(db *gorm.DB) *IngredientRepository {
	return &IngredientRepository{
		db: db,
	}
}

func (r *IngredientRepository) GetIngredientByName(ctx context.Context, name string) (*models.Ingredient, error) {
	var ingredient models.Ingredient

	err := r.db.WithContext(ctx).Where("name = ?", name).First(&ingredient).Error
	if err != nil {
		return nil, domain.ErrConflict
	}

	return &ingredient, nil
}

func (r *IngredientRepository) CreateIngredient(ctx context.Context, ingredient *models.Ingredient) error {
	return r.db.WithContext(ctx).Create(ingredient).Error
}

func (r *IngredientRepository) GetIngredientByID(ctx context.Context, ingredientID uint) (*models.Ingredient, error) {
	var ingredient models.Ingredient

	err := r.db.WithContext(ctx).Where("id = ?", ingredientID).First(&ingredient).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, domain.ErrIngredientNotFound
	}
	return &ingredient, nil
}

func (r *IngredientRepository) UpdateIngredient(ctx context.Context, ingredient *models.Ingredient) error {
	return r.db.WithContext(ctx).Save(ingredient).Error
}

func (r *IngredientRepository) IngredientCount(ctx context.Context, ingredientID uint) (int64, error) {
	var count int64

	r.db.WithContext(ctx).Model(&models.MenuItemIngredient{}).Where("ingredient_id = ?", ingredientID).Count(&count)

	return count, nil
}

func (r *IngredientRepository) DeleteIngredient(ctx context.Context, ingredientID uint) error {
	var ingredient models.Ingredient

	result := r.db.WithContext(ctx).Where("id = ?", ingredientID).Delete(&ingredient)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *IngredientRepository) GetAllIngredients(ctx context.Context, page, pageSize int) ([]models.Ingredient, int64, error) {
	var ingredients []models.Ingredient
	var count int64

	offset := utils.Pagination(page, pageSize)

	r.db.Model(&models.Ingredient{}).Count(&count)

	err := r.db.
		WithContext(ctx).
		Order("created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&ingredients).Error
	if err != nil {
		return nil, 0, err
	}

	return ingredients, count, nil
}

func (r *IngredientRepository) CompareCurrentStockAgainstMinTheshold(ctx context.Context) ([]models.Ingredient, error) {
	var ingredients []models.Ingredient

	err := r.db.WithContext(ctx).Where("current_stock <= min_threshold").Find(&ingredients).Error
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (r *IngredientRepository) UpdateThreshHoldLimit(ctx context.Context, ingredientID uint, threshHold float64) error {
	ingredient, err := r.GetIngredientByID(ctx, ingredientID)
	if err != nil {
		return err
	}

	err = r.db.Model(ingredient).
		WithContext(ctx).
		Where("id = ?", ingredientID).
		Update("min_threshold", threshHold).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *IngredientRepository) TsvectorSearchIngredients(ctx context.Context, req *dto.IngredientSearchRequest) ([]models.IngredientWithRank, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	offset := utils.Pagination(req.Page, req.Limit)

	// build query
	query := r.db.Model(&models.Ingredient{}).WithContext(ctx).
		Select("ingredients.*, ts_rank(search_vector, plainto_tsquery('english', ?)) AS rank", req.Query).
		Where("search_vector @@ to_tsquery('english', ? || ':*')", req.Query).
		Offset(offset).Limit(req.Limit)

	if req.MinThreshold != nil {
		query.Where("min_threshold >= ?", req.MinThreshold)
	}
	if req.CurrentStock != nil {
		query.Where("current_stock >= ?", req.CurrentStock)
	}

	var count int64
	query.Count(&count)

	// Execute query with ranking
	// Crreate rank struct

	var rows []models.IngredientWithRank
	err :=
		query.Order("rank DESC, created_at DESC").
			Preload("StockMovement").
			Preload("MenuItemIngredients").
			Offset(offset).Limit(req.Limit).
			Find(&rows).Error

	if err != nil {
		return nil, 0, err
	}

	return rows, count, nil
}
