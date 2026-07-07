package server

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

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
