package services

import (
	"context"
	"testing"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	redisStore "github.com/AboloreDev/geritcht-restaurant/internals/redis"
	"github.com/stretchr/testify/assert"
)

var testCategoryCtx = context.Background()

// ─── MockCategoryRepository

type MockCategoryRepository struct {
	category   *models.MenuCategory
	categories []models.MenuCategory
	searchRank []models.MenuCategoryWithRank
	total      int64
	getErr     error
	createErr  error
	updateErr  error
	deleteErr  error
	countErr   error
	menuCount  int64
}

func (m *MockCategoryRepository) Create(_ context.Context, category *models.MenuCategory) error {
	category.ID = 1
	return m.createErr
}
func (m *MockCategoryRepository) GetByID(_ context.Context, categoryID uint) (*models.MenuCategory, error) {
	return m.category, m.getErr
}
func (m *MockCategoryRepository) GetByName(_ context.Context, name string) (*models.MenuCategory, error) {
	return m.category, m.getErr
}
func (m *MockCategoryRepository) GetAll(_ context.Context, page, pageSize int) ([]models.MenuCategory, int64, error) {
	return m.categories, m.total, m.getErr
}
func (m *MockCategoryRepository) Update(_ context.Context, category *models.MenuCategory) error {
	return m.updateErr
}
func (m *MockCategoryRepository) Delete(_ context.Context, categoryID uint) error {
	return m.deleteErr
}
func (m *MockCategoryRepository) CountMenuItems(_ context.Context, categoryID uint) (int64, error) {
	return m.menuCount, m.countErr
}
func (m *MockCategoryRepository) TsvectorSearchCategories(_ context.Context, req *dto.CategorySearchRequest) ([]models.MenuCategoryWithRank, int64, error) {
	return m.searchRank, m.menuCount, m.countErr
}

// ─── Helper

func newCategoryService(repo *MockCategoryRepository) *CategoryService {
	return NewCategoryService(redisStore.NewNopCache(), repo)
}

// ─── CreateCategory Tests

func TestCreateCategory_Success(t *testing.T) {
	service := newCategoryService(&MockCategoryRepository{
		getErr: domain.ErrNotFound, // name not found, proceed to create
	})

	req := &dto.CreateCategoryRequest{
		Name:        "Soups",
		Description: "Nigerian soups",
	}

	response, err := service.CreateCategoryService(testCategoryCtx, req, "image.jpg")

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, req.Name, response.Name)
	assert.True(t, response.IsActive)
}

func TestCreateCategory_DuplicateName(t *testing.T) {
	service := newCategoryService(&MockCategoryRepository{
		category: &models.MenuCategory{ID: 1, Name: "Soups"},
		getErr:   nil, // found, there is a conflict
	})

	req := &dto.CreateCategoryRequest{Name: "Soups"}

	response, err := service.CreateCategoryService(testCategoryCtx, req, "")

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrNameConflict, err)
}

func TestCreateCategory_CreateFails(t *testing.T) {
	service := newCategoryService(&MockCategoryRepository{
		getErr:    domain.ErrNotFound,
		createErr: domain.ErrInternalServerError,
	})

	req := &dto.CreateCategoryRequest{Name: "Soups"}

	response, err := service.CreateCategoryService(testCategoryCtx, req, "")

	assert.Nil(t, response)
	assert.Error(t, err)
}

// ─── UpdateCategory Tests

func TestUpdateCategory(t *testing.T) {
	tests := []struct {
		name        string
		category    *models.MenuCategory
		getErr      error
		updateErr   error
		req         *dto.UpdateCategoryRequest
		expectedErr error
	}{
		{
			name:        "success",
			category:    &models.MenuCategory{ID: 1, Name: "Old Name"},
			req:         &dto.UpdateCategoryRequest{Name: "New Name"},
			expectedErr: nil,
		},
		{
			name:        "category not found",
			getErr:      domain.ErrNotFound,
			req:         &dto.UpdateCategoryRequest{Name: "New Name"},
			expectedErr: domain.ErrNotFound,
		},
		{
			name:        "update fails",
			category:    &models.MenuCategory{ID: 1},
			updateErr:   domain.ErrInternalServerError,
			req:         &dto.UpdateCategoryRequest{Name: "New Name"},
			expectedErr: domain.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newCategoryService(&MockCategoryRepository{
				category:  tt.category,
				getErr:    tt.getErr,
				updateErr: tt.updateErr,
			})

			response, err := service.UpdateCategoryService(testCategoryCtx, 1, tt.req)

			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				assert.NotNil(t, response)
			} else {
				assert.Nil(t, response)
			}
		})
	}
}

// ─── DeleteCategory Tests

func TestDeleteCategory(t *testing.T) {
	tests := []struct {
		name        string
		menuCount   int64
		deleteErr   error
		expectedErr error
	}{
		{
			name:        "success",
			menuCount:   0,
			expectedErr: nil,
		},
		{
			name:        "has linked menu items",
			menuCount:   3,
			expectedErr: nil, // fmt.Errorf will check error separately
		},
		{
			name:        "delete fails",
			menuCount:   0,
			deleteErr:   domain.ErrNotFound,
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newCategoryService(&MockCategoryRepository{
				menuCount: tt.menuCount,
				deleteErr: tt.deleteErr,
			})

			err := service.DeleteCategoryService(testCategoryCtx, 1)

			if tt.name == "has linked menu items" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot delete category")
				return
			}

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

// ─── GetCategory Tests

func TestGetCategory_Success(t *testing.T) {
	service := newCategoryService(&MockCategoryRepository{
		category: &models.MenuCategory{ID: 1, Name: "Soups"},
	})

	response, err := service.GetCategoryService(testCategoryCtx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Soups", response.Name)
}

func TestGetCategory_NotFound(t *testing.T) {
	service := newCategoryService(&MockCategoryRepository{
		getErr: domain.ErrNotFound,
	})

	response, err := service.GetCategoryService(testCategoryCtx, 999)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrNotFound, err)
}

// ─── GetCategories Tests

func TestGetCategories_Success(t *testing.T) {
	service := newCategoryService(&MockCategoryRepository{
		categories: []models.MenuCategory{
			{ID: 1, Name: "Soups"},
			{ID: 2, Name: "Rice"},
		},
		total: 2,
	})

	response, meta, err := service.GetCategoriesService(testCategoryCtx, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, int64(2), meta.Total)
	assert.Equal(t, 1, meta.TotalPages)
	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 10, meta.Limit)
}

func TestGetCategories_Empty(t *testing.T) {
	service := newCategoryService(&MockCategoryRepository{
		categories: []models.MenuCategory{},
		total:      0,
	})

	response, meta, err := service.GetCategoriesService(testCategoryCtx, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, response, 0)
	assert.Equal(t, int64(0), meta.Total)
	assert.Equal(t, 0, meta.TotalPages)
	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 10, meta.Limit)
}
