package services

import (
	"context"
	"testing"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	redisStore "github.com/AboloreDev/geritcht-restaurant/internals/redis"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testMenuCtx = context.Background()

// ─── MockMenuRepository

type MockMenuRepository struct {
	category    *models.MenuCategory
	allergens   []models.Allergen
	dietaryTags []models.DietaryTag
	searchRank []models.MenuWithRank
	menu        *models.Menu
	menus       []models.Menu
	total       int64
	nameCount   int64
	image       *models.MenuImage
	imageCount  int64
	nextImage   *models.MenuImage

	categoryErr error
	getErr      error
	createErr   error
	updateErr   error
	deleteErr   error
	countErr    error
	imageErr    error
}

func (m *MockMenuRepository) GetCategoryByID(_ context.Context, categoryID uint) (*models.MenuCategory, error) {
	return m.category, m.categoryErr
}
func (m *MockMenuRepository) GetAllergensByIDs(_ context.Context, ids []uint) ([]models.Allergen, error) {
	return m.allergens, m.getErr
}
func (m *MockMenuRepository) GetDietaryTagsByIDs(_ context.Context, ids []uint) ([]models.DietaryTag, error) {
	return m.dietaryTags, m.getErr
}
func (m *MockMenuRepository) CountByNameAndCategory(_ context.Context, name string, categoryID uint) (int64, error) {
	return m.nameCount, m.countErr
}
func (m *MockMenuRepository) Create(_ context.Context, menu *models.Menu) error {
	menu.ID = 1
	return m.createErr
}
func (m *MockMenuRepository) GetByID(_ context.Context, menuID uint) (*models.Menu, error) {
	return m.menu, m.getErr
}
func (m *MockMenuRepository) GetByIDAvailable(_ context.Context, menuID uint) (*models.Menu, error) {
	return m.menu, m.getErr
}
func (m *MockMenuRepository) Update(_ context.Context, menu *models.Menu) error { return m.updateErr }
func (m *MockMenuRepository) ReplaceAllergens(_ context.Context, menu *models.Menu, allergens []models.Allergen) error {
	return nil
}
func (m *MockMenuRepository) ReplaceDietaryTags(_ context.Context, menu *models.Menu, tags []models.DietaryTag) error {
	return nil
}
func (m *MockMenuRepository) Delete(_ context.Context, menuID uint) error { return m.deleteErr }
func (m *MockMenuRepository) GetAll(_ context.Context, filter dto.MenuFilterRequest) ([]models.Menu, int64, error) {
	return m.menus, m.total, m.getErr
}
func (m *MockMenuRepository) CountImages(_ context.Context, menuID uint) (int64, error) {
	return m.imageCount, m.countErr
}
func (m *MockMenuRepository) CreateImage(_ context.Context, image *models.MenuImage) error {
	return m.imageErr
}
func (m *MockMenuRepository) GetImageByID(_ context.Context, imageID uint) (*models.MenuImage, error) {
	return m.image, m.imageErr
}
func (m *MockMenuRepository) DeleteImage(_ context.Context, image *models.MenuImage) error {
	return m.imageErr
}
func (m *MockMenuRepository) GetNextPrimaryImage(_ context.Context, menuID uint, excludeID uint) (*models.MenuImage, error) {
	return m.nextImage, m.imageErr
}
func (m *MockMenuRepository) SetImagePrimary(_ context.Context, image *models.MenuImage) error {
	return nil
}

func (m *MockMenuRepository) TsvectorSearchMenuItems(ctx context.Context, req *dto.MenuSearchRequest) ([]models.MenuWithRank, int64, error) {
	return m.searchRank, m.total, m.getErr
}

func newMenuService(repo *MockMenuRepository) *MenuService {
	return NewMenuService(repo, redisStore.NewNopCache(), zerolog.Nop().With().Logger())
}

// ─── CreateMenu Tests

func TestCreateMenu_Success(t *testing.T) {
	service := newMenuService(&MockMenuRepository{
		category:  &models.MenuCategory{ID: 1, Name: "Soups"},
		nameCount: 0,
	})

	req := &dto.CreateMenuRequest{
		CategoryID: 1,
		Name:       "Jollof Rice",
		Price:      3500,
	}

	response, err := service.CreateMenuService(testMenuCtx, req)

	assert.NoError(t, err)
	assert.Equal(t, "Jollof Rice", response.Name)
}

func TestCreateMenu_CategoryNotFound(t *testing.T) {
	service := newMenuService(&MockMenuRepository{
		categoryErr: domain.ErrNotFound,
	})

	req := &dto.CreateMenuRequest{CategoryID: 999, Name: "Jollof Rice"}

	response, err := service.CreateMenuService(testMenuCtx, req)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestCreateMenu_DuplicateName(t *testing.T) {
	service := newMenuService(&MockMenuRepository{
		category:  &models.MenuCategory{ID: 1},
		nameCount: 1, // already exists
	})

	req := &dto.CreateMenuRequest{CategoryID: 1, Name: "Jollof Rice"}

	response, err := service.CreateMenuService(testMenuCtx, req)

	assert.Nil(t, response)
	assert.Equal(t, domain.ErrNameConflict, err)
}

// ─── GetMenu Tests

func TestGetMenu_Success(t *testing.T) {
	service := newMenuService(&MockMenuRepository{
		menu: &models.Menu{ID: 1, Name: "Jollof Rice", IsAvailable: true},
	})

	response, err := service.GetMenu(testMenuCtx, 1)

	assert.NoError(t, err)
	assert.Equal(t, "Jollof Rice", response.Name)
}

func TestGetMenu_NotFound(t *testing.T) {
	service := newMenuService(&MockMenuRepository{
		getErr: domain.ErrNotFound,
	})

	response, err := service.GetMenu(testMenuCtx, 999)

	assert.Nil(t, response)
	assert.Error(t, err)
}

// ─── DeleteMenu Tests

func TestDeleteMenu_Success(t *testing.T) {
	service := newMenuService(&MockMenuRepository{})

	err := service.DeleteMenu(testMenuCtx, 1)

	assert.NoError(t, err)
}

func TestDeleteMenu_NotFound(t *testing.T) {
	service := newMenuService(&MockMenuRepository{
		deleteErr: gorm.ErrRecordNotFound,
	})

	err := service.DeleteMenu(testMenuCtx, 999)

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

// ─── AddMenuImage Tests

func TestAddMenuImage_Success(t *testing.T) {
	service := newMenuService(&MockMenuRepository{
		imageCount: 2,
	})

	err := service.AddMenuImageService(testMenuCtx, 1, "alt text", "image.jpg")

	assert.NoError(t, err)
}

func TestAddMenuImage_MaxReached(t *testing.T) {
	service := newMenuService(&MockMenuRepository{
		imageCount: 4, // already at max
	})

	err := service.AddMenuImageService(testMenuCtx, 1, "alt", "img.jpg")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maximum number of images")
}

// ─── RemoveMenuImage Tests

func TestRemoveMenuImage_Success(t *testing.T) {
	service := newMenuService(&MockMenuRepository{
		image:      &models.MenuImage{ID: 1, MenuID: 1, IsPrimary: false},
		imageCount: 3,
	})

	err := service.RemoveMenuImageService(testMenuCtx, 1)

	assert.NoError(t, err)
}

func TestRemoveMenuImage_LastImage(t *testing.T) {
	service := newMenuService(&MockMenuRepository{
		image:      &models.MenuImage{ID: 1, MenuID: 1},
		imageCount: 1, // only one image
	})

	err := service.RemoveMenuImageService(testMenuCtx, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete the only image")
}

func TestRemoveMenuImage_PrimaryReassigned(t *testing.T) {
	service := newMenuService(&MockMenuRepository{
		image:      &models.MenuImage{ID: 1, MenuID: 1, IsPrimary: true},
		imageCount: 2,
		nextImage:  &models.MenuImage{ID: 2, MenuID: 1},
	})

	err := service.RemoveMenuImageService(testMenuCtx, 1)

	assert.NoError(t, err)
}

// ─── ToggleAvailability Tests

func TestToggleMenuAvailability_Success(t *testing.T) {
	service := newMenuService(&MockMenuRepository{
		menu: &models.Menu{ID: 1, IsAvailable: true},
	})

	available := false
	err := service.ToggleMenuAvailabilityService(testMenuCtx, 1, &available)

	assert.NoError(t, err)
}

// ─── GetAllMenu Tests

func TestGetAllMenu_Success(t *testing.T) {
	service := newMenuService(&MockMenuRepository{
		menus: []models.Menu{
			{ID: 1, Name: "Jollof Rice"},
			{ID: 2, Name: "Fried Rice"},
		},
		total: 2,
	})

	filter := dto.MenuFilterRequest{Page: 1, PageSize: 10}
	response, meta, err := service.GetAllMenuService(testMenuCtx, filter)

	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, int64(2), meta.Total)
}
