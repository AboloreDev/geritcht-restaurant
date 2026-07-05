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

var testTagCtx = context.Background()

// ─── MockDietaryTagRepository

type MockDietaryTagRepository struct {
	tag       *models.DietaryTag
	tags      []models.DietaryTag
	total     int64
	menuCount int64
	getErr    error
	createErr error
	updateErr error
	deleteErr error
	countErr  error
}

func (m *MockDietaryTagRepository) Create(_ context.Context, tag *models.DietaryTag) error {
	tag.ID = 1
	return m.createErr
}
func (m *MockDietaryTagRepository) GetByName(_ context.Context, name string) (*models.DietaryTag, error) {
	return m.tag, m.getErr
}
func (m *MockDietaryTagRepository) GetByID(_ context.Context, tagID uint) (*models.DietaryTag, error) {
	return m.tag, m.getErr
}
func (m *MockDietaryTagRepository) GetAll(_ context.Context, page, pageSize int) ([]models.DietaryTag, int64, error) {
	return m.tags, m.total, m.getErr
}
func (m *MockDietaryTagRepository) Update(_ context.Context, tag *models.DietaryTag) error {
	return m.updateErr
}
func (m *MockDietaryTagRepository) Delete(_ context.Context, tagID uint) error {
	return m.deleteErr
}
func (m *MockDietaryTagRepository) CountMenuItemsUsingTag(_ context.Context, tagID uint) (int64, error) {
	return m.menuCount, m.countErr
}

func newDietaryTagsService(repo *MockDietaryTagRepository) *DietaryTagsService {
	return NewDietaryTagsService(repo, redisStore.NewNopCache())
}

// ─── Tests

func TestCreateDietaryTag_Success(t *testing.T) {
	service := newDietaryTagsService(&MockDietaryTagRepository{
		getErr: gorm.ErrRecordNotFound,
	})

	req := &dto.CreateDietaryTagRequest{Name: "Vegan"}
	response, err := service.CreateDietaryTagService(testTagCtx, req)

	assert.NoError(t, err)
	assert.Equal(t, "Vegan", response.Name)
}

func TestCreateDietaryTag_DuplicateName(t *testing.T) {
	service := newDietaryTagsService(&MockDietaryTagRepository{
		tag: &models.DietaryTag{ID: 1, Name: "Vegan"},
	})

	req := &dto.CreateDietaryTagRequest{Name: "Vegan"}
	response, err := service.CreateDietaryTagService(testTagCtx, req)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrNameConflict, err)
}

func TestDeleteDietaryTag_HasLinkedMenuItems(t *testing.T) {
	service := newDietaryTagsService(&MockDietaryTagRepository{
		menuCount: 5,
	})

	err := service.DeleteDietaryTagService(testTagCtx, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete tags")
}

func TestDeleteDietaryTag_Success(t *testing.T) {
	service := newDietaryTagsService(&MockDietaryTagRepository{
		menuCount: 0,
	})

	err := service.DeleteDietaryTagService(testTagCtx, 1)

	assert.NoError(t, err)
}

func TestUpdateDietaryTag_NotFound(t *testing.T) {
	service := newDietaryTagsService(&MockDietaryTagRepository{
		getErr: gorm.ErrRecordNotFound,
	})

	req := &dto.UpdateDietaryTagRequest{Name: "New"}
	response, err := service.UpdateDietaryTagService(testTagCtx, 999, req)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrNotFound, err)
}
