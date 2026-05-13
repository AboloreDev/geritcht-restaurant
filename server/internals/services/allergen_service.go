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

type AllergenService struct {
	db         *gorm.DB
	redisStore interfaces.Cacher
}

func NewAllergenService(db *gorm.DB, redisStore interfaces.Cacher) *AllergenService {
	return &AllergenService{
		db:         db,
		redisStore: redisStore,
	}
}

func (s *AllergenService) ConvertToAllergenResponse(allergen *models.Allergen) *dto.AllergenResponse {
	return &dto.AllergenResponse{
		ID:   allergen.ID,
		Name: allergen.Name,
	}
}

func (s *AllergenService) CreateAllergenServices(req *dto.CreateAllergenRequest) (*dto.AllergenResponse, error) {
	var allergen models.Allergen

	err := s.db.Where("name = ?", req.Name).First(&allergen).Error
	if err == nil {
		return nil, domain.ErrNameConflict
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	allergen = models.Allergen{
		Name: req.Name,
	}

	err = s.db.Create(&allergen).Error
	if err != nil {
		return nil, err
	}

	s.redisStore.FlushByPattern(ctx, "menu:allergen:*")

	return s.ConvertToAllergenResponse(&allergen), nil
}

func (s *AllergenService) GetAllAllergenService(page, pageSize int) ([]*dto.AllergenResponse, *utils.PaginatedMeta, error) {
	cacheKey := fmt.Sprintf("menu:allergen:p%d:s%d", page, pageSize)

	cached, err := s.redisStore.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cachedResponse struct {
			Data []*dto.AllergenResponse `json:"data"`
			Meta *utils.PaginatedMeta    `json:"meta"`
		}
		if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
			return cachedResponse.Data, cachedResponse.Meta, nil
		}
	}
	var allergens []models.Allergen
	var count int64

	offset := utils.Pagination(page, pageSize)

	s.db.Model(models.Allergen{}).Count(&count)

	err = s.db.Offset(offset).Limit(pageSize).Find(&allergens).Error
	if err != nil {
		return nil, nil, err
	}

	response := make([]*dto.AllergenResponse, 0, len(allergens))

	for _, allergen := range allergens {
		response = append(response, s.ConvertToAllergenResponse(&allergen))
	}

	totalPages := int((count + int64(pageSize) - 1) / int64(pageSize))

	meta := &utils.PaginatedMeta{
		Page:       page,
		Limit:      pageSize,
		Total:      count,
		TotalPages: totalPages,
	}

	cacheData := struct {
		Data []*dto.AllergenResponse `json:"data"`
		Meta *utils.PaginatedMeta    `json:"meta"`
	}{Data: response, Meta: meta}

	// store in cache
	data, err := json.Marshal(&cacheData)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to set data: %d", err)
	}
	s.redisStore.Set(ctx, cacheKey, string(data), 1*time.Hour)

	return response, meta, nil
}

func (s *AllergenService) UpdateAllergenService(allergenID uint, req *dto.UpdateAllergenRequest) (*dto.AllergenResponse, error) {
	var allergen models.Allergen

	err := s.db.Where("id = ? ", allergenID).First(&allergen).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	allergen.Name = req.Name

	err = s.db.Save(&allergen).Error
	if err != nil {
		return nil, err
	}

	s.redisStore.Delete(ctx,
		fmt.Sprintf("menu:allergen:%d", allergenID),
	)

	s.redisStore.FlushByPattern(ctx, "menu:allergen:*")

	return s.ConvertToAllergenResponse(&allergen), nil
}

func (s *AllergenService) DeleteAllergenService(allergenID uint) error {
	var allergen models.Allergen
	var count int64

	s.db.Model(models.Menu{}).Where("allergens = ?", allergenID).Count(&count)
	if count > 0 {
		return fmt.Errorf("cannot delete tags with %d menu", count)
	}

	result := s.db.Where("id = ?", allergenID).Delete(&allergen)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	s.redisStore.Delete(ctx, fmt.Sprintf("menu:allergen:%d", allergenID))
	s.redisStore.FlushByPattern(ctx, "menu:allergen:*")

	return nil
}
