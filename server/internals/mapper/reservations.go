package mapper

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
)

// add to mapper
func ReservationResponse(reservation *models.Reservation) *dto.ReservationResponse {
	return &dto.ReservationResponse{
		ID:              reservation.ID,
		UserID:          reservation.UserID,
		TableID:         reservation.TableID,
		PartySize:       reservation.PartySize,
		Status:          string(reservation.Status),
		SpecialRequests: reservation.SpecialRequests,
		CheckedInAt:     reservation.CheckedInAt,
		CreatedAt:       reservation.CreatedAt,
		Date:            reservation.Date.Format("2006-01-02"),
		TimeSlot:        utils.FormatDataTypesTime(reservation.TimeSlot),
		User: &dto.UserResponse{
			ID:          reservation.User.ID,
			FirstName:   reservation.User.FirstName,
			LastName:    reservation.User.LastName,
			Role: string(reservation.User.Role),
			Email:       reservation.User.Email,
			PhoneNumber: reservation.User.PhoneNumber,
		},
		Table: dto.TableResponse{
			ID:       reservation.Table.ID,
			Name:     reservation.Table.Name,
			Capacity: reservation.Table.Capacity,
			Location: reservation.Table.Location,
			Status:   string(reservation.Table.Status),
		},
	}
}
