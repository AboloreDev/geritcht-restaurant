package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
)

type MenuService struct {
	menuRepo   repositories.MenuRepositoryInterface
	redisStore interfaces.Cacher
}

func NewMenuService(menuRepo repositories.MenuRepositoryInterface, redisStore interfaces.Cacher) *MenuService {
	return &MenuService{
		menuRepo:   menuRepo,
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

func (s *MenuService) CreateMenuService(ctx context.Context, req *dto.CreateMenuRequest) (*dto.MenuResponse, error) {
	category, err := s.menuRepo.GetCategoryByID(ctx, req.CategoryID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	allergens, err := s.menuRepo.GetAllergensByIDs(ctx, req.AllergenIDs)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	dietaryTags, err := s.menuRepo.GetDietaryTagsByIDs(ctx, req.DietaryTagIDs)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	// Check duplicate menu
	count, err := s.menuRepo.CountByNameAndCategory(ctx, req.Name, req.CategoryID)
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
	err = s.menuRepo.Create(ctx, &menu)
	if err != nil {
		return nil, err
	}

	s.redisStore.FlushByPattern(ctx, "menu:all:*")
	s.redisStore.Delete(ctx, fmt.Sprintf("menu:category:%d", menu.MenuCategoryID))

	return s.ConvertToMenuResponse(&menu), nil
}

func (s *MenuService) UpdateMenuService(ctx context.Context, menuID uint, req *dto.UpdateMenuRequest) (*dto.MenuResponse, error) {
	menu, err := s.menuRepo.GetByID(ctx, menuID)
	if err != nil {
		return nil, err
	}

	allergen, err := s.menuRepo.GetAllergensByIDs(ctx, req.AllergenIDs)
	dietaryTags, err := s.menuRepo.GetDietaryTagsByIDs(ctx, req.DietaryTagIDs)

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

	s.menuRepo.ReplaceAllergens(ctx, menu, allergen)
	s.menuRepo.ReplaceDietaryTags(ctx, menu, dietaryTags)

	err = s.menuRepo.Update(ctx, menu)
	if err != nil {
		return nil, err
	}

	s.redisStore.FlushByPattern(ctx, "menu:all:*")
	s.redisStore.Delete(ctx, "menu:all")
	s.redisStore.Delete(ctx, fmt.Sprintf("menu:item:%d", menuID))

	return s.ConvertToMenuResponse(menu), nil
}

func (s *MenuService) GetMenu(ctx context.Context, menuID uint) (*dto.MenuResponse, error) {
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

	menu, err := s.menuRepo.GetByIDAvailable(ctx, menuID)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(&menu)
	if err != nil {
		return nil, fmt.Errorf("Failed to set data: %d", err)
	}
	s.redisStore.Set(ctx, cachedKey, string(data), 1*time.Hour)

	return s.ConvertToMenuResponse(menu), nil
}

func (s *MenuService) DeleteMenu(ctx context.Context, menuID uint) error {
	err := s.menuRepo.Delete(ctx, menuID)
	if err != nil {
		return err
	}

	s.redisStore.FlushByPattern(ctx, "menu:all:*")
	s.redisStore.Delete(ctx, fmt.Sprintf("menu:item:%d", menuID))

	return nil
}

func (s *MenuService) AddMenuImageService(ctx context.Context, menuID uint, altText, url string) error {

	count, err := s.menuRepo.CountImages(ctx, menuID)
	if err != nil {
		return err
	}

	if count >= 4 {
		return errors.New("maximum number of images reached for this menu item")
	}

	menuImage := models.MenuImage{
		MenuID:    menuID,
		AltText:   altText,
		URL:       url,
		IsPrimary: count == 0,
		CreatedAt: time.Now(),
	}

	err = s.menuRepo.CreateImage(ctx, &menuImage)
	if err != nil {
		return err
	}

	s.redisStore.FlushByPattern(ctx, "menu:all:*")
	s.redisStore.Delete(ctx, fmt.Sprintf("menu:item:%d", menuID))

	return nil
}

func (s *MenuService) RemoveMenuImageService(ctx context.Context, menuImageID uint) error {
	image, err := s.menuRepo.GetImageByID(ctx, menuImageID)
	if err != nil {
		return err
	}

	totalImages, err := s.menuRepo.CountImages(ctx, image.MenuID)
	if err != nil {
		return err
	}
	if totalImages <= 1 {
		return errors.New("cannot delete the only image. Please delete the entire menu item or upload a new image first")
	}

	if err := s.menuRepo.DeleteImage(ctx, image); err != nil {
		return err
	}

	if image.IsPrimary {
		nextImage, err := s.menuRepo.GetNextPrimaryImage(ctx, image.MenuID, menuImageID)
		if err == nil {
			s.menuRepo.SetImagePrimary(ctx, nextImage)
		}
	}

	s.redisStore.FlushByPattern(ctx, "menu:all:*")
	s.redisStore.Delete(ctx, fmt.Sprintf("menu:item:%d", image.MenuID))

	return nil
}

func (s *MenuService) ToggleMenuAvailabilityService(ctx context.Context, menuID uint, isAvailable *bool) error {

	menu, err := s.menuRepo.GetByID(ctx, menuID)
	if err != nil {
		return err
	}

	if isAvailable != nil {
		menu.IsAvailable = *isAvailable
	}

	err = s.menuRepo.Update(ctx, menu)
	if err != nil {
		return err
	}

	s.redisStore.FlushByPattern(ctx, "menu:all:*")

	return nil
}

func (s *MenuService) GetAllMenuService(ctx context.Context, filter dto.MenuFilterRequest) ([]*dto.MenuResponse, *utils.PaginatedMeta, error) {
	cacheKey := fmt.Sprintf("menu:all:page:%d:size:%d", filter.Page, filter.PageSize)
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

	menus, count, err := s.menuRepo.GetAll(ctx, filter)
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
	s.redisStore.Set(ctx, cacheKey, string(data), 1*time.Hour)

	return response, meta, nil
}

func (s *MenuService) SearchProduct(ctx context.Context, req *dto.MenuSearchRequest) ([]*dto.MenuSearchResponse, *utils.PaginatedMeta, error) {
	cacheKey := utils.BuildMenuCacheKey(req)
	cached, err := s.redisStore.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cachedResponse struct {
			Data []*dto.MenuSearchResponse `json:"data"`
			Meta *utils.PaginatedMeta      `json:"meta"`
		}
		if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
			return cachedResponse.Data, cachedResponse.Meta, nil
		}
	}

	menus, count, err := s.menuRepo.TsvectorSearchMenuItems(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	response := make([]*dto.MenuSearchResponse, len(menus))

	for i, menu := range menus {
		response[i] = &dto.MenuSearchResponse{
			MenuResponse: *s.ConvertToMenuResponse(&menu),
			Rank:         0.0,
		}
	}

	totalPages := int((count + int64(req.Limit) - 1) / int64(req.Limit))
	meta := &utils.PaginatedMeta{
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      count,
		TotalPages: totalPages,
	}

	cacheData := struct {
		Data []*dto.MenuSearchResponse `json:"data"`
		Meta *utils.PaginatedMeta      `json:"meta"`
	}{Data: response, Meta: meta}

	data, _ := json.Marshal(&cacheData)
	s.redisStore.Set(ctx, cacheKey, string(data), utils.GetCacheTTL(req))

	return response, meta, nil
}
