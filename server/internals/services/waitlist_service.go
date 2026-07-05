package services

import (
	"context"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
)

type WaitlistService struct {
	waitlistRepo repositories.WaitlistRepositoryInterface
}

func NewWaitlistService(waitlistRepo repositories.WaitlistRepositoryInterface) *WaitlistService {
	return &WaitlistService{waitlistRepo: waitlistRepo}
}

func (s *WaitlistService) JoinWaitlist(ctx context.Context, userID uint, req *dto.JoinWaitlistRequest) (*dto.WaitlistResponse, error) {
	// validate first
	if req.Date == "" {
		return nil, domain.ErrInvalidDate
	}
	if req.PartySize <= 0 {
		return nil, domain.ErrInvalidTableCapacity
	}
	if !utils.IsValidTimeSlots(req.TimeSlot) {
		return nil, domain.ErrInvalidTimeSlot
	}

	// then parse
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, domain.ErrInvalidDate
	}

	parsedTimeSlot, err := utils.ParseToDataTypesTime(req.TimeSlot)
	if err != nil {
		return nil, domain.ErrInvalidTimeSlot
	}

	// check if table is available, no need for waitlist
	count, err := s.waitlistRepo.CountAvailableTables(ctx, req.Date, req.TimeSlot, req.PartySize)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, domain.ErrTableAvailable
	}

	// check if already on waitlist
	_, err = s.waitlistRepo.GetByUserDateSlot(ctx, userID, req.Date, req.TimeSlot)
	if err == nil {
		return nil, domain.ErrAlreadyOnWaitlist
	}

	// join waitlist
	waitlist := models.Waitlist{
		UserID:    userID,
		TimeSlot:  parsedTimeSlot,
		Date:      parsedDate,
		PartySize: req.PartySize,
		Status:    models.WaitlistStatusWaiting,
		CreatedAt: time.Now(),
	}

	if err := s.waitlistRepo.Create(ctx, &waitlist); err != nil {
		return nil, err
	}

	return &dto.WaitlistResponse{
		ID:         waitlist.ID,
		UserID:     waitlist.UserID,
		Date:       waitlist.Date.Format("2006-01-02"),
		TimeSlot:   utils.FormatDataTypesTime(waitlist.TimeSlot),
		PartySize:  waitlist.PartySize,
		Status:     string(models.WaitlistStatusWaiting),
		NotifiedAt: waitlist.NotifiedAt,
		CreatedAt:  waitlist.CreatedAt,
	}, nil
}

func (s *WaitlistService) GetWaitlistPosition(ctx context.Context, userID uint, date, timeSlot string) (int, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	parsedTimeSlot, err := utils.ParseToDataTypesTime(timeSlot)
	if err != nil {
		return 0, domain.ErrInvalidTimeSlot
	}

	waitlist, err := s.waitlistRepo.GetByUserDateSlot(ctx, userID, date, timeSlot)
	if err != nil {
		return 0, domain.ErrNotFound
	}

	position, err := s.waitlistRepo.GetPosition(ctx, date, parsedTimeSlot, waitlist.CreatedAt)
	if err != nil {
		return 0, err
	}

	return int(position) + 1, nil
}
