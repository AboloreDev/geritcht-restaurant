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
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
)

var ctx = context.Background()

type CategoryService struct {
	redisStore   interfaces.Cacher
	categoryRepo repositories.CategoryRepositoryInterface
}

func NewCategoryService(
	redisStore interfaces.Cacher,
	categoryRepo repositories.CategoryRepositoryInterface) *CategoryService {
	return &CategoryService{
		redisStore:   redisStore,
		categoryRepo: categoryRepo,
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
			CategoryID: item.MenuCategoryID,
			Name:       item.Name,
		})
	}

	return &dto.MenuCategoryResponse{
		ID:           category.ID,
		Name:         category.Name,
		Description:  category.Description,
		ImageURL:     category.ImageURL,
		MenuItems:    menuItem,
		IsActive:     category.IsActive,
		CreatedAt:    category.CreatedAt,
		DisplayOrder: category.DisplayOrder,
	}
}

func (s *CategoryService) CreateCategoryService(ctx context.Context, req *dto.CreateCategoryRequest, imageURL string) (*dto.MenuCategoryResponse, error) {

	_, err := s.categoryRepo.GetByName(ctx, req.Name)
	if err == nil {
		return nil, domain.ErrNameConflict
	}

	category := models.MenuCategory{
		Name:         req.Name,
		Description:  req.Description,
		IsActive:     true,
		CreatedAt:    time.Now(),
		ImageURL:     imageURL,
		DisplayOrder: req.DisplayOrder,
	}

	err = s.categoryRepo.Create(ctx, &category)
	if err != nil {
		return nil, err
	}

	s.redisStore.FlushByPattern(ctx, "menu:categories:*")

	return s.ConvertCategoryResponse(&category), nil
}

func (s *CategoryService) UpdateCategoryService(ctx context.Context, categoryID uint, req *dto.UpdateCategoryRequest) (*dto.MenuCategoryResponse, error) {

	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.DisplayOrder != 0 {
		category.DisplayOrder = req.DisplayOrder
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	err = s.categoryRepo.Update(ctx, category)
	if err != nil {
		return nil, err
	}

	s.redisStore.Delete(ctx,
		fmt.Sprintf("menu:category:%d", categoryID),
	)

	s.redisStore.FlushByPattern(ctx, "menu:categories:*")

	return s.ConvertCategoryResponse(category), nil

}

func (s *CategoryService) DeleteCategoryService(ctx context.Context, categoryID uint) error {

	count, err := s.categoryRepo.CountMenuItems(ctx, categoryID)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete category with %d menu items", count)
	}

	if err := s.categoryRepo.Delete(ctx, categoryID); err != nil {
		return err
	}

	s.redisStore.Delete(ctx, fmt.Sprintf("menu:category:%d", categoryID))
	s.redisStore.FlushByPattern(ctx, "menu:categories:*")

	return nil
}

func (s *CategoryService) GetCategoriesService(ctx context.Context, page, pageSize int) ([]*dto.MenuCategoryResponse, *utils.PaginatedMeta, error) {
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

	categories, count, err := s.categoryRepo.GetAll(ctx, page, pageSize)
	if err != nil {
		return nil, nil, err
	}

	response := make([]*dto.MenuCategoryResponse, 0, len(categories))

	for _, category := range categories {
		response = append(response, s.ConvertCategoryResponse(&category))
	}

	totalPages := int((count + int64(pageSize) - 1) / int64(pageSize))

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

func (s *CategoryService) GetCategoryService(ctx context.Context, categoryID uint) (*dto.MenuCategoryResponse, error) {
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

	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	data, err := json.Marshal(&category)
	if err != nil {
		return nil, fmt.Errorf("Failed to set data: %d", err)
	}

	s.redisStore.Set(ctx, cachedKey, string(data), 1*time.Hour)

	return s.ConvertCategoryResponse(category), nil
}
