package services

import (
	"context"
	"testing"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	redisStore "github.com/AboloreDev/geritcht-restaurant/internals/redis"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type MockAllergenRepository struct {
	allergen  *models.Allergen
	allergens []models.Allergen
	total     int64
	menuCount int64
	getErr    error
	createErr error
	updateErr error
	deleteErr error
	countErr  error
}

func (m *MockAllergenRepository) Create(_ context.Context, allergen *models.Allergen) error {
	allergen.ID = 1
	return m.createErr
}
func (m *MockAllergenRepository) GetByName(_ context.Context, name string) (*models.Allergen, error) {
	return m.allergen, m.getErr
}
func (m *MockAllergenRepository) GetByID(_ context.Context, allergenID uint) (*models.Allergen, error) {
	return m.allergen, m.getErr
}
func (m *MockAllergenRepository) GetAll(_ context.Context, page, pageSize int) ([]models.Allergen, int64, error) {
	return m.allergens, m.total, m.getErr
}
func (m *MockAllergenRepository) Update(_ context.Context, allergen *models.Allergen) error {
	return m.updateErr
}
func (m *MockAllergenRepository) Delete(_ context.Context, allergenID uint) error {
	return m.deleteErr
}
func (m *MockAllergenRepository) CountMenuItemsUsingAllergen(_ context.Context, allergenID uint) (int64, error) {
	return m.menuCount, m.countErr
}

func newAllergenService(repo *MockAllergenRepository) *AllergenService {
	return NewAllergenService(repo, redisStore.NewNopCache())
}

func TestCreateAllergen_Success(t *testing.T) {
	service := newAllergenService(&MockAllergenRepository{
		getErr: gorm.ErrRecordNotFound,
	})

	req := &dto.CreateAllergenRequest{Name: "Peanuts"}
	response, err := service.CreateAllergenServices(testTagCtx, req)

	assert.NoError(t, err)
	assert.Equal(t, "Peanuts", response.Name)
}

func TestCreateAllergen_DuplicateName(t *testing.T) {
	service := newAllergenService(&MockAllergenRepository{
		allergen: &models.Allergen{ID: 1, Name: "Peanuts"},
	})

	req := &dto.CreateAllergenRequest{Name: "Peanuts"}
	response, err := service.CreateAllergenServices(testTagCtx, req)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrNameConflict, err)
}

func TestDeleteAllergen_HasLinkedMenuItems(t *testing.T) {
	service := newAllergenService(&MockAllergenRepository{menuCount: 3})

	err := service.DeleteAllergenService(testTagCtx, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete tags")
}

func TestDeleteAllergen_Success(t *testing.T) {
	service := newAllergenService(&MockAllergenRepository{menuCount: 0})

	err := service.DeleteAllergenService(testTagCtx, 1)

	assert.NoError(t, err)
}
