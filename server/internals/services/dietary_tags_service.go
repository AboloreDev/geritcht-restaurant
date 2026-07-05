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
	"gorm.io/gorm"
)

type DietaryTagsService struct {
	tagRepo    repositories.DietaryTagRepositoryInterface
	redisStore interfaces.Cacher
}

func NewDietaryTagsService(tagRepo repositories.DietaryTagRepositoryInterface, redisStore interfaces.Cacher) *DietaryTagsService {
	return &DietaryTagsService{tagRepo: tagRepo, redisStore: redisStore}
}

func (s *DietaryTagsService) ConvertToDietaryTagResponse(tag *models.DietaryTag) *dto.DietaryTagResponse {
	return &dto.DietaryTagResponse{ID: tag.ID, Name: tag.Name}
}

func (s *DietaryTagsService) CreateDietaryTagService(ctx context.Context, req *dto.CreateDietaryTagRequest) (*dto.DietaryTagResponse, error) {
	_, err := s.tagRepo.GetByName(ctx, req.Name)
	if err == nil {
		return nil, domain.ErrNameConflict
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	tag := models.DietaryTag{Name: req.Name}
	if err := s.tagRepo.Create(ctx, &tag); err != nil {
		return nil, err
	}

	s.redisStore.FlushByPattern(ctx, "menu:tags:*")
	return s.ConvertToDietaryTagResponse(&tag), nil
}

func (s *DietaryTagsService) GetAllDietaryTagService(ctx context.Context, page, pageSize int) ([]*dto.DietaryTagResponse, *utils.PaginatedMeta, error) {
	cacheKey := fmt.Sprintf("menu:tags:p%d:s%d", page, pageSize)

	cached, err := s.redisStore.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cachedResponse struct {
			Data []*dto.DietaryTagResponse `json:"data"`
			Meta *utils.PaginatedMeta      `json:"meta"`
		}
		if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
			return cachedResponse.Data, cachedResponse.Meta, nil
		}
	}

	tags, count, err := s.tagRepo.GetAll(ctx, page, pageSize)
	if err != nil {
		return nil, nil, err
	}

	response := make([]*dto.DietaryTagResponse, 0, len(tags))
	for _, tag := range tags {
		response = append(response, s.ConvertToDietaryTagResponse(&tag))
	}

	totalPages := int((count + int64(pageSize) - 1) / int64(pageSize))
	meta := &utils.PaginatedMeta{Page: page, Limit: pageSize, Total: count, TotalPages: totalPages}

	cacheData := struct {
		Data []*dto.DietaryTagResponse `json:"data"`
		Meta *utils.PaginatedMeta      `json:"meta"`
	}{Data: response, Meta: meta}
	data, _ := json.Marshal(&cacheData)
	s.redisStore.Set(ctx, cacheKey, string(data), 1*time.Hour)

	return response, meta, nil
}

func (s *DietaryTagsService) UpdateDietaryTagService(ctx context.Context, tagID uint, req *dto.UpdateDietaryTagRequest) (*dto.DietaryTagResponse, error) {
	tag, err := s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	tag.Name = req.Name
	if err := s.tagRepo.Update(ctx, tag); err != nil {
		return nil, err
	}

	s.redisStore.Delete(ctx, fmt.Sprintf("menu:tags:%d", tagID))
	s.redisStore.FlushByPattern(ctx, "menu:tags:*")
	return s.ConvertToDietaryTagResponse(tag), nil
}

func (s *DietaryTagsService) DeleteDietaryTagService(ctx context.Context, tagID uint) error {
	count, err := s.tagRepo.CountMenuItemsUsingTag(ctx, tagID)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete tags with %d menu", count)
	}

	if err := s.tagRepo.Delete(ctx, tagID); err != nil {
		return domain.ErrNotFound
	}

	s.redisStore.Delete(ctx, fmt.Sprintf("menu:tags:%d", tagID))
	s.redisStore.FlushByPattern(ctx, "menu:tags:*")
	return nil
}
