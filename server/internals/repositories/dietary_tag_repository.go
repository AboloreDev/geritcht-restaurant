package repositories

import (
	"context"
	"errors"

	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
)

type DietaryTagRepository struct {
	db *gorm.DB
}

func NewDietaryTagRepository(db *gorm.DB) *DietaryTagRepository {
	return &DietaryTagRepository{db: db}
}

func (r *DietaryTagRepository) Create(ctx context.Context, tag *models.DietaryTag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *DietaryTagRepository) GetByName(ctx context.Context, name string) (*models.DietaryTag, error) {
	var tag models.DietaryTag
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &tag, nil
}

func (r *DietaryTagRepository) GetByID(ctx context.Context, tagID uint) (*models.DietaryTag, error) {
	var tag models.DietaryTag
	err := r.db.WithContext(ctx).Where("id = ?", tagID).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *DietaryTagRepository) GetAll(ctx context.Context, page, pageSize int) ([]models.DietaryTag, int64, error) {
	var tags []models.DietaryTag
	var count int64

	offset := utils.Pagination(page, pageSize)

	r.db.WithContext(ctx).Model(&models.DietaryTag{}).Count(&count)

	err := r.db.WithContext(ctx).
		Offset(offset).Limit(pageSize).
		Find(&tags).Error
	if err != nil {
		return nil, 0, err
	}

	return tags, count, nil
}

func (r *DietaryTagRepository) Update(ctx context.Context, tag *models.DietaryTag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

func (r *DietaryTagRepository) Delete(ctx context.Context, tagID uint) error {
	result := r.db.WithContext(ctx).Where("id = ?", tagID).Delete(&models.DietaryTag{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *DietaryTagRepository) CountMenuItemsUsingTag(ctx context.Context, tagID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Menu{}).
		Where("dietary_tags = ?", tagID).
		Count(&count).Error
	return count, err
}
