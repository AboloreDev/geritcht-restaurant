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
		CategoryID: menu.MenuCategoryID,
		Category: &dto.MenuCategoryResponse{
			ID:           menu.MenuCategory.ID,
			Name:         menu.MenuCategory.Name,
			Description:  menu.MenuCategory.Description,
			ImageURL:     menu.MenuCategory.ImageURL,
			IsActive:     menu.MenuCategory.IsActive,
			DisplayOrder: menu.MenuCategory.DisplayOrder,
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
		DisplayOrder:    menu.DisplayOrder,
		Images:          menuImages,
		CreatedAt:       menu.CreatedAt,
		UpdatedAt:       menu.UpdatedAt,
	}
}

func (s *MenuService) CreateMenuService(req *dto.CreateMenuRequest) (*dto.MenuResponse, error) {
	var category models.MenuCategory
	var allergens []models.Allergen
	var dietaryTags []models.DietaryTag
	var count int64

	// Check category exists
	err := s.db.First(&category, req.CategoryID).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	// Fetch allergens
	err = s.db.Find(&allergens, req.AllergenIDs).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	// Fetch dietary tags
	err = s.db.Find(&dietaryTags, req.DietaryTagIDs).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	// Check duplicate menu
	err = s.db.Model(&models.Menu{}).
		Where("name = ? AND menu_category_id = ?", req.Name, req.CategoryID).
		Count(&count).Error

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, domain.ErrNameConflict
	}

	menu := models.Menu{
		MenuCategoryID:  category.ID,
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		PrepTimeMinutes: req.PrepTimeMinutes,
		SpiceLevel:      req.SpiceLevel,
		Allergens:       allergens,
		DietaryTags:     dietaryTags,
		DisplayOrder:    req.DisplayOrder,
	}

	// Create menu
	err = s.db.Create(&menu).Error
	if err != nil {
		return nil, err
	}

	s.redisStore.FlushByPattern(ctx, "menu:all:*")
	s.redisStore.Delete(ctx, fmt.Sprintf("menu:category:%d", menu.MenuCategoryID))

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

	if req.CategoryID != 0 {
		menu.MenuCategoryID = req.CategoryID
	}

	if req.Name != "" {
		menu.Name = req.Name
	}

	if req.Description != "" {
		menu.Description = req.Description
	}

	if req.Price != 0 {
		menu.Price = req.Price
	}

	if req.PrepTimeMinutes != 0 {
		menu.PrepTimeMinutes = req.PrepTimeMinutes
	}

	if req.SpiceLevel != 0 {
		menu.SpiceLevel = req.SpiceLevel
	}

	if req.DisplayOrder != 0 {
		menu.DisplayOrder = req.DisplayOrder
	}

	if req.IsAvailable != nil {
		menu.IsAvailable = *req.IsAvailable
	}
	s.db.Model(&menu).Association("Allergens").Replace(allergens)
	s.db.Model(&menu).Association("DietaryTags").Replace(dietaryTags)

	err = s.db.Save(&menu).Error
	if err != nil {
		return nil, err
	}

	s.redisStore.FlushByPattern(ctx, "menu:all:*")
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

	err := s.db.Preload("Images").Preload("Allergens").Preload("MenuCategory").
		Preload("MenuItemIngredients").Preload("DietaryTags").
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
	result := s.db.Where("id = ?", menuID).Delete(&models.Menu{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	s.redisStore.FlushByPattern(ctx, "menu:all:*")
	s.redisStore.Delete(ctx, fmt.Sprintf("menu:item:%d", menuID))

	return nil
}

func (s *MenuService) AddMenuImageService(menuID uint, altText, url string) error {
	var menuImage models.MenuImage
	var count int64

	err := s.db.Model(&models.MenuImage{}).Where("menu_id = ?", menuID).Count(&count).Error
	if err != nil {
		return err
	}

	if count >= 4 {
		return errors.New("maximum number of images reached for this menu item")
	}

	menuImage = models.MenuImage{
		MenuID:    menuID,
		AltText:   altText,
		URL:       url,
		IsPrimary: count == 0,
		CreatedAt: time.Now(),
	}

	err = s.db.Create(&menuImage).Error
	if err != nil {
		return err
	}

	s.redisStore.FlushByPattern(ctx, "menu:all:*")
	s.redisStore.Delete(ctx, fmt.Sprintf("menu:item:%d", menuID))

	return nil
}

func (s *MenuService) RemoveMenuImageService(menuImageID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var menuImage models.MenuImage

		if err := tx.First(&menuImage, menuImageID).Error; err != nil {
			return err
		}

		var totalImages int64
		tx.Model(&models.MenuImage{}).Where("menu_id = ?", menuImage.MenuID).Count(&totalImages)

		if totalImages <= 1 {
			return errors.New("cannot delete the only image. Please delete the entire menu item or upload a new image first")
		}

		if err := tx.Delete(&menuImage).Error; err != nil {
			return err
		}

		if menuImage.IsPrimary {
			var nextPrimaryImage models.MenuImage
			err := tx.Where("menu_id = ? AND id != ?", menuImage.MenuID, menuImageID).First(&nextPrimaryImage).Error
			if err == nil {
				tx.Model(&nextPrimaryImage).Update("is_primary", true)
			}
		}

		s.redisStore.FlushByPattern(ctx, "menu:all:*")
		s.redisStore.Delete(ctx, fmt.Sprintf("menu:item:%d", menuImage.MenuID))

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

	err = s.db.Save(&menu).Error
	if err != nil {
		return err
	}

	s.redisStore.FlushByPattern(ctx, "menu:all:*")

	return nil
}

func (s *MenuService) GetAllMenuService(filter dto.MenuFilterRequest) ([]*dto.MenuResponse, *utils.PaginatedMeta, error) {
	cacheKey := utils.BuildMenuCacheKey(filter)
	cached, err := s.redisStore.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cachedResponse struct {
			Data []*dto.MenuResponse  `json:"data"`
			Meta *utils.PaginatedMeta `json:"meta"`
		}
		if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
			return cachedResponse.Data, cachedResponse.Meta, nil
		}
	}

	offset := utils.Pagination(filter.Page, filter.PageSize)

	query := s.db.Model(&models.Menu{}).Where("is_available = ?", true)

	query = utils.ApplyMenuFilters(query, filter)

	var count int64
	query.Count(&count)

	query = utils.ApplyMenuSorting(query, filter)

	var menus []models.Menu
	err = query.
		Preload("Images").
		Preload("Allergens").
		Preload("MenuCategory").
		Preload("DietaryTags").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&menus).Error
	if err != nil {
		return nil, nil, err
	}

	response := make([]*dto.MenuResponse, 0, len(menus))
	for _, menu := range menus {
		response = append(response, s.ConvertToMenuResponse(&menu))
	}

	totalPages := int((count + int64(filter.PageSize) - 1) / int64(filter.PageSize))
	meta := &utils.PaginatedMeta{
		Page:       filter.Page,
		Limit:      filter.PageSize,
		Total:      count,
		TotalPages: totalPages,
	}

	cacheData := struct {
		Data []*dto.MenuResponse  `json:"data"`
		Meta *utils.PaginatedMeta `json:"meta"`
	}{Data: response, Meta: meta}

	data, _ := json.Marshal(&cacheData)
	s.redisStore.Set(ctx, cacheKey, string(data), utils.GetCacheTTL(filter))

	return response, meta, nil
}
