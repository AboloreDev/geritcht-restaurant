package server

import (
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

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
		default:
			utils.InternalServerError(ctx, "Internal server error", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "Check In successfully", response)
}

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
