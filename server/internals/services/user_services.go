package services

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) ConvertToUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Role:        string(user.Role),
		IsActive:    user.IsActive,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
	}
}

func (s *UserService) GetUserProfileService(userID uint) (*dto.UserResponse, error) {
	return s.getUserByRole(userID, models.RoleCustomer)
}

func (s *UserService) GetStaffProfileService(userID uint) (*dto.UserResponse, error) {
	return s.getUserByRole(userID, models.RoleStaff)
}

func (s *UserService) GetAllUsersService(page int, pageSize int) ([]*dto.UserResponse, *utils.PaginatedMeta, error) {
	var users []models.User
	var total int64

	offset := utils.Pagination(page, pageSize)

	err := s.db.Model(&models.User{}).
		Where("role = ?", models.RoleCustomer).
		Count(&total).Offset(offset).
		Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, nil, domain.ErrNotFound
	}

	response := make([]*dto.UserResponse, 0, len(users))

	for _, user := range users {
		response = append(response, s.ConvertToUserResponse(&user))
	}

	totalPages := int(total + int64(pageSize) - 1/int64(pageSize))

	meta := &utils.PaginatedMeta{
		Total:      total,
		Page:       page,
		Limit:      pageSize,
		TotalPages: totalPages,
	}

	return response, meta, nil
}

func (s *UserService) DeactivateUserService(userID uint) error {
	return s.deactivateUserByRole(userID, models.RoleCustomer)
}

func (s *UserService) DeactivateStaffService(userID uint) error {
	return s.deactivateUserByRole(userID, models.RoleStaff)
}

func (s *UserService) ActivateUserService(userID uint) error {
	return s.activateUserByRole(userID, models.RoleCustomer)
}

func (s *UserService) ActivateStaffService(userID uint) error {
	return s.activateUserByRole(userID, models.RoleStaff)
}

func (s *UserService) GetAllStaffService(page int, pageSize int) ([]*dto.UserResponse, *utils.PaginatedMeta, error) {
	var users []models.User
	var total int64

	offset := utils.Pagination(page, pageSize)

	err := s.db.Model(&models.User{}).
		Where("role = ?", models.RoleStaff).
		Count(&total).Offset(offset).
		Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, nil, domain.ErrNotFound
	}

	response := make([]*dto.UserResponse, 0, len(users))

	for _, user := range users {
		response = append(response, s.ConvertToUserResponse(&user))
	}

	totalPages := int(total + int64(pageSize) - 1/int64(pageSize))

	meta := &utils.PaginatedMeta{
		Total:      total,
		Page:       page,
		Limit:      pageSize,
		TotalPages: totalPages,
	}

	return response, meta, nil
}

func (s *UserService) UpdateProfileService(userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	var user models.User
	err := s.db.Where("id = ? AND role = ?", userID, models.RoleCustomer).First(&user).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}
	err = s.db.Save(&user).Error
	if err != nil {
		return nil, err
	}

	return s.ConvertToUserResponse(&user), nil
}

func (s *UserService) UpdateStaffService(userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	var user models.User
	err := s.db.Where("id = ? AND role = ? ", userID, models.RoleStaff).First(&user).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}

	err = s.db.Save(&user).Error
	if err != nil {
		return nil, err
	}

	return s.ConvertToUserResponse(&user), nil
}

// ReUseable Func
func (s *UserService) getUserByRole(userID uint, role models.UserRole) (*dto.UserResponse, error) {
	var user models.User
	err := s.db.Where("id = ? AND is_active = ? AND role = ?", userID, true, role).First(&user).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return s.ConvertToUserResponse(&user), nil
}

func (s *UserService) deactivateUserByRole(userID uint, role models.UserRole) error {
	result := s.db.Model(&models.User{}).
		Where("id = ? AND is_active = ? AND role = ?", userID, true, role).
		Update("is_active", false)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *UserService) activateUserByRole(userID uint, role models.UserRole) error {
	result := s.db.Model(&models.User{}).
		Where("id = ? AND is_active = ? AND role = ?", userID, false, role).
		Update("is_active", true)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
