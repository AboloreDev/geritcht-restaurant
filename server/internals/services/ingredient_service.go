package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/events"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
)

type IngredientService struct {
	redisStore     interfaces.Cacher
	eventPublisher interfaces.Publisher
	ingredientRepo repositories.IngredientRepositoryInterface
	userRepo       repositories.UserRepositoryInterface
	paymentRepo    repositories.PaymentRepositoryInterface
}

func NewIngredientService(
	redisStore interfaces.Cacher,
	eventPublisher interfaces.Publisher,
	ingredientRepo repositories.IngredientRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
	paymentRepo repositories.PaymentRepositoryInterface) *IngredientService {
	return &IngredientService{
		redisStore:     redisStore,
		eventPublisher: eventPublisher,
		ingredientRepo: ingredientRepo,
		userRepo:       userRepo,
		paymentRepo:    paymentRepo,
	}
}

func (s *IngredientService) CreateIngredientService(ctx context.Context, req *dto.CreateIngredientRequest) (*dto.IngredientResponse, error) {
	_, err := s.ingredientRepo.GetIngredientByName(ctx, req.Name)
	if err == nil {
		return nil, domain.ErrNameConflict
	}

	ingredient := &models.Ingredient{
		Name:         req.Name,
		Unit:         req.Unit,
		CurrentStock: req.CurrentStock,
		MinThreshold: req.MinThreshold,
	}

	err = s.ingredientRepo.CreateIngredient(ctx, ingredient)
	if err != nil {
		return nil, err
	}

	s.redisStore.FlushByPattern(ctx, "menu:ingredients:*")

	return &dto.IngredientResponse{
		ID:           ingredient.ID,
		Name:         ingredient.Name,
		Unit:         ingredient.Unit,
		CurrentStock: ingredient.CurrentStock,
		MinThreshold: ingredient.MinThreshold,
		CreatedAt:    time.Now(),
	}, nil
}

func (s *IngredientService) UpdateIngredientService(ctx context.Context, ingredientID uint, req *dto.UpdateIngredientRequest) (*dto.IngredientResponse, error) {

	ingredient, err := s.ingredientRepo.GetIngredientByID(ctx, ingredientID)
	if err != nil {
		return nil, domain.ErrIngredientNotFound
	}

	if req.Name != "" {
		ingredient.Name = req.Name
	}

	if req.Unit != "" {
		ingredient.Unit = req.Unit
	}

	if req.MinThreshold >= 0 {
		ingredient.MinThreshold = req.MinThreshold
	}

	err = s.ingredientRepo.UpdateIngredient(ctx, ingredient)
	if err != nil {
		return nil, err
	}

	s.redisStore.Delete(ctx,
		fmt.Sprintf("menu:ingredient:%d", ingredientID),
	)

	s.redisStore.FlushByPattern(ctx, "menu:ingredients:*")

	return &dto.IngredientResponse{
		ID:           ingredient.ID,
		Name:         ingredient.Name,
		Unit:         ingredient.Unit,
		CurrentStock: ingredient.CurrentStock,
		MinThreshold: ingredient.MinThreshold,
	}, nil
}

func (s *IngredientService) DeleteIngredientService(ctx context.Context, ingredientID uint) error {

	count, err := s.ingredientRepo.IngredientCount(ctx, ingredientID)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("cannot delete ingredients with %d menuitems", count)
	}

	err = s.ingredientRepo.DeleteIngredient(ctx, ingredientID)
	if err != nil {
		return err
	}

	s.redisStore.Delete(ctx, fmt.Sprintf("menu:ingredient:%d", ingredientID))
	s.redisStore.FlushByPattern(ctx, "menu:ingredients:*")

	return nil
}

func (s *IngredientService) GetIngredientService(ctx context.Context, ingredientID uint) (*dto.IngredientResponse, error) {
	cachedKey := fmt.Sprintf("menu:ingredient:%d", ingredientID)

	exists, _ := s.redisStore.Exists(ctx, cachedKey)
	if exists {
		cached, err := s.redisStore.Get(ctx, cachedKey)
		if err == nil && cached != "" {
			var ingredient models.Ingredient
			isLow := ingredient.CurrentStock <= ingredient.MinThreshold
			err = json.Unmarshal([]byte(cached), &ingredient)
			if err != nil {
				return nil, err
			}
			return &dto.IngredientResponse{
				ID:           ingredient.ID,
				Name:         ingredient.Name,
				Unit:         ingredient.Unit,
				CurrentStock: ingredient.CurrentStock,
				IsLow:        isLow,
				MinThreshold: ingredient.MinThreshold,
				CreatedAt:    ingredient.CreatedAt,
			}, nil
		}
	}

	ingredient, err := s.ingredientRepo.GetIngredientByID(ctx, ingredientID)
	if err != nil {
		return nil, domain.ErrIngredientNotFound
	}

	isLow := ingredient.CurrentStock <= ingredient.MinThreshold

	data, err := json.Marshal(&ingredient)
	if err != nil {
		return nil, fmt.Errorf("Failed to set data: %d", err)
	}
	s.redisStore.Set(ctx, cachedKey, string(data), 1*time.Hour)

	return &dto.IngredientResponse{
		ID:           ingredient.ID,
		Name:         ingredient.Name,
		Unit:         ingredient.Unit,
		CurrentStock: ingredient.CurrentStock,
		IsLow:        isLow,
		MinThreshold: ingredient.MinThreshold,
		CreatedAt:    ingredient.CreatedAt,
	}, nil
}

func (s *IngredientService) GetAllIngredientService(ctx context.Context, page, pageSize int) ([]*dto.IngredientResponse, *utils.PaginatedMeta, error) {
	cacheKey := fmt.Sprintf("menu:ingredients:p%d:s%d", page, pageSize)

	cached, err := s.redisStore.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cachedResponse struct {
			Data []*dto.IngredientResponse `json:"data"`
			Meta *utils.PaginatedMeta      `json:"meta"`
		}
		if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
			return cachedResponse.Data, cachedResponse.Meta, nil
		}
	}

	ingredients, count, err := s.ingredientRepo.GetAllIngredients(ctx, page, pageSize)

	response := make([]*dto.IngredientResponse, 0, len(ingredients))

	for _, ingredient := range ingredients {
		isLow := ingredient.CurrentStock <= ingredient.MinThreshold
		response = append(response, &dto.IngredientResponse{
			ID:           ingredient.ID,
			Name:         ingredient.Name,
			Unit:         ingredient.Unit,
			CurrentStock: ingredient.CurrentStock,
			MinThreshold: ingredient.MinThreshold,
			IsLow:        isLow,
			CreatedAt:    ingredient.CreatedAt,
		})
	}

	totalPages := int((count + int64(pageSize) - 1) / int64(pageSize))

	meta := &utils.PaginatedMeta{
		Page:       page,
		Limit:      pageSize,
		Total:      count,
		TotalPages: totalPages,
	}

	cacheData := struct {
		Data []*dto.IngredientResponse `json:"data"`
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

func (s *IngredientService) GetLowStockIngredientsService(ctx context.Context) ([]*dto.IngredientResponse, error) {
	ingredients, err := s.ingredientRepo.CompareCurrentStockAgainstMinTheshold(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.IngredientResponse, 0, len(ingredients))

	for _, ingredient := range ingredients {
		isLow := ingredient.CurrentStock <= ingredient.MinThreshold
		response = append(response, &dto.IngredientResponse{
			ID:           ingredient.ID,
			Name:         ingredient.Name,
			Unit:         ingredient.Unit,
			CurrentStock: ingredient.CurrentStock,
			MinThreshold: ingredient.MinThreshold,
			IsLow:        isLow,
			CreatedAt:    ingredient.CreatedAt,
		})
	}

	return response, nil
}

func (s *IngredientService) SetThresholdLimit(ctx context.Context, ingredientID uint, req *dto.ThresholdRequest) error {
	ingredient, err := s.ingredientRepo.GetIngredientByID(ctx, ingredientID)
	if err != nil {
		return domain.ErrIngredientNotFound
	}

	if req.Threshold <= 0 {
		return domain.ErrNegativeThreshold
	}

	err = s.ingredientRepo.UpdateThreshHoldLimit(ctx, ingredient.ID, req.Threshold)
	if err != nil {
		return fmt.Errorf("failed to update threshold: %w", err)
	}

	return nil
}

func (s *IngredientService) CheckLowStock(ctx context.Context, userID, ingredientID uint) error {
	ingredient, err := s.ingredientRepo.GetIngredientByID(ctx, ingredientID)
	if err != nil {
		return domain.ErrIngredientNotFound
	}

	user, err := s.userRepo.GetByIDAndRole(ctx, userID, models.RoleAdmin)
	if err != nil {
		return domain.ErrUserNotFound
	}

	isLow := ingredient.CurrentStock <= ingredient.MinThreshold

	if isLow {
		return s.sendLowStockAlert(ctx, user, ingredient)
	}

	return nil
}

func (s *IngredientService) sendLowStockAlert(ctx context.Context, user *models.User, ingredient *models.Ingredient) error {
	lowStockAlert, err := json.Marshal(events.LowStockAlertPayload{
		AdminEmail: user.Email,
		AdminName:  user.FirstName,
		Items: []events.LowStockPayload{
			{
				Name:         ingredient.Name,
				CurrentStock: ingredient.CurrentStock,
				MinThreshold: ingredient.MinThreshold,
			},
		},
	})
	if err != nil {
		return err
	}
	err = s.paymentRepo.CreateOutboxEvent(ctx, nil, &models.OutboxEvent{
		EventType: events.ChannelEmailLowStockAlert,
		Payload:   string(lowStockAlert),
		Status:    "pending",
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	err = s.eventPublisher.PublishMessage(
		events.ChannelEmailLowStockAlert,
		&events.LowStockAlertPayload{
			AdminEmail: user.Email,
			AdminName:  user.FirstName,
			Items: []events.LowStockPayload{
				{
					Name:         ingredient.Name,
					CurrentStock: ingredient.CurrentStock,
					MinThreshold: ingredient.MinThreshold,
				},
			},
		},
		map[string]string{"Priority": "Important Mail"},
	)
	if err != nil {
		return nil
	}

	s.paymentRepo.MarkOutboxPublished(ctx, events.ChannelEmailLowStockAlert)

	return nil
}

func (s *IngredientService) SearchIngredients(ctx context.Context, req *dto.IngredientSearchRequest) ([]*dto.IngredientSearchResponse, *utils.PaginatedMeta, error) {
	cacheKey := fmt.Sprintf("menu:ingredients:%s:p%d:s%d", req.Query, req.Page, req.Limit)
	cached, err := s.redisStore.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cachedResponse struct {
			Data []*dto.IngredientSearchResponse `json:"data"`
			Meta *utils.PaginatedMeta            `json:"meta"`
		}
		if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
			return cachedResponse.Data, cachedResponse.Meta, nil
		}
	}

	rows, count, err := s.ingredientRepo.TsvectorSearchIngredients(ctx, req)
	if err != nil {
		return nil, nil, domain.ErrIngredientSearchNotFound
	}

	response := make([]*dto.IngredientSearchResponse, len(rows))

	for i := range rows {
		response[i] = &dto.IngredientSearchResponse{
			IngredientResponse: dto.IngredientResponse{
				ID:           rows[i].ID,
				Name:         rows[i].Name,
				Unit:         rows[i].Unit,
				CurrentStock: rows[i].CurrentStock,
				MinThreshold: rows[i].MinThreshold,
				CreatedAt:    rows[i].CreatedAt,
			},
			Rank: (rows[i].Rank),
		}
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	totalPages := int((count + int64(req.Limit) - 1) / int64(req.Limit))
	meta := &utils.PaginatedMeta{
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      count,
		TotalPages: totalPages,
	}

	cacheData := struct {
		Data []*dto.IngredientSearchResponse `json:"data"`
		Meta *utils.PaginatedMeta            `json:"meta"`
	}{Data: response, Meta: meta}

	data, _ := json.Marshal(&cacheData)
	s.redisStore.Set(ctx, cacheKey, string(data), 1*time.Hour)

	return response, meta, nil
}
