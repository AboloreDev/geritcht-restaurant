package server

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateAllergenHandler(ctx *gin.Context) {
	var req dto.CreateAllergenRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	err = s.allergenServices.CreateAllergenServices(&req)
	if err != nil {
		switch err {
		case domain.ErrNameConflict:
			utils.ConflictResponse(ctx, "Allergen already exists", err)
		default:
			utils.InternalServerError(ctx, "Failed to create allergen", err)
		}
		return
	}

	utils.CreatedResponse(ctx, "Allergen created successfully", nil)
}

func (s *Server) CreateDietaryTagsHandler(ctx *gin.Context) {
	var req dto.CreateDietaryTagRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	err = s.dietaryTagsService.CreateDietaryTags(&req)
	if err != nil {
		switch err {
		case domain.ErrNameConflict:
			utils.ConflictResponse(ctx, "DietaryTags already exists", err)
		default:
			utils.InternalServerError(ctx, "Failed to create allergen", err)
		}
		return
	}

	utils.CreatedResponse(ctx, "DietaryTags created successfully", nil)
}
