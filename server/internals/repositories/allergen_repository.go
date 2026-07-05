package repositories

import (
	"context"
	"errors"

	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
)

type AllergenRepository struct {
	db *gorm.DB
}

func NewAllergenRepository(db *gorm.DB) *AllergenRepository {
	return &AllergenRepository{db: db}
}

func (r *AllergenRepository) Create(ctx context.Context, allergen *models.Allergen) error {
	return r.db.WithContext(ctx).Create(allergen).Error
}

func (r *AllergenRepository) GetByName(ctx context.Context, name string) (*models.Allergen, error) {
	var allergen models.Allergen
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&allergen).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &allergen, nil
}

func (r *AllergenRepository) GetByID(ctx context.Context, allergenID uint) (*models.Allergen, error) {
	var allergen models.Allergen
	err := r.db.WithContext(ctx).Where("id = ?", allergenID).First(&allergen).Error
	if err != nil {
		return nil, err
	}
	return &allergen, nil
}

func (r *AllergenRepository) GetAll(ctx context.Context, page, pageSize int) ([]models.Allergen, int64, error) {
	var allergens []models.Allergen
	var count int64

	offset := utils.Pagination(page, pageSize)

	r.db.WithContext(ctx).Model(&models.Allergen{}).Count(&count)

	err := r.db.WithContext(ctx).
		Offset(offset).Limit(pageSize).
		Find(&allergens).Error
	if err != nil {
		return nil, 0, err
	}

	return allergens, count, nil
}

func (r *AllergenRepository) Update(ctx context.Context, allergen *models.Allergen) error {
	return r.db.WithContext(ctx).Save(allergen).Error
}

func (r *AllergenRepository) Delete(ctx context.Context, allergenID uint) error {
	result := r.db.WithContext(ctx).Where("id = ?", allergenID).Delete(&models.Allergen{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *AllergenRepository) CountMenuItemsUsingAllergen(ctx context.Context, allergenID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Menu{}).
		Where("allergens = ?", allergenID).
		Count(&count).Error
	return count, err
}
