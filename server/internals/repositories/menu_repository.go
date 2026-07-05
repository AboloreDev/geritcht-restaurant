package repositories

import (
	"context"

	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) GetCategoryByID(ctx context.Context, categoryID uint) (*models.MenuCategory, error) {
	var category models.MenuCategory
	err := r.db.WithContext(ctx).First(&category, categoryID).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *MenuRepository) GetAllergensByIDs(ctx context.Context, ids []uint) ([]models.Allergen, error) {
	var allergens []models.Allergen
	err := r.db.WithContext(ctx).Find(&allergens, ids).Error
	return allergens, err
}

func (r *MenuRepository) GetDietaryTagsByIDs(ctx context.Context, ids []uint) ([]models.DietaryTag, error) {
	var tags []models.DietaryTag
	err := r.db.WithContext(ctx).Find(&tags, ids).Error
	return tags, err
}

func (r *MenuRepository) CountByNameAndCategory(ctx context.Context, name string, categoryID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Menu{}).
		Where("name = ? AND menu_category_id = ?", name, categoryID).
		Count(&count).Error
	return count, err
}

func (r *MenuRepository) Create(ctx context.Context, menu *models.Menu) error {
	return r.db.WithContext(ctx).Create(menu).Error
}

func (r *MenuRepository) GetByID(ctx context.Context, menuID uint) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.WithContext(ctx).First(&menu, menuID).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *MenuRepository) GetByIDAvailable(ctx context.Context, menuID uint) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.WithContext(ctx).
		Preload("Images").Preload("Allergens").Preload("MenuCategory").
		Preload("MenuItemIngredients").Preload("DietaryTags").
		Where("id = ? AND is_available = ?", menuID, true).
		First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *MenuRepository) Update(ctx context.Context, menu *models.Menu) error {
	return r.db.WithContext(ctx).Save(menu).Error
}

func (r *MenuRepository) ReplaceAllergens(ctx context.Context, menu *models.Menu, allergens []models.Allergen) error {
	return r.db.WithContext(ctx).Model(menu).Association("Allergens").Replace(allergens)
}

func (r *MenuRepository) ReplaceDietaryTags(ctx context.Context, menu *models.Menu, tags []models.DietaryTag) error {
	return r.db.WithContext(ctx).Model(menu).Association("DietaryTags").Replace(tags)
}

func (r *MenuRepository) Delete(ctx context.Context, menuID uint) error {
	result := r.db.WithContext(ctx).Where("id = ?", menuID).Delete(&models.Menu{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *MenuRepository) GetAll(ctx context.Context, filter dto.MenuFilterRequest) ([]models.Menu, int64, error) {
	offset := utils.Pagination(filter.Page, filter.PageSize)

	query := r.db.WithContext(ctx).Model(&models.Menu{}).Where("is_available = ?", true)
	query = utils.ApplyMenuFilters(query, filter)

	var count int64
	query.Count(&count)

	query = utils.ApplyMenuSorting(query, filter)

	var menus []models.Menu
	err := query.
		Preload("Images").Preload("Allergens").
		Preload("MenuCategory").Preload("DietaryTags").
		Offset(offset).Limit(filter.PageSize).
		Find(&menus).Error
	if err != nil {
		return nil, 0, err
	}

	return menus, count, nil
}

// ─── Images

func (r *MenuRepository) CountImages(ctx context.Context, menuID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.MenuImage{}).
		Where("menu_id = ?", menuID).Count(&count).Error
	return count, err
}

func (r *MenuRepository) CreateImage(ctx context.Context, image *models.MenuImage) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *MenuRepository) GetImageByID(ctx context.Context, imageID uint) (*models.MenuImage, error) {
	var image models.MenuImage
	err := r.db.WithContext(ctx).First(&image, imageID).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *MenuRepository) DeleteImage(ctx context.Context, image *models.MenuImage) error {
	return r.db.WithContext(ctx).Delete(image).Error
}

func (r *MenuRepository) GetNextPrimaryImage(ctx context.Context, menuID uint, excludeID uint) (*models.MenuImage, error) {
	var image models.MenuImage
	err := r.db.WithContext(ctx).
		Where("menu_id = ? AND id != ?", menuID, excludeID).
		First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *MenuRepository) SetImagePrimary(ctx context.Context, image *models.MenuImage) error {
	return r.db.WithContext(ctx).Model(image).Update("is_primary", true).Error
}
