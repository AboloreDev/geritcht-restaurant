package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

type PaginatedMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Response
	Meta PaginatedMeta `json:"meta"`
}

func SuccessResponse(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func CreatedResponse(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(ctx *gin.Context, statusCode int, message string, err error) {
	response := Response{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	ctx.JSON(statusCode, response)
}

func BadRequest(ctx *gin.Context, message string, err error) {
	ErrorResponse(ctx, http.StatusBadRequest, message, err)
}

func UnAuthorized(ctx *gin.Context, message string, err error) {
	ErrorResponse(ctx, http.StatusUnauthorized, message, err)
}

func InternalServerError(ctx *gin.Context, message string, err error) {
	ErrorResponse(ctx, http.StatusInternalServerError, message, err)
}

func ConflictResponse(ctx *gin.Context, message string, err error) {
	ErrorResponse(ctx, http.StatusConflict, message, err)
}

func NotFound(ctx *gin.Context, message string, err error) {
	ErrorResponse(ctx, http.StatusNotFound, message, err)
}

func Forbidden(ctx *gin.Context, message string, err error) {
	ErrorResponse(ctx, http.StatusForbidden, message, err)
}

func PaginatedSuccessResponse(ctx *gin.Context, message string, data interface{}, meta PaginatedMeta) {
	ctx.JSON(http.StatusOK, PaginatedResponse{
		Response: Response{
			Success: true,
			Message: message,
			Data:    data,
		},
		Meta: meta,
	})
}

func TooManyRequests(ctx *gin.Context, message string, err error) {
	ErrorResponse(ctx, http.StatusTooManyRequests, message, err)
}
