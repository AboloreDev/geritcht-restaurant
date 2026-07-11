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

var testIngredientCtx = context.Background()

type MockIngredientRepository struct {
	ingredients  []models.Ingredient
	ingredient   *models.Ingredient
	admin        *models.User
	rowsAffected int64
	lowStock     []models.Ingredient
	count        int64
	serachRank   []models.IngredientWithRank

	ingredientsErr error
	ingredientErr  error
	adminErr       error
	updateErr      error
	outboxErr      error
	lowStockErr    error
	createErr      error
	countErr       error
}

func (r *MockIngredientRepository) GetIngredientByName(_ context.Context, name string) (*models.Ingredient, error) {
	return r.ingredient, r.ingredientErr
}

func (r *MockIngredientRepository) CreateIngredient(_ context.Context, ingredient *models.Ingredient) error {
	ingredient.ID = 1
	return r.createErr
}

func (r *MockIngredientRepository) GetIngredientByID(_ context.Context, ingredientID uint) (*models.Ingredient, error) {
	return r.ingredient, r.ingredientErr
}

func (r *MockIngredientRepository) UpdateIngredient(_ context.Context, ingredient *models.Ingredient) error {
	return r.updateErr
}

func (r *MockIngredientRepository) IngredientCount(ctx context.Context, ingredientID uint) (int64, error) {
	return r.count, r.countErr
}

func (r *MockIngredientRepository) DeleteIngredient(ctx context.Context, ingredientID uint) error {
	return r.ingredientErr
}

func (r *MockIngredientRepository) GetAllIngredients(ctx context.Context, page, pageSize int) ([]models.Ingredient, int64, error) {
	return r.ingredients, r.count, r.ingredientsErr
}

func (r *MockIngredientRepository) CompareCurrentStockAgainstMinTheshold(ctx context.Context) ([]models.Ingredient, error) {
	return r.ingredients, r.ingredientsErr
}

func (r *MockIngredientRepository) UpdateThreshHoldLimit(ctx context.Context, ingredientID uint, threshHold float64) error {
	return r.ingredientErr
}
func (r *MockIngredientRepository) TsvectorSearchIngredients(ctx context.Context, req *dto.IngredientSearchRequest) ([]models.IngredientWithRank, int64, error) {
	return r.serachRank, r.count, r.ingredientsErr
}

func newIngredientService(repo *MockIngredientRepository) *IngredientService {
	return NewIngredientService(
		redisStore.NewNopCache(),
		&MockPublisher{},
		repo,
		&MockUserRepository{},
		&MockPaymentRepository{},
	)
}

func TestCreateIngredient_Success(t *testing.T) {
	service := newIngredientService(&MockIngredientRepository{
		ingredientErr: domain.ErrIngredientNotFound, // name not found, proceed to create
	})

	req := &dto.CreateIngredientRequest{
		Name:         "Soups",
		Unit:         "kg",
		CurrentStock: 20.0,
		MinThreshold: 5.0,
	}

	response, err := service.CreateIngredientService(testIngredientCtx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, req.Name, response.Name)
	assert.Equal(t, req.Unit, response.Unit)
	assert.Equal(t, req.CurrentStock, response.CurrentStock)
	assert.Equal(t, req.MinThreshold, response.MinThreshold)
	assert.Equal(t, req.CurrentStock <= req.MinThreshold, response.IsLow)
}

func TestCreateIngredient_DuplicateName(t *testing.T) {
	service := newIngredientService(&MockIngredientRepository{
		ingredient: &models.Ingredient{ID: 1, Name: "Soups"},
		createErr:  nil, // found, there is a conflict
	})

	req := &dto.CreateIngredientRequest{
		Name:         "Soups",
		Unit:         "kg",
		CurrentStock: 20.0,
		MinThreshold: 5.0,
	}

	response, err := service.CreateIngredientService(testCategoryCtx, req)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrNameConflict, err)
}

func TestCreateIngredient_CreateFails(t *testing.T) {
	service := newIngredientService(&MockIngredientRepository{
		createErr: domain.ErrInternalServerError,
	})

	req := &dto.CreateIngredientRequest{
		Name:         "Soups",
		Unit:         "kg",
		CurrentStock: 20.0,
		MinThreshold: 5.0,
	}

	response, err := service.CreateIngredientService(testCategoryCtx, req)

	assert.Nil(t, response)
	assert.Error(t, err)
}

func TestUpdateIngredients(t *testing.T) {
	tests := []struct {
		name        string
		ingredient  *models.Ingredient
		getErr      error
		updateErr   error
		req         *dto.UpdateIngredientRequest
		expectedErr error
	}{
		{
			name: "success",
			ingredient: &models.Ingredient{
				ID:           1,
				Name:         "Old Name",
				Unit:         "kg",
				CurrentStock: 20.0,
				MinThreshold: 5.0,
			},
			req: &dto.UpdateIngredientRequest{
				Name:         "New Name",
				Unit:         "litre",
				MinThreshold: 6.0,
			},
			expectedErr: nil,
		},
		{
			name:        "ingredient not found",
			getErr:      domain.ErrIngredientNotFound,
			req:         &dto.UpdateIngredientRequest{Name: "New Name"},
			expectedErr: domain.ErrIngredientNotFound,
		},
		{
			name:        "update fails",
			ingredient:  &models.Ingredient{ID: 1},
			updateErr:   domain.ErrInternalServerError,
			req:         &dto.UpdateIngredientRequest{Name: "New Name"},
			expectedErr: domain.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newIngredientService(&MockIngredientRepository{
				ingredient:    tt.ingredient,
				ingredientErr: tt.getErr,
				updateErr:     tt.updateErr,
			})

			response, err := service.UpdateIngredientService(testCategoryCtx, 1, tt.req)

			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				assert.NotNil(t, response)
			} else {
				assert.Nil(t, response)
			}
		})
	}
}

func TestGetIngredient_Success(t *testing.T) {
	service := newIngredientService(&MockIngredientRepository{
		ingredient: &models.Ingredient{
			ID:           1,
			Name:         "Soups",
			Unit:         "kg",
			CurrentStock: 20.0,
			MinThreshold: 5.0,
		},
	})

	response, err := service.GetIngredientService(testCategoryCtx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Soups", response.Name)
}

func TestGetIngredient_NotFound(t *testing.T) {
	service := newIngredientService(&MockIngredientRepository{
		ingredientErr: domain.ErrIngredientNotFound,
	})

	response, err := service.GetIngredientService(testCategoryCtx, 999)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrIngredientNotFound, err)
}

// ─── GetCategories Tests

func TestGetIngredients_Success(t *testing.T) {
	service := newIngredientService(&MockIngredientRepository{
		ingredients: []models.Ingredient{
			{
				ID:           1,
				Name:         "Soups",
				Unit:         "kg",
				CurrentStock: 20.0,
				MinThreshold: 5.0,
			},
			{
				ID:           2,
				Name:         "Rice",
				Unit:         "kg",
				CurrentStock: 30.0,
				MinThreshold: 10.0,
			},
		},
		count: 2,
	})

	response, meta, err := service.GetAllIngredientService(testCategoryCtx, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, int64(2), meta.Total)
	assert.Equal(t, 1, meta.TotalPages)
	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 10, meta.Limit)
}

func TestGetIngredients_Empty(t *testing.T) {
	service := newIngredientService(&MockIngredientRepository{
		ingredients: []models.Ingredient{},
		count:       0,
	})

	response, meta, err := service.GetAllIngredientService(testCategoryCtx, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, response, 0)
	assert.Equal(t, int64(0), meta.Total)
	assert.Equal(t, 0, meta.TotalPages)
	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 10, meta.Limit)
}
