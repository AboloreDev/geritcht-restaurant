package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
)

type MenuService struct {
	db         *gorm.DB
	redisStore interfaces.Cacher
}

func NewMenuService(db *gorm.DB, redisStore interfaces.Cacher) *MenuService {
	return &MenuService{
		db:         db,
		redisStore: redisStore,
	}
}

func (s *MenuService) ConvertToMenuResponse(menu *models.Menu) *dto.MenuResponse {
	if menu == nil {
		return nil
	}

	// Image Attchement
	menuImages := make([]dto.MenuImageResponse, len(menu.Images))

	for i := range menu.Images {
		menuImages[i] = dto.MenuImageResponse{
			ID:        menu.Images[i].ID,
			URL:       menu.Images[i].URL,
			AltText:   menu.Images[i].AltText,
			IsPrimary: menu.Images[i].IsPrimary,
			CreatedAt: menu.Images[i].CreatedAt,
		}
	}

	// Allergens Attachment
	allergens := make([]dto.AllergenResponse, 0, len(menu.Allergens))

	for _, allergen := range menu.Allergens {
		allergens = append(allergens, dto.AllergenResponse{
			ID:   allergen.ID,
			Name: allergen.Name,
		})
	}

	// Dietary Tags Attachment
	dietaryTags := make([]dto.DietaryTagResponse, 0, len(menu.DietaryTags))

	for _, dietaryTag := range menu.DietaryTags {
		dietaryTags = append(dietaryTags, dto.DietaryTagResponse{
			ID:   dietaryTag.ID,
			Name: dietaryTag.Name,
		})
	}

	return &dto.MenuResponse{
		ID:         menu.ID,
		CategoryID: menu.CategoryID,
		Category: &dto.MenuCategoryResponse{
			ID:          menu.Category.ID,
			Name:        menu.Category.Name,
			Description: menu.Category.Description,
			ImageURL:    menu.Category.ImageURL,
			IsActive:    menu.Category.IsActive,
		},
		Name:            menu.Name,
		Description:     menu.Description,
		Price:           menu.Price,
		ImageURL:        menu.ImageURL,
		IsAvailable:     menu.IsAvailable,
		PrepTimeMinutes: menu.PrepTimeMinutes,
		SpiceLevel:      menu.SpiceLevel,
		Allergens:       allergens,
		DietaryTags:     dietaryTags,
		Images:          menuImages,
		CreatedAt:       menu.CreatedAt,
		UpdatedAt:       menu.UpdatedAt,
	}
}

func (s *MenuService) CreateMenuService(req *dto.CreateMenuRequest) (*dto.MenuResponse, error) {
	var category models.MenuCategory
	var allergens []models.Allergen
	var dietaryTags []models.DietaryTag
	var existingMenu models.Menu
	var count int64

	err := s.db.First(&category, req.CategoryID).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	err = s.db.Find(&allergens, req.AllergenIDs).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	err = s.db.Find(&dietaryTags, req.DietaryTagIDs).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	s.db.Where("name = ? AND category_id = ?", req.Name, req.CategoryID).First(&existingMenu).Count(&count)

	if count > 0 {
		return nil, domain.ErrNameConflict
	}

	menu := models.Menu{
		CategoryID:      category.ID,
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		PrepTimeMinutes: req.PrepTimeMinutes,
		SpiceLevel:      req.SpiceLevel,
		Allergens:       allergens,
		DietaryTags:     dietaryTags,
	}

	err = s.db.Create(&menu).Error
	if err != nil {
		return nil, err
	}

	s.redisStore.Delete(ctx, "menu:all")
	s.redisStore.Delete(ctx, fmt.Sprintf("menu:category:%d", menu.CategoryID))

	return s.ConvertToMenuResponse(&menu), nil
}

func (s *MenuService) UpdateMenuService(menuID uint, req *dto.UpdateMenuRequest) (*dto.MenuResponse, error) {
	var menu models.Menu
	err := s.db.First(&menu, menuID).Error
	if err != nil {
		return nil, err
	}

	var allergens []models.Allergen
	s.db.Find(&allergens, req.AllergenIDs)

	var dietaryTags []models.DietaryTag
	s.db.Find(&dietaryTags, req.DietaryTagIDs)

	menu.CategoryID = req.CategoryID
	menu.Name = req.Name
	menu.Description = req.Description
	menu.Price = req.Price
	menu.PrepTimeMinutes = req.PrepTimeMinutes
	menu.SpiceLevel = req.SpiceLevel
	if req.IsAvailable != nil {
		menu.IsAvailable = *req.IsAvailable
	}
	s.db.Model(&menu).Association("Allergens").Replace(allergens)
	s.db.Model(&menu).Association("DietaryTags").Replace(dietaryTags)

	err = s.db.Save(&menu).Error
	if err != nil {
		return nil, err
	}

	s.redisStore.Delete(ctx, "menu:all")
	s.redisStore.Delete(ctx, fmt.Sprintf("menu:item:%d", menuID))

	return s.ConvertToMenuResponse(&menu), nil
}

func (s *MenuService) GetMenu(menuID uint) (*dto.MenuResponse, error) {
	var menu models.Menu
	cachedKey := fmt.Sprintf("menu:item:%d", menuID)

	exists, _ := s.redisStore.Exists(ctx, cachedKey)
	if exists {
		cache, err := s.redisStore.Get(ctx, cachedKey)
		if err == nil && cache != "" {
			var menu models.Menu
			err := json.Unmarshal([]byte(cache), &menu)
			if err != nil {
				return nil, err
			}
			return s.ConvertToMenuResponse(&menu), nil
		}
	}

	err := s.db.Preload("Images").Preload("Allergens").Preload("Category").
		Preload("Ingredients").Preload("DietaryTags").
		Where("id = ? AND is_available = ? ", menuID, true).First(&menu).Error
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(&menu)
	if err != nil {
		return nil, fmt.Errorf("Failed to set data: %d", err)
	}
	s.redisStore.Set(ctx, cachedKey, string(data), 1*time.Hour)

	return s.ConvertToMenuResponse(&menu), nil
}

func (s *MenuService) DeleteMenu(menuID uint) error {
	var menu models.Menu

	err := s.db.Where("id = ? ", menuID).Delete(&menu).Error
	if err != nil {
		return domain.ErrNotFound
	}

	if err != nil {
		return err
	}
	if s.db.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	s.redisStore.Delete(ctx, "menu:all")
	s.redisStore.Delete(ctx, fmt.Sprintf("menu:item:%d", menuID))

	return nil
}

func (s *MenuService) AddMenuImageService(menuID uint, altText, url string) error {
	var menuImage models.MenuImage
	var count int64

	err := s.db.Model(&models.MenuImage{}).Where("menu_item_id = ?", menuID).Count(&count).Error
	if err != nil {
		return err
	}

	if count >= 4 {
		return errors.New("maximum number of images reached for this menu item")
	}

	menuImage = models.MenuImage{
		MenuItemID: menuID,
		AltText:    altText,
		URL:        url,
		IsPrimary:  count == 0,
		CreatedAt:  time.Now(),
	}

	err = s.db.Create(&menuImage).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *MenuService) RemoveMenuImageService(menuImageID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var menuImage models.MenuImage

		if err := tx.First(&menuImage, menuImageID).Error; err != nil {
			return err
		}

		var totalImages int64
		tx.Model(&models.MenuImage{}).Where("menu_item_id = ?", menuImage.MenuItemID).Count(&totalImages)

		if totalImages <= 1 {
			return errors.New("cannot delete the only image. Please delete the entire menu item or upload a new image first")
		}

		if err := tx.Delete(&menuImage).Error; err != nil {
			return err
		}

		if menuImage.IsPrimary {
			var nextPrimaryImage models.MenuImage
			err := tx.Where("menu_item_id = ? AND id != ?", menuImage.MenuItemID, menuImageID).First(&nextPrimaryImage).Error
			if err == nil {
				tx.Model(&nextPrimaryImage).Update("is_primary", true)
			}
		}

		s.redisStore.Delete(ctx, fmt.Sprintf("menu:item:%d", menuImage.MenuItemID))

		return nil
	})
}

func (s *MenuService) ToggleMenuAvailabilityService(menuID uint, isAvailable *bool) error {
	var menu models.Menu

	err := s.db.Where("id = ? ", menuID).First(&menu).Error
	if err != nil {
		return err
	}

	if isAvailable != nil {
		menu.IsAvailable = *isAvailable
	}

	return s.db.Save(&menu).Error
}

func (s *MenuService) GetAllMenuService(page, pageSize int) ([]*dto.MenuResponse, *utils.PaginatedMeta, error) {
	var menus []models.Menu
	cacheKey := fmt.Sprintf("menu:all:p%d:s%d", page, pageSize)
	var count int64
	offset := utils.Pagination(page, pageSize)

	cached, err := s.redisStore.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cachedResponse struct {
			Data []*dto.MenuResponse  `json:"data"`
			Meta *utils.PaginatedMeta `json:"meta"`
		}
		err := json.Unmarshal([]byte(cached), &cachedResponse)
		if err == nil {
			return cachedResponse.Data, cachedResponse.Meta, nil
		}
	}

	s.db.Model(&models.Menu{}).Where("is_available = ?", true).Count(&count)

	err = s.db.Preload("Images").Preload("Allergens").Preload("Category").
		Preload("Ingredients").Preload("DietaryTags").
		Where("is_available = ?", true).Find(&menus).
		Order("created_at DESC").Offset(offset).Limit(pageSize).Error
	if err != nil {
		return nil, nil, err
	}

	response := make([]*dto.MenuResponse, 0, len(menus))

	for _, menu := range menus {
		response = append(response, s.ConvertToMenuResponse(&menu))
	}

	totalPages := int(count + int64(pageSize) - 1/int64(pageSize))

	meta := &utils.PaginatedMeta{
		Page:       page,
		Limit:      pageSize,
		Total:      count,
		TotalPages: totalPages,
	}

	cacheData := struct {
		Data []*dto.MenuResponse  `json:"data"`
		Meta *utils.PaginatedMeta `json:"meta"`
	}{Data: response, Meta: meta}

	data, err := json.Marshal(&cacheData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal data: %v", err)
	}

	s.redisStore.Set(ctx, cacheKey, string(data), 1*time.Hour)

	return response, meta, nil
}
