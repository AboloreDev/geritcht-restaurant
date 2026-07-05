package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/events"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
	"gorm.io/gorm"
)

type InventoryService struct {
	db             *gorm.DB
	eventPublisher interfaces.Publisher
	redisStore     interfaces.Cacher
	inventoryRepo  repositories.InventoryRepositoryInterface
}

func NewInventoryService(
	db *gorm.DB,
	eventPublisher interfaces.Publisher,
	redisStore interfaces.Cacher,
	inventoryRepo repositories.InventoryRepositoryInterface) *InventoryService {
	return &InventoryService{
		db:             db,
		eventPublisher: eventPublisher,
		redisStore:     redisStore,
		inventoryRepo:  inventoryRepo,
	}
}

func (s *InventoryService) DeductStock(ctx context.Context, tx *gorm.DB, orderItems []models.OrderItem, orderID uint, createdBy uint) error {
	for _, orderItem := range orderItems {
		// get recipe for this menu item
		recipes, err := s.inventoryRepo.GetRecipesByMenuItemID(ctx, tx, orderItem.MenuID)
		if err != nil {
			return err
		}

		// no recipe → skip (drinks, packaged items)
		if len(recipes) == 0 {
			continue
		}

		for _, recipe := range recipes {
			required := recipe.Quantity * float64(orderItem.Quantity)

			// fetch current stock
			ingredient, err := s.inventoryRepo.GetIngredientByID(ctx, tx, recipe.IngredientID)
			if err != nil {
				return err
			}

			// check BEFORE deducting
			if ingredient.CurrentStock < required {
				return domain.ErrInsufficientStock
			}

			// atomic deduction with double check
			rowsAffected, err := s.inventoryRepo.DeductIngredientStock(ctx, tx, recipe.IngredientID, required)
			if err != nil {
				return err
			}
			if rowsAffected == 0 {
				return domain.ErrInsufficientStock
			}

			// log stock movement
			if err := s.inventoryRepo.CreateStockMovement(ctx, tx, &models.StockMovement{
				IngredientID: recipe.IngredientID,
				Quantity:     required,
				Type:         models.StockMovementOut,
				Reason:       fmt.Sprintf("Order #%d deduction", orderID),
				CreatedBy:    createdBy,
				CreatedAt:    time.Now(),
			}); err != nil {
				return err
			}
		}
	}

	// check thresholds after all deductions
	return s.CheckAndAlertThreshold(ctx, tx)
}

func (s *InventoryService) CheckAndAlertThreshold(ctx context.Context, tx *gorm.DB) error {
	admin, err := s.inventoryRepo.GetAdminUser(ctx, tx)
	if err != nil {
		return err
	}

	lowIngredients, err := s.inventoryRepo.GetLowStockIngredients(ctx, tx)
	if err != nil {
		return err
	}

	for _, ingredient := range lowIngredients {
		payload, err := json.Marshal(events.LowStockAlertPayload{
			AdminEmail: admin.Email,
			AdminName:  admin.FirstName,
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

		outbox := &models.OutboxEvent{
			EventType: events.ChannelEmailLowStockAlert,
			Payload:   string(payload),
			Status:    "pending",
			CreatedAt: time.Now(),
		}

		if err := s.inventoryRepo.CreateOutboxEvent(ctx, tx, outbox); err != nil {
			return err
		}

		err = s.eventPublisher.PublishMessage(
			events.ChannelEmailLowStockAlert,
			&events.LowStockAlertPayload{
				AdminEmail: admin.Email,
				AdminName:  admin.FirstName,
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
			continue // outbox worker retries
		}

		s.inventoryRepo.MarkOutboxPublished(ctx, tx, outbox.ID)
	}

	outOfStock, err := s.inventoryRepo.GetOutOfStockIngredients(ctx, tx)
	if err != nil {
		return err
	}

	for _, ingredient := range outOfStock {
		menuItemIDs, err := s.inventoryRepo.GetMenuItemIDsByIngredient(ctx, tx, ingredient.ID)
		if err != nil {
			continue
		}
		if len(menuItemIDs) > 0 {
			s.inventoryRepo.DisableMenuItems(ctx, tx, menuItemIDs)
		}
	}

	if len(outOfStock) > 0 {
		s.redisStore.FlushByPattern(ctx, "menu:all:*")
		s.redisStore.FlushByPattern(ctx, "menu:categories:*")
	}

	return nil
}
