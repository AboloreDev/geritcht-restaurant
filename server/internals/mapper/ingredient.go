package mapper

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
)

func ConvertToInventoryAlertResponse(lowStockIngredients []models.Ingredient, outOfStockMenuItems []models.Menu,
) *dto.InventoryAlertResponse {
	lowStock := make([]dto.IngredientResponse, 0, len(lowStockIngredients))
	for _, ing := range lowStockIngredients {
		lowStock = append(lowStock, dto.IngredientResponse{
			ID:           ing.ID,
			Name:         ing.Name,
			CurrentStock: ing.CurrentStock,
			Unit:         ing.Unit,
		})
	}

	outOfStock := make([]dto.MenuResponse, 0, len(outOfStockMenuItems))
	for _, menu := range outOfStockMenuItems {
		images := []dto.MenuImageResponse{}
		if len(menu.Images) > 0 {
			images = append(images, dto.MenuImageResponse{
				ID:      menu.Images[0].ID,
				URL:     menu.Images[0].URL,
				AltText: menu.Images[0].AltText,
			})
		}
		outOfStock = append(outOfStock, dto.MenuResponse{
			ID:          menu.ID,
			Name:        menu.Name,
			IsAvailable: menu.IsAvailable,
			Images:      images,
		})
	}

	return &dto.InventoryAlertResponse{
		LowStockIngredients: lowStock,
		OutOfStockItems:     outOfStock,
		TotalLowStock:       len(lowStock),
		TotalOutOfStock:     len(outOfStock),
	}
}
