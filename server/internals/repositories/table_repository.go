package repositories

import (
	"context"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
)

type TableRepository struct {
	db *gorm.DB
}

func NewTableRepository(db *gorm.DB) *TableRepository {
	return &TableRepository{db: db}
}

func (r *TableRepository) GetByName(ctx context.Context, name string) (*models.Table, error) {
	var table models.Table
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&table).Error
	if err != nil {
		return nil, err
	}
	return &table, nil
}

func (r *TableRepository) GetByID(ctx context.Context, tableID uint) (*models.Table, error) {
	var table models.Table
	err := r.db.WithContext(ctx).Where("id = ?", tableID).First(&table).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &table, nil
}

func (r *TableRepository) GetByIDWithRelations(ctx context.Context, tableID uint) (*models.Table, error) {
	var table models.Table
	err := r.db.WithContext(ctx).
		Preload("Reservations").Preload("Orders").
		Where("id = ?", tableID).First(&table).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &table, nil
}

func (r *TableRepository) Create(ctx context.Context, table *models.Table) error {
	return r.db.WithContext(ctx).Create(table).Error
}

func (r *TableRepository) Update(ctx context.Context, table *models.Table) error {
	return r.db.WithContext(ctx).Save(table).Error
}

func (r *TableRepository) Delete(ctx context.Context, tableID uint) error {
	result := r.db.WithContext(ctx).Where("id = ?", tableID).Delete(&models.Table{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TableRepository) GetAll(ctx context.Context, page, pageSize int) ([]models.Table, int64, error) {
	var tables []models.Table
	var count int64
	offset := utils.Pagination(page, pageSize)

	r.db.WithContext(ctx).Model(&models.Table{}).Count(&count)

	err := r.db.WithContext(ctx).
		Order("name ASC").
		Offset(offset).Limit(pageSize).
		Find(&tables).Error
	if err != nil {
		return nil, 0, err
	}

	return tables, count, nil
}
