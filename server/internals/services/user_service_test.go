package services

import (
	"context"
	"testing"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/stretchr/testify/assert"
)

var testUserCtx = context.Background()

// ─── MockUserRepository extended from auth service test.go

// ─── Helper
func newUserService(repo *MockUserRepository) *UserService {
	return NewUserService(repo)
}

// ─── GetUserProfile Tests

func TestGetUserProfile_Success(t *testing.T) {
	service := newUserService(&MockUserRepository{
		user: &models.User{
			ID:    1,
			Email: "test@test.com",
			Role:  models.RoleCustomer,
		},
	})

	response, err := service.GetUserProfileService(testUserCtx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "test@test.com", response.Email)
}

func TestGetUserProfile_NotFound(t *testing.T) {
	service := newUserService(&MockUserRepository{
		getErr: domain.ErrNotFound,
	})

	response, err := service.GetUserProfileService(testUserCtx, 999)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrNotFound, err)
}

// ─── GetStaffProfile Tests

func TestGetStaffProfile_Success(t *testing.T) {
	service := newUserService(&MockUserRepository{
		user: &models.User{
			ID:   1,
			Role: models.RoleStaff,
		},
	})

	response, err := service.GetStaffProfileService(testUserCtx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, response)
}

// ─── DeactivateUser Tests

func TestDeactivateUser(t *testing.T) {
	tests := []struct {
		name        string
		updateErr   error
		expectedErr error
	}{
		{name: "success", updateErr: nil, expectedErr: nil},
		{name: "user not found", updateErr: domain.ErrNotFound, expectedErr: domain.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newUserService(&MockUserRepository{
				updateErr: tt.updateErr,
			})

			err := service.DeactivateUserService(testUserCtx, 1)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestActivateUser_Success(t *testing.T) {
	service := newUserService(&MockUserRepository{})

	err := service.ActivateUserService(testUserCtx, 1)

	assert.NoError(t, err)
}

// ─── UpdateProfile Tests

func TestUpdateProfile_Success(t *testing.T) {
	service := newUserService(&MockUserRepository{
		user: &models.User{
			ID:        1,
			FirstName: "Old",
			LastName:  "Name",
			Role:      models.RoleCustomer,
		},
	})

	req := &dto.UpdateProfileRequest{
		FirstName: "New",
		LastName:  "Name",
	}

	response, err := service.UpdateProfileService(testUserCtx, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, "New", response.FirstName)
	assert.Equal(t, "Name", response.LastName)
}

func TestUpdateProfile_NotFound(t *testing.T) {
	service := newUserService(&MockUserRepository{
		getErr: domain.ErrNotFound,
	})

	req := &dto.UpdateProfileRequest{FirstName: "New"}

	response, err := service.UpdateProfileService(testUserCtx, 999, req)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestUpdateProfile_PartialUpdate(t *testing.T) {
	// only firstName provided — lastName and phone should stay unchanged
	service := newUserService(&MockUserRepository{
		user: &models.User{
			ID:          1,
			FirstName:   "Old",
			LastName:    "Original",
			PhoneNumber: "08000000000",
			Role:        models.RoleCustomer,
		},
	})

	req := &dto.UpdateProfileRequest{
		FirstName: "Updated",
		// LastName and PhoneNumber empty
	}

	response, err := service.UpdateProfileService(testUserCtx, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, "Updated", response.FirstName)
	assert.Equal(t, "Original", response.LastName)       // unchanged
	assert.Equal(t, "08000000000", response.PhoneNumber) // unchanged
}

// ─── GetAllUsers Tests

func TestGetAllUsers_Success(t *testing.T) {
	service := newUserService(&MockUserRepository{
		users: []*models.User{
			{ID: 1, Email: "user1@test.com", Role: models.RoleCustomer},
			{ID: 2, Email: "user2@test.com", Role: models.RoleCustomer},
		},
		total: 2,
	})

	response, meta, err := service.GetAllUsersService(testUserCtx, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, int64(2), meta.Total)
}

func TestGetAllStaff_Success(t *testing.T) {
	service := newUserService(&MockUserRepository{
		users: []*models.User{
			{ID: 1, Email: "staff1@test.com", Role: models.RoleStaff},
		},
		total: 1,
	})

	response, meta, err := service.GetAllStaffService(testUserCtx, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, int64(1), meta.Total)
}
