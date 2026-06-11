package server

import (
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateTakeoutOrderHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	var req dto.CreateTakeoutOrderRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.orderService.CreateTakeoutOrder(userID, &req)
	if err != nil {
		switch err {
		case domain.ErrCartIsEmpty:
			utils.BadRequest(ctx, "Cart is empty", err)
		case domain.ErrCartNotFound:
			utils.NotFound(ctx, "Cart not found", err)
		case domain.ErrMenuNotAvailable:
			utils.BadRequest(ctx, "One or more menu items in the cart are not available", err)
		default:
			utils.InternalServerError(ctx, "Failed to create order", err)
			return
		}
	}

	utils.CreatedResponse(ctx, "Order created successfully", response)
}

func (s *Server) GetAllTakeoutOrdersHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	response, meta, err := s.orderService.GetAllTakeoutOrders(userID, page, pageSize)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to fetch orders", err)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Orders fetched successfully", response, *meta)
}

func (s *Server) GetTakeoutOrderHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid id", err)
		return
	}
	orderID := uint(id)

	response, err := s.orderService.GetTakeoutOrder(userID, orderID)
	if err != nil {
		switch err {
		case domain.ErrOrderNotFound:
			utils.NotFound(ctx, "Order not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to fetch order", err)
		}
	}

	utils.SuccessResponse(ctx, "Order fetched successfully", response)
}
