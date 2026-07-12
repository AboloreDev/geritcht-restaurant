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

// @Summary Initialize payment
// @Description Initialize a payment for an existing order and return the payment authorization details.
// @Tags Payments
// @Accept json
// @Produce json
// @Param input body dto.InitializePaymentRequest true "Payment initialization request"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.PaymentResponse} "Payment initialized successfully"
// @Failure 400 {object} utils.Response "Invalid request data, invalid order status, or order already paid"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Order or payment not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /payments/initialize [post]
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
		return
	}

	utils.SuccessResponse(ctx, "Payment initialized successfully", response)
}

// @Summary Verify payment
// @Description Verify the status of a payment using its transaction reference.
// @Tags Payments
// @Accept json
// @Produce json
// @Param input body dto.VerifyPaymentRequest true "Payment verification request"
// @Success 200 {object} utils.Response{data=dto.PaymentResponse} "Payment verified successfully"
// @Failure 400 {object} utils.Response "Invalid request data or payment not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /payments/verify/{ref} [post]
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
		return
	}

	utils.SuccessResponse(ctx, "Payment verified successfully", response)
}

// @Summary Handle Paystack webhook
// @Description Internal endpoint to handle incoming Paystack webhook requests for payment updates.
// @Tags Payments
// @Accept json
// @Produce json

// @Success 200 "Webhook processed successfully"
// @Failure 400 {object} utils.Response "Invalid signature or payment not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /payments/webhook [post]
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
		return
	}

	utils.SuccessResponse(ctx, "success", nil)
}

// @Summary Get payment history
// @Description Retrieve a paginated list of all payments made by the authenticated user.
// @Tags Payments
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Security BearerAuth
// @Success 200 {object} utils.PaginatedResponse{data=[]dto.PaymentResponse} "Payment history retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /payments/history [get]
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

// @Summary Get payment by reference
// @Description Retrieve the details of a payment using its unique payment reference.
// @Tags Payments
// @Accept json
// @Produce json
// @Param reference path string true "Payment reference"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.PaymentResponse} "Payment retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Payment not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /payments/ref/{reference} [get]
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
		return
	}

	utils.SuccessResponse(ctx, "Payment fetched successfully", response)
}

// @Summary Get payment details
// @Description Retrieve the details of a payment using its unique ID.
// @Tags Payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.PaymentResponse} "Payment retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Payment not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /payments/{id} [get]
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
		return
	}

	utils.SuccessResponse(ctx, "Payment fetched successfully", response)
}

// @Summary Get refund details
// @Description Retrieve the details of a refund using its unique ID.
// @Tags Payments
// @Accept json
// @Produce json
// @Param id path string true "Refund ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.RefundResponse} "Refund retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Refund not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /refunds/{id} [get]
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
		return
	}

	utils.SuccessResponse(ctx, "Refund fetched successfully", response)
}
