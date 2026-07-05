package server

import (
	"io"
	"net/http"
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) InitilaisePaymentHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	var req dto.InitializePaymentRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.paymentService.InitialisePayment(ctx.Request.Context(), userID, &req)
	if err != nil {
		switch err {
		case domain.ErrOrderNotFound:
			utils.NotFound(ctx, "Order not found", err)
		case domain.ErrInvalidOrderStatus:
			utils.BadRequest(ctx, "Invalid order status", err)
		case domain.ErrOrderAlreadyPaid:
			utils.BadRequest(ctx, "Order already paid", err)
		case domain.ErrPaymentNotFound:
			utils.NotFound(ctx, "Payment not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to initialize payment", err)
			return
		}
	}

	utils.SuccessResponse(ctx, "Payment initialized successfully", response)
}

func (s *Server) VerifyPaymentHandler(ctx *gin.Context) {
	var req dto.VerifyPaymentRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.paymentService.VerifyPayment(ctx.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrPaymentNotFound:
			utils.NotFound(ctx, "Payment not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to verify payment", err)
			return
		}
	}

	utils.SuccessResponse(ctx, "Payment verified successfully", response)
}

func (s *Server) WebhookHandler(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		utils.BadRequest(ctx, "Failed to read body", err)
		return
	}

	signature := ctx.GetHeader("x-paystack-signature")
	if signature == "" {
		utils.UnAuthorized(ctx, "Missing signature", nil)
		return
	}

	err = s.paymentService.HandlePaystackWebhook(ctx.Request.Context(), body, signature)
	if err != nil {
		switch err {
		case domain.ErrInvalidSignature:
			utils.BadRequest(ctx, "Invalid signature", err)
		case domain.ErrPaymentNotFound:
			utils.NotFound(ctx, "Payment not found", err)
		case domain.ErrPaymentAmountMismatch:
			utils.BadRequest(ctx, "Payment amount mismatch", err)
		default:
			ctx.Status(http.StatusOK)
			return
		}
	}

	utils.SuccessResponse(ctx, "success", nil)
}

func (s *Server) GetAllPaymentHistory(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	response, meta, err := s.paymentService.GetAllPaymentHistory(ctx.Request.Context(), userID, page, pageSize)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to fetch payment history", err)
		return
	}

	utils.PaginatedSuccessResponse(ctx, "Payments fetched successfully", response, *meta)
}

func (s *Server) GetPaymentByReferenceHandler(ctx *gin.Context) {
	reference := ctx.Param("reference")

	response, err := s.paymentService.GetPaymentByReference(ctx.Request.Context(), reference)
	if err != nil {
		switch err {
		case domain.ErrPaymentNotFound:
			utils.NotFound(ctx, "Payment not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to fetch payment history", err)
		}
	}

	utils.SuccessResponse(ctx, "Payment fetched successfully", response)
}

func (s *Server) GetPaymentDetailsHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid id", err)
		return
	}
	paymentID := uint(id)

	response, err := s.paymentService.GetPaymentDetails(ctx.Request.Context(), paymentID)
	if err != nil {
		switch err {
		case domain.ErrPaymentNotFound:
			utils.NotFound(ctx, "Payment not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to fetch payment history", err)
		}
	}

	utils.SuccessResponse(ctx, "Payment fetched successfully", response)
}

func (s *Server) GetRefundDetailsHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid id", err)
		return
	}
	refundID := uint(id)

	response, err := s.paymentService.GetRefundDetails(ctx.Request.Context(), refundID)
	if err != nil {
		switch err {
		case domain.ErrRefundNotFound:
			utils.NotFound(ctx, "refund not found", err)
		default:
			utils.InternalServerError(ctx, "Failed to fetch refund", err)
		}
	}

	utils.SuccessResponse(ctx, "Refund fetched successfully", response)
}
