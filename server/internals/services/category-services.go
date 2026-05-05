package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
)

var ctx = context.Background()

type CategoryService struct {
	db         *gorm.DB
	redisStore interfaces.Cacher
}

func NewCategoryService(db *gorm.DB, redisStore interfaces.Cacher) *CategoryService {
	return &CategoryService{
		db:         db,
		redisStore: redisStore,
	}
}

func (s *CategoryService) ConvertCategoryResponse(category *models.MenuCategory) *dto.MenuCategoryResponse {
	if category == nil {
		return nil
	}

	menuItem := make([]dto.MenuResponse, 0, len(category.Menu))

	for _, item := range category.Menu {
		menuItem = append(menuItem, dto.MenuResponse{
			ID:         item.ID,
			CategoryID: item.CategoryID,
			Name:       item.Name,
		})
	}

	return &dto.MenuCategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		ImageURL:    category.ImageURL,
		MenuItems:   menuItem,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
	}
}

func (s *CategoryService) CreateCategoryService(req *dto.CreateCategoryRequest, imageURL string) (*dto.MenuCategoryResponse, error) {
	var existingCategory models.MenuCategory

	err := s.db.Where("name = ?", req.Name).First(&existingCategory).Error
	if err == nil {
		return nil, domain.ErrNameConflict
	}

	category := models.MenuCategory{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
		CreatedAt:   time.Now(),
		ImageURL:    imageURL,
	}

	err = s.db.Create(&category).Error
	if err != nil {
		return nil, err
	}

	s.redisStore.Delete(ctx, "menu:categories")

	return s.ConvertCategoryResponse(&category), nil
}

func (s *CategoryService) UpdateCategoryService(categoryID uint, req *dto.UpdateCategoryRequest) (*dto.MenuCategoryResponse, error) {
	var category models.MenuCategory

	err := s.db.Where("id = ? ", categoryID).First(&category).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	category.Name = req.Name
	category.Description = req.Description
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	err = s.db.Save(&category).Error
	if err != nil {
		return nil, err
	}

	s.redisStore.Delete(ctx,
		"menu:categories",
		fmt.Sprintf("menu:category:%d", categoryID),
	)

	return s.ConvertCategoryResponse(&category), nil

}

func (s *CategoryService) DeleteCategoryService(categoryID uint) error {
	var category models.MenuCategory
	var count int64

	s.db.Model(models.Menu{}).Where("category_id = ?", categoryID).Count(&count)
	if count > 0 {
		return fmt.Errorf("cannot delete category with %d menu", count)
	}

	err := s.db.Delete(&category, categoryID).Error
	if err != nil {
		return err
	}

	s.redisStore.Delete(ctx,
		"menu:categories",
		fmt.Sprintf("menu:category:%d", categoryID),
	)

	return nil
}

func (s *CategoryService) GetCategoriesService(page, pageSize int) ([]*dto.MenuCategoryResponse, *utils.PaginatedMeta, error) {
	cacheKey := fmt.Sprintf("menu:categories:p%d:s%d", page, pageSize)

	cached, err := s.redisStore.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cachedResponse struct {
			Data []*dto.MenuCategoryResponse `json:"data"`
			Meta *utils.PaginatedMeta        `json:"meta"`
		}
		if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
			return cachedResponse.Data, cachedResponse.Meta, nil
		}
	}
	var categories []models.MenuCategory
	var count int64

	offset := utils.Pagination(page, pageSize)

	s.db.Model(models.MenuCategory{}).Count(&count)

	err = s.db.Preload("Menu").Offset(offset).Limit(pageSize).Find(&categories).Error
	if err != nil {
		return nil, nil, err
	}

	response := make([]*dto.MenuCategoryResponse, 0, len(categories))

	for _, category := range categories {
		response = append(response, s.ConvertCategoryResponse(&category))
	}

	totalPages := int(count + int64(pageSize) - 1/int64(pageSize))

	meta := &utils.PaginatedMeta{
		Page:       page,
		Limit:      pageSize,
		Total:      count,
		TotalPages: totalPages,
	}

	cacheData := struct {
		Data []*dto.MenuCategoryResponse `json:"data"`
		Meta *utils.PaginatedMeta        `json:"meta"`
	}{Data: response, Meta: meta}

	// store in cache
	data, err := json.Marshal(&cacheData)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to set data: %d", err)
	}
	s.redisStore.Set(ctx, cacheKey, string(data), 1*time.Hour)

	return response, meta, nil
}

func (s *CategoryService) GetCategoryService(categoryID uint) (*dto.MenuCategoryResponse, error) {
	cachedKey := fmt.Sprintf("menu:category:%d", categoryID)

	exists, _ := s.redisStore.Exists(ctx, cachedKey)
	if exists {
		cached, err := s.redisStore.Get(ctx, cachedKey)
		if err == nil && cached != "" {
			var category models.MenuCategory
			err = json.Unmarshal([]byte(cached), &category)
			if err != nil {
				return nil, err
			}
			return s.ConvertCategoryResponse(&category), nil
		}
	}

	var category models.MenuCategory
	err := s.db.Preload("Menu").Where("id = ? AND is_active = ?", categoryID, true).First(&category).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	data, err := json.Marshal(&category)
	if err != nil {
		return nil, fmt.Errorf("Failed to set data: %d", err)
	}
	s.redisStore.Set(ctx, cachedKey, string(data), 1*time.Hour)

	return s.ConvertCategoryResponse(&category), nil
}
