package server

import (
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Add item to cart
// @Description Add a menu item to the current user's shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param input body dto.AddToCartRequest true "Menu item and quantity to add"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Item added to cart successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Menu item not found"
// @Router /cart [post]
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

// @Summary Update cart item
// @Description Update the quantity of a specific item in the current user's shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart Item ID"
// @Param input body dto.UpdateCartItemRequest true "Updated cart item details"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Cart updated successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Cart item not found or menu item not found"
// @Router /cart/{id} [patch]
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
		return
	}

	utils.SuccessResponse(ctx, "cart updated successfully", cartResponse)
}

// @Summary Remove item from cart
// @Description Remove a specific item from the current user's shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart Item ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response "Item removed from cart successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Cart item not found"
// @Router /cart/{id} [delete]
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
		return
	}

	utils.SuccessResponse(ctx, "Item removed from cart successfully", nil)
}

// @Summary Clear user's cart
// @Description Remove all items from the current user's shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response "Cart cleared successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Cart not found"
// @Router /cart/clear [delete]
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
		return
	}

	utils.SuccessResponse(ctx, "Cart cleared successfully", nil)
}

// @Summary Get user's cart
// @Description Retrieve current user's shopping cart with all items
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Cart retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Cart not found"
// @Router /cart [get]
func (s *Server) GetUserCart(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	response, err := s.cartServices.GetUserCart(ctx.Request.Context(), userID)
	if err != nil {
		utils.BadRequest(ctx, "Failed to get cart", err)
		return
	}

	utils.SuccessResponse(ctx, "Cart fetched successfully", response)
}
