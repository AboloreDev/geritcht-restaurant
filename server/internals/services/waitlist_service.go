package services

import (
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/gorm"
)

type WaitlistService struct {
	db *gorm.DB
}

func NewWaitlistService(db *gorm.DB) *WaitlistService {
	return &WaitlistService{db: db}
}

func (s *WaitlistService) JoinWaitlist(userID uint, req *dto.JoinWaitlistRequest) (*dto.WaitlistResponse, error) {
	var count int64

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, err
	}

	parsedTimeSlot, err := utils.ParseToDataTypesTime(req.TimeSlot)
	if err != nil {
		return nil, err
	}

	if req.PartySize <= 0 {
		return nil, domain.ErrInvalidTableCapacity
	}

	if req.Date == "" {
		return nil, domain.ErrInvalidDate
	}

	if !utils.IsValidTimeSlots(req.TimeSlot) {
		return nil, domain.ErrInvalidTimeSlot
	}

	err = s.db.Model(&models.Table{}).
		Where("capacity >= ?", req.PartySize).
		Where("status = ?", models.TableStatusAvailable).
		Where("id NOT IN (?)",
			s.db.Model(&models.Reservation{}).
				Select("table_id").
				Where("date = ? AND time_slot = ?", req.Date, req.TimeSlot).
				Where("status NOT IN ?", []string{
					"cancelled",
					"no_show",
				}),
		).
		Count(&count).Error

	if err == nil {
		return nil, domain.ErrAlreadyOnWaitlist
	}

	if count > 0 {
		return nil, domain.ErrTableAvailable
	}

	err = s.db.Where("user_id = ? AND date = ? AND time_slot = ?", userID, req.Date, req.TimeSlot).
		First(&models.Waitlist{}).Error

	if err != nil {
		return nil, err
	}

	waitlist := models.Waitlist{
		UserID:    userID,
		TimeSlot:  parsedTimeSlot,
		Date:      parsedDate,
		PartySize: req.PartySize,
		Status:    models.WaitlistStatusWaiting,
		CreatedAt: time.Now(),
	}

	err = s.db.Create(&waitlist).Error
	if err != nil {
		return nil, err
	}

	return &dto.WaitlistResponse{
		ID:         waitlist.ID,
		UserID:     waitlist.UserID,
		Date:       waitlist.Date.Format("2000-05-08"),
		TimeSlot:   utils.FormatDataTypesTime(waitlist.TimeSlot),
		PartySize:  waitlist.PartySize,
		Status:     string(models.WaitlistStatusWaiting),
		NotifiedAt: waitlist.NotifiedAt,
		CreatedAt:  waitlist.CreatedAt,
	}, nil
}

func (s *WaitlistService) GetWaitlistPosition(userID uint, date, timeSlot string) (int, error) {
	var waitlist models.Waitlist

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	parsedTimeSlot, err := utils.ParseToDataTypesTime(timeSlot)
	if err != nil {
		return 0, err
	}

	err = s.db.Where("user_id = ? AND date = ? AND time_slot = ? AND status = ?",
		userID, date, parsedTimeSlot, models.WaitlistStatusWaiting).
		First(&waitlist).Error
	if err != nil {
		return 0, domain.ErrNotFound
	}

	var position int64
	s.db.Model(&models.Waitlist{}).
		Where("date = ? AND time_slot = ? AND status = ? AND created_at < ?",
			date, timeSlot, models.WaitlistStatusWaiting, waitlist.CreatedAt).
		Count(&position)

	return int(position) + 1, nil
}
