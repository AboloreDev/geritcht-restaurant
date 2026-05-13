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

type DietaryTagsService struct {
	db         *gorm.DB
	redisStore interfaces.Cacher
}

func NewDietaryTagsService(db *gorm.DB, redisStore interfaces.Cacher) *DietaryTagsService {
	return &DietaryTagsService{
		db:         db,
		redisStore: redisStore,
	}
}

func (s *DietaryTagsService) ConvertToDietaryTagResponse(dietaryTags *models.DietaryTag) *dto.DietaryTagResponse {
	return &dto.DietaryTagResponse{
		ID:   dietaryTags.ID,
		Name: dietaryTags.Name,
	}
}

func (s *DietaryTagsService) CreateDietaryTagService(req *dto.CreateDietaryTagRequest) (*dto.DietaryTagResponse, error) {
	var tags models.DietaryTag

	err := s.db.Where("name = ?", req.Name).First(&tags).Error
	if err == nil {
		return nil, domain.ErrNameConflict
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	tags = models.DietaryTag{
		Name: req.Name,
	}

	err = s.db.Create(&tags).Error
	if err != nil {
		return nil, err
	}

	s.redisStore.FlushByPattern(ctx, "menu:tags:*")

	return s.ConvertToDietaryTagResponse(&tags), nil
}

func (s *DietaryTagsService) GetAllDietaryTagService(page, pageSize int) ([]*dto.DietaryTagResponse, *utils.PaginatedMeta, error) {
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
	var dietaryTags []models.DietaryTag
	var count int64

	offset := utils.Pagination(page, pageSize)

	s.db.Model(models.DietaryTag{}).Count(&count)

	err = s.db.Offset(offset).Limit(pageSize).Find(&dietaryTags).Error
	if err != nil {
		return nil, nil, err
	}

	response := make([]*dto.DietaryTagResponse, 0, len(dietaryTags))

	for _, tag := range dietaryTags {
		response = append(response, s.ConvertToDietaryTagResponse(&tag))
	}

	totalPages := int((count + int64(pageSize) - 1) / int64(pageSize))

	meta := &utils.PaginatedMeta{
		Page:       page,
		Limit:      pageSize,
		Total:      count,
		TotalPages: totalPages,
	}

	cacheData := struct {
		Data []*dto.DietaryTagResponse `json:"data"`
		Meta *utils.PaginatedMeta      `json:"meta"`
	}{Data: response, Meta: meta}

	// store in cache
	data, err := json.Marshal(&cacheData)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to set data: %d", err)
	}
	s.redisStore.Set(ctx, cacheKey, string(data), 1*time.Hour)

	return response, meta, nil
}

func (s *DietaryTagsService) UpdateDietaryTagService(tagID uint, req *dto.UpdateDietaryTagRequest) (*dto.DietaryTagResponse, error) {
	var dietaryTag models.DietaryTag

	err := s.db.Where("id = ? ", tagID).First(&dietaryTag).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	dietaryTag.Name = req.Name

	err = s.db.Save(&dietaryTag).Error
	if err != nil {
		return nil, err
	}

	s.redisStore.Delete(ctx,
		fmt.Sprintf("menu:tags:%d", tagID),
	)

	s.redisStore.FlushByPattern(ctx, "menu:tags:*")

	return s.ConvertToDietaryTagResponse(&dietaryTag), nil
}

func (s *DietaryTagsService) DeleteDietaryTagService(tagID uint) error {
	var dietaryTag models.DietaryTag
	var count int64

	s.db.Model(models.Menu{}).Where("dietary_tags = ?", tagID).Count(&count)
	if count > 0 {
		return fmt.Errorf("cannot delete tags with %d menu", count)
	}

	result := s.db.Where("id = ?", tagID).Delete(&dietaryTag)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	s.redisStore.Delete(ctx, fmt.Sprintf("menu:tags:%d", tagID))
	s.redisStore.FlushByPattern(ctx, "menu:tags:*")

	return nil
}
