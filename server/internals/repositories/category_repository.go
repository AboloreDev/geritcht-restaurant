package repositories

import (
	"context"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, category *models.MenuCategory) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *CategoryRepository) GetByID(ctx context.Context, categoryID uint) (*models.MenuCategory, error) {
	var category models.MenuCategory
	err := r.db.WithContext(ctx).
		Where("id = ? AND is_active = ?", categoryID, true).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetByName(ctx context.Context, name string) (*models.MenuCategory, error) {
	var category models.MenuCategory
	err := r.db.WithContext(ctx).
		Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetAll(ctx context.Context, page, pageSize int) ([]models.MenuCategory, int64, error) {
	var categories []models.MenuCategory
	var count int64

	offset := utils.Pagination(page, pageSize)

	r.db.WithContext(ctx).Model(&models.MenuCategory{}).Count(&count)

	err := r.db.WithContext(ctx).
		Order("display_order ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *models.MenuCategory) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *CategoryRepository) Delete(ctx context.Context, categoryID uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", categoryID).Delete(&models.MenuCategory{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CategoryRepository) CountMenuItems(ctx context.Context, categoryID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Menu{}).
		Where("menu_category_id = ?", categoryID).
		Count(&count).Error
	return count, err
}

// TSvector search
func (r *CategoryRepository) TsvectorSearchCategories(ctx context.Context, req *dto.CategorySearchRequest) ([]models.MenuCategoryWithRank, int64, error) {

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	offset := utils.Pagination(req.Page, req.Limit)

	// build query
	query := r.db.Model(&models.MenuCategory{}).
		Select("menu_categories.*, ts_rank(search_vector, plainto_tsquery('english', ?)) AS rank", req.Query).
		Where("search_vector @@ to_tsquery('english', ? || ':*')", req.Query).
		Where("is_active = ?", true).
		Offset(offset).Limit(req.Limit)

	var count int64
	query.Count(&count)

	// Execute query with ranking
	var rows []models.MenuCategoryWithRank
	err :=
		query.Order("rank DESC, created_at DESC").
			Offset(offset).Limit(req.Limit).
			Find(&rows).Error

	if err != nil {
		return nil, 0, err
	}

	return rows, count, nil
}
