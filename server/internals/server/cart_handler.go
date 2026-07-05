package server

import (
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) AddToCartHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	var req dto.AddToCartRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	cartResponse, err := s.cartServices.AddItemToCart(ctx.Request.Context(), userID, &req)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to add item to cart", err)
		return
	}

	utils.SuccessResponse(ctx, "Item Added to cart successfully", cartResponse)
}

func (s *Server) UpdateCartItemHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Item ID", err)
		return
	}
	itemID := uint(id)

	var req dto.UpdateCartItemRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	cartResponse, err := s.cartServices.UpdateCartItem(ctx.Request.Context(), userID, itemID, &req)
	if err != nil {
		switch err {
		case domain.ErrCartItemNotFound:
			utils.NotFound(ctx, "Cart item not found", err)
		case domain.ErrMenuNotFound:
			utils.NotFound(ctx, "Menu item not found", err)
		default:
			utils.InternalServerError(ctx, "Something went wrong", err)
			return
		}
	}

	utils.SuccessResponse(ctx, "cart updated successfully", cartResponse)
}

func (s *Server) RemoveCartItemHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid Item ID", err)
		return
	}
	itemID := uint(id)

	err = s.cartServices.RemoveCartItem(ctx.Request.Context(), userID, itemID)
	if err != nil {
		switch err {
		case domain.ErrCartItemNotFound:
			utils.NotFound(ctx, "Cart item not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to add item to cart", err)
			return
		}
	}

	utils.SuccessResponse(ctx, "Item removed from cart successfully", nil)
}

func (s *Server) ClearCartHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	err := s.cartServices.ClearCart(ctx.Request.Context(), userID)
	if err != nil {
		switch err {
		case domain.ErrCartNotFound:
			utils.NotFound(ctx, "Cart not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to add item to cart", err)
			return
		}
	}

	utils.SuccessResponse(ctx, "Cart cleared successfully", nil)
}

func (s *Server) GetUserCart(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	response, err := s.cartServices.GetUserCart(ctx.Request.Context(), userID)
	if err != nil {
		utils.BadRequest(ctx, "Failed to get cart", err)
		return
	}

	utils.SuccessResponse(ctx, "Cart fetched successfully", response)
}
