package services

import (
	"context"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
)

type UserService struct {
	userRepo repositories.UserRepositoryInterface
}

func NewUserService(userRepo repositories.UserRepositoryInterface) *UserService {
	return &UserService{
		userRepo: userRepo,
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

func (s *UserService) GetUserProfileService(ctx context.Context, userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByIDAndRole(ctx, userID, models.RoleCustomer)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return s.ConvertToUserResponse(user), nil
}

func (s *UserService) GetStaffProfileService(ctx context.Context, userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByIDAndRole(ctx, userID, models.RoleStaff)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return s.ConvertToUserResponse(user), nil
}

func (s *UserService) GetAllUsersService(ctx context.Context, page, pageSize int) ([]*dto.UserResponse, *utils.PaginatedMeta, error) {
	users, total, err := s.userRepo.GetAllByRole(ctx, models.RoleCustomer, page, pageSize)
	if err != nil {
		return nil, nil, err
	}
	return s.buildUserListResponse(users, total, page, pageSize)
}

func (s *UserService) DeactivateUserService(ctx context.Context, userID uint) error {
	return s.userRepo.UpdateActiveByRole(ctx, userID, models.RoleCustomer, false)
}

func (s *UserService) DeactivateStaffService(ctx context.Context, userID uint) error {
	return s.userRepo.UpdateActiveByRole(ctx, userID, models.RoleStaff, false)
}

func (s *UserService) ActivateUserService(ctx context.Context, userID uint) error {
	return s.userRepo.UpdateActiveByRole(ctx, userID, models.RoleCustomer, true)
}

func (s *UserService) ActivateStaffService(ctx context.Context, userID uint) error {
	return s.userRepo.UpdateActiveByRole(ctx, userID, models.RoleStaff, true)
}

func (s *UserService) GetAllStaffService(ctx context.Context, page int, pageSize int) ([]*dto.UserResponse, *utils.PaginatedMeta, error) {
	users, total, err := s.userRepo.GetAllByRole(ctx, models.RoleStaff, page, pageSize)
	if err != nil {
		return nil, nil, err
	}
	return s.buildUserListResponse(users, total, page, pageSize)
}

func (s *UserService) UpdateProfileService(ctx context.Context, userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByIDAndRole(ctx, userID, models.RoleCustomer)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return s.updateUser(ctx, user, req)
}

func (s *UserService) UpdateStaffService(ctx context.Context, userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByIDAndRole(ctx, userID, models.RoleStaff)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return s.updateUser(ctx, user, req)
}

// ReUseable helpers
func (s *UserService) updateUser(ctx context.Context, user *models.User, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return s.ConvertToUserResponse(user), nil
}

func (s *UserService) buildUserListResponse(users []*models.User, total int64, page, pageSize int) ([]*dto.UserResponse, *utils.PaginatedMeta, error) {
	response := make([]*dto.UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, s.ConvertToUserResponse(user))
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	meta := &utils.PaginatedMeta{
		Total:      total,
		Page:       page,
		Limit:      pageSize,
		TotalPages: totalPages,
	}

	return response, meta, nil
}
