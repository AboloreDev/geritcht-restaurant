package server

import (
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Create a reservation
// @Description Create a new table reservation for the authenticated user.
// @Tags Reservations
// @Accept json
// @Produce json
// @Param input body dto.CreateReservationRequest true "Reservation details"
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=dto.ReservationResponse} "Reservation created successfully"
// @Failure 400 {object} utils.Response "Invalid request data, invalid time slot, table unavailable, or invalid table capacity"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Table not found"
// @Failure 429 {object} utils.Response "Too many requests"
// @Failure 500 {object} utils.Response "Internal server error"
// @Header 201 {string} X-RateLimit-Limit "Maximum requests allowed"
// @Header 201 {string} X-RateLimit-Remaining "Remaining requests in the current window"
// @Header 429 {string} X-RateLimit-Limit "Maximum requests allowed"
// @Header 429 {string} X-RateLimit-Remaining "Remaining requests (0 when limited)"
// @Router /reservations [post]
func (s *Server) CreateReservationHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var reservation *dto.CreateReservationRequest

	err := ctx.ShouldBindJSON(&reservation)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.reservationServices.CreateReservation(ctx.Request.Context(), reservation, userID)
	if err != nil {
		switch err {
		case domain.ErrInvalidTimeSlot:
			utils.BadRequest(ctx, "Invalid time slot", err)
		case domain.ErrTableAlreadyBooked:
			utils.BadRequest(ctx, "Table already booked", err)
		case domain.ErrPastDates:
			utils.BadRequest(ctx, "You cant book for past dates", err)
		case domain.ErrInvalidTableCapacity:
			utils.BadRequest(ctx, "Invalid table capacity", err)
		case domain.ErrNotFound:
			utils.NotFound(ctx, "Table not found", err)
		default:
			utils.InternalServerError(ctx, "Internal server error", err)
		}
		return
	}

	utils.CreatedResponse(ctx, "Reservation created successfully", response)
}

// @Summary Check table availability
// @Description Check available tables for a specified date, time, and party size.
// @Tags Reservations
// @Accept json
// @Produce json
// @Param date query string true "Reservation date (YYYY-MM-DD)"
// @Param time query string true "Reservation time (HH:MM)"
// @Param guests query int true "Number of guests"
// @Success 200 {object} utils.Response{data=[]dto.TableResponse} "Availability retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid request data or time slot"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /availability [get]
func (s *Server) CheckAvailabilityHandler(ctx *gin.Context) {
	var availability dto.CheckAvailabilityRequest

	err := ctx.ShouldBindQuery(&availability)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.reservationServices.CheckTableAvailability(ctx.Request.Context(), &availability)
	if err != nil {
		switch err {
		case domain.ErrInvalidTimeSlot:
			utils.BadRequest(ctx, "Invalid time slot", err)
		default:
			utils.InternalServerError(ctx, "Internal server error", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Availability fetched successfully", response)
}

// @Summary Get my reservations
// @Description Retrieve all reservations belonging to the authenticated user.
// @Tags Reservations
// @Accept json
// @Produce json
// @Param status query string false "Reservation status"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]dto.ReservationResponse} "Reservations retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /reservations [get]
func (s *Server) GetAllUserReservationsHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var reservationFiler dto.ReservationFilterRequest

	err := ctx.ShouldBindQuery(&reservationFiler)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.reservationServices.GetAllUserReservations(ctx.Request.Context(), userID, &reservationFiler)
	if err != nil {
		switch err {
		case domain.ErrInvalidDate:
			utils.BadRequest(ctx, "Invalid date", err)
		default:
			utils.InternalServerError(ctx, "Internal server error", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "User Reservation fetched successfully", response)
}

// @Summary Get reservation
// @Description Retrieve a reservation by its ID for the authenticated user.
// @Tags Reservations
// @Accept json
// @Produce json
// @Param id path int true "Reservation ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.ReservationResponse} "Reservation retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid reservation ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Reservation not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /reservations/{id} [get]
func (s *Server) GetUserReservationHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid ID", err)
		return
	}
	reservationID := uint(id)

	response, err := s.reservationServices.GetUserReservation(ctx.Request.Context(), userID, reservationID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "Reservation not found", err)
		case domain.ErrUserNotFound:
			utils.NotFound(ctx, "User not found", err)
		case domain.ErrInvalidDate:
			utils.BadRequest(ctx, "Invalid date", err)
		default:
			utils.InternalServerError(ctx, "Internal server error", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "User Reservation fetched successfully", response)
}

// @Summary List all reservations
// @Description Retrieve all reservations. Admin access required.
// @Tags Reservations
// @Accept json
// @Produce json
// @Param status query string false "Reservation status"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]dto.ReservationResponse} "Reservations retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /reservations/admin [get]
func (s *Server) GetAllReservationsHandler(ctx *gin.Context) {
	var reservationFiler dto.ReservationFilterRequest

	err := ctx.ShouldBindQuery(&reservationFiler)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.reservationServices.GetAllReservations(ctx.Request.Context(), &reservationFiler)
	if err != nil {
		switch err {
		case domain.ErrInvalidDate:
			utils.BadRequest(ctx, "Invalid date", err)
		default:
			utils.InternalServerError(ctx, "Internal server error", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "User Reservation fetched successfully", response)
}

// @Summary Get today's reservations
// @Description Retrieve all reservations for today. Admin access required.
// @Tags Reservations
// @Accept json
// @Produce json
// @Param status query string false "Reservation status"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]dto.ReservationResponse} "Reservations retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /reservations/today [get]
func (s *Server) GetTodayReservationHandler(ctx *gin.Context) {
	var reservationFiler dto.ReservationFilterRequest

	err := ctx.ShouldBindQuery(&reservationFiler)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.reservationServices.GetTodayReservations(ctx.Request.Context(), &reservationFiler)
	if err != nil {
		switch err {
		case domain.ErrInvalidDate:
			utils.BadRequest(ctx, "Invalid date", err)
		default:
			utils.InternalServerError(ctx, "Internal server error", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Today Reservation fetched successfully", response)
}

// @Summary Check in reservation
// @Description Check in the authenticated user for an existing reservation.
// @Tags Reservations
// @Accept json
// @Produce json
// @Param id path int true "Reservation ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.ReservationResponse} "Checked in successfully"
// @Failure 400 {object} utils.Response "Cannot check in or reservation already checked in"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Reservation not found"
// @Failure 429 {object} utils.Response "Too many requests"
// @Failure 500 {object} utils.Response "Internal server error"
// @Header 200 {string} X-RateLimit-Limit "Maximum requests allowed"
// @Header 200 {string} X-RateLimit-Remaining "Remaining requests in the current window"
// @Header 429 {string} X-RateLimit-Limit "Maximum requests allowed"
// @Header 429 {string} X-RateLimit-Remaining "Remaining requests (0 when limited)"
// @Router /reservations/{id}/check-in [post]
func (s *Server) CheckInReservationHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid ID", err)
		return
	}
	reservationID := uint(id)

	response, err := s.reservationServices.CheckInReservation(ctx.Request.Context(), reservationID, userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "Reservation not found", err)
		case domain.ErrForbidden:
			utils.Forbidden(ctx, "Forbidden", err)
		case domain.ErrAlreadyCheckedIn:
			utils.BadRequest(ctx, "You are already checked in", err)
		case domain.ErrCannotCheckIn:
			utils.BadRequest(ctx, "You cannot checkin a reservation within 15minutes of the reserved time", err)
		default:
			utils.InternalServerError(ctx, "Internal server error", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Check In successfully", response)
}

// @Summary Cancel reservation
// @Description Cancel an existing reservation owned by the authenticated user.
// @Tags Reservations
// @Accept json
// @Produce json
// @Param id path int true "Reservation ID"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.ReservationResponse} "Reservation cancelled successfully"
// @Failure 400 {object} utils.Response "Reservation cannot be cancelled"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Reservation or user not found"
// @Failure 429 {object} utils.Response "Too many requests"
// @Failure 500 {object} utils.Response "Internal server error"
// @Header 200 {string} X-RateLimit-Limit "Maximum requests allowed"
// @Header 200 {string} X-RateLimit-Remaining "Remaining requests in the current window"
// @Header 429 {string} X-RateLimit-Limit "Maximum requests allowed"
// @Header 429 {string} X-RateLimit-Remaining "Remaining requests (0 when limited)"
// @Router /reservations/{id}/cancel [patch]
func (s *Server) CancelReservationHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid ID", err)
		return
	}
	reservationID := uint(id)

	response, err := s.reservationServices.CancelReservation(ctx.Request.Context(), userID, reservationID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.NotFound(ctx, "Reservation not found", err)
		case domain.ErrForbidden:
			utils.Forbidden(ctx, "Forbidden", err)
		case domain.ErrCannotCancel:
			utils.BadRequest(ctx, "You cannot cancel a reservation within 24hours of the reserved time", err)
		case domain.ErrUserNotFound:
			utils.NotFound(ctx, "User not found", err)
		default:
			utils.InternalServerError(ctx, "Internal server error", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Successfully cancelled", response)
}
