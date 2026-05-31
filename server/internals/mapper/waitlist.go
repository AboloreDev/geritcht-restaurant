package mapper

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
)

func WaitlistResponse(waitlist *models.Waitlist) *dto.WaitlistResponse {
	return &dto.WaitlistResponse{
		ID:         waitlist.ID,
		UserID:     waitlist.UserID,
		Date:       waitlist.Date.Format("2000-05-08"),
		TimeSlot:   utils.FormatDataTypesTime(waitlist.TimeSlot),
		PartySize:  waitlist.PartySize,
		Status:     string(waitlist.Status),
		NotifiedAt: waitlist.NotifiedAt,
		CreatedAt:  waitlist.CreatedAt,
	}
}
