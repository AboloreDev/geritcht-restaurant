package server

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Join waitlist
// @Description Join the reservation waitlist for a preferred date, time slot, and party size when no suitable table is available.
// @Tags Waitlist
// @Accept json
// @Produce json
// @Param input body dto.JoinWaitlistRequest true "Waitlist request"
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=dto.WaitlistResponse} "Successfully joined the waitlist"
// @Failure 400 {object} utils.Response "Invalid request data, invalid date/time, table available, invalid table capacity, or user already on the waitlist"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /waitlist [post]
func (s *Server) JoinWaitlistHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	var waitlist *dto.JoinWaitlistRequest

	err := ctx.ShouldBindJSON(&waitlist)
	if err != nil {
		utils.BadRequest(ctx, "Invalid request data", err)
		return
	}

	response, err := s.waitlistService.JoinWaitlist(ctx.Request.Context(), userID, waitlist)
	if err != nil {
		switch err {
		case domain.ErrInvalidTableCapacity:
			utils.BadRequest(ctx, "Invalid table capacity", err)
		case domain.ErrAlreadyOnWaitlist:
			utils.BadRequest(ctx, "User already in waitlist", err)
		case domain.ErrTableAvailable:
			utils.BadRequest(ctx, "Table is available", err)
		case domain.ErrInvalidTimeSlot:
			utils.BadRequest(ctx, "Invalid time slot", err)
		case domain.ErrInvalidDate:
			utils.BadRequest(ctx, "Invalid date", err)
		default:
			utils.InternalServerError(ctx, "Error joining waitlist", err)
		}
		return
	}

	utils.CreatedResponse(ctx, "Joined waitlist success", response)
}

// @Summary Get waitlist position
// @Description Retrieve the authenticated user's current position on the waitlist for a specific date and time slot.
// @Tags Waitlist
// @Accept json
// @Produce json
// @Param date query string true "Reservation date (YYYY-MM-DD)"
// @Param time_slot query string true "Reservation time slot"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.WaitlistPositionResponse} "Waitlist position retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid request or user is not on the waitlist"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /waitlist/position [get]
func (s *Server) GetWaitlistPositionHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	date := ctx.Query("date")
	timeSlot := ctx.Query("time_slot")

	response, err := s.waitlistService.GetWaitlistPosition(ctx.Request.Context(), userID, date, timeSlot)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			utils.BadRequest(ctx, "User not found on waitlist", err)
		default:
			utils.InternalServerError(ctx, "Error joining waitlist", err)
		}
		return
	}

	utils.SuccessResponse(ctx, "success", response)
}
