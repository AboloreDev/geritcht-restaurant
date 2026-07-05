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

type AllergenService struct {
	allergenRepo repositories.AllergenRepositoryInterface
	redisStore   interfaces.Cacher
}

func NewAllergenService(allergenRepo repositories.AllergenRepositoryInterface, redisStore interfaces.Cacher) *AllergenService {
	return &AllergenService{allergenRepo: allergenRepo, redisStore: redisStore}
}

func (s *AllergenService) ConvertToAllergenResponse(allergen *models.Allergen) *dto.AllergenResponse {
	return &dto.AllergenResponse{ID: allergen.ID, Name: allergen.Name}
}

func (s *AllergenService) CreateAllergenServices(ctx context.Context, req *dto.CreateAllergenRequest) (*dto.AllergenResponse, error) {
	_, err := s.allergenRepo.GetByName(ctx, req.Name)
	if err == nil {
		return nil, domain.ErrNameConflict
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	allergen := models.Allergen{Name: req.Name}
	if err := s.allergenRepo.Create(ctx, &allergen); err != nil {
		return nil, err
	}

	s.redisStore.FlushByPattern(ctx, "menu:allergen:*")
	return s.ConvertToAllergenResponse(&allergen), nil
}

func (s *AllergenService) GetAllAllergenService(ctx context.Context, page, pageSize int) ([]*dto.AllergenResponse, *utils.PaginatedMeta, error) {
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

	allergens, count, err := s.allergenRepo.GetAll(ctx, page, pageSize)
	if err != nil {
		return nil, nil, err
	}

	response := make([]*dto.AllergenResponse, 0, len(allergens))
	for _, allergen := range allergens {
		response = append(response, s.ConvertToAllergenResponse(&allergen))
	}

	totalPages := int((count + int64(pageSize) - 1) / int64(pageSize))
	meta := &utils.PaginatedMeta{Page: page, Limit: pageSize, Total: count, TotalPages: totalPages}

	cacheData := struct {
		Data []*dto.AllergenResponse `json:"data"`
		Meta *utils.PaginatedMeta    `json:"meta"`
	}{Data: response, Meta: meta}
	data, _ := json.Marshal(&cacheData)
	s.redisStore.Set(ctx, cacheKey, string(data), 1*time.Hour)

	return response, meta, nil
}

func (s *AllergenService) UpdateAllergenService(ctx context.Context, allergenID uint, req *dto.UpdateAllergenRequest) (*dto.AllergenResponse, error) {
	allergen, err := s.allergenRepo.GetByID(ctx, allergenID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	allergen.Name = req.Name
	if err := s.allergenRepo.Update(ctx, allergen); err != nil {
		return nil, err
	}

	s.redisStore.Delete(ctx, fmt.Sprintf("menu:allergen:%d", allergenID))
	s.redisStore.FlushByPattern(ctx, "menu:allergen:*")
	return s.ConvertToAllergenResponse(allergen), nil
}

func (s *AllergenService) DeleteAllergenService(ctx context.Context, allergenID uint) error {
	count, err := s.allergenRepo.CountMenuItemsUsingAllergen(ctx, allergenID)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete tags with %d menu", count)
	}

	if err := s.allergenRepo.Delete(ctx, allergenID); err != nil {
		return domain.ErrNotFound
	}

	s.redisStore.Delete(ctx, fmt.Sprintf("menu:allergen:%d", allergenID))
	s.redisStore.FlushByPattern(ctx, "menu:allergen:*")
	return nil
}
