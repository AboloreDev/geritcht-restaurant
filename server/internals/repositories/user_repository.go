package repositories

import (
	"context"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).
		Create(user).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).
		Save(user).Error
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Delete(&models.User{}, id).Error
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	var users []*models.User

	err := r.db.WithContext(ctx).
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetByIdAndActive(ctx context.Context, id uint, active bool) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		Where("id = ? AND is_active = ?", id, active).
		First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAllByRole(ctx context.Context, role models.UserRole, page, pageSize int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64
	offset := utils.Pagination(page, pageSize)

	r.db.WithContext(ctx).Model(&models.User{}).
		Where("role = ?", role).Count(&total)

	err := r.db.WithContext(ctx).
		Where("role = ?", role).
		Offset(offset).Limit(pageSize).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) GetByIDAndRole(ctx context.Context, id uint, role models.UserRole) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("id = ? AND role = ? AND is_active = ?", id, role, true).
		First(&user).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &user, nil
}

func (r *UserRepository) UpdateActiveByRole(ctx context.Context, id uint, role models.UserRole, active bool) error {
	// active=false means deactivate, active=true means activate
	// query checks opposite active state to prevent no-op
	result := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ? AND role = ? AND is_active = ?", id, role, !active).
		Update("is_active", active)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
