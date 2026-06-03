package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/events"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/mapper"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ReservationService struct {
	db         *gorm.DB
	redisStore interfaces.Cacher
	publisher  interfaces.Publisher
}

func NewReservationService(
	db *gorm.DB,
	redisStore interfaces.Cacher,
	publisher interfaces.Publisher) *ReservationService {
	return &ReservationService{
		db:         db,
		redisStore: redisStore,
		publisher:  publisher,
	}
}

func (s *ReservationService) CheckTableAvailability(req *dto.CheckAvailabilityRequest) (*dto.AvailabilityResponse, error) {
	if !utils.IsValidTimeSlots(req.TimeSlot) {
		return nil, domain.ErrInvalidTimeSlot
	}

	cacheKey := fmt.Sprintf("availability:%s:%s:%d",
		req.Date,
		req.TimeSlot,
		req.PartySize,
	)

	cached, err := s.redisStore.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var response dto.AvailabilityResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}
	}

	var allTables []models.Table
	err = s.db.Where("capacity >= ?", req.PartySize).
		Find(&allTables).Error
	if err != nil {
		return nil, err
	}

	var reservations []models.Reservation
	timeSlot, _ := utils.ParseToDataTypesTime(req.TimeSlot)
	err = s.db.Select("table_id").
		Where("date = ? AND time_slot = ? AND status NOT IN ?",
			req.Date,
			timeSlot,
			[]string{"cancelled", "no_show"},
		).
		Find(&reservations).Error
	if err != nil {
		return nil, err
	}

	bookedTableIDs := make(map[uint]bool, len(reservations))
	for _, r := range reservations {
		bookedTableIDs[r.TableID] = true
	}

	tables := make([]dto.TableAvailabilityResponse, 0, len(allTables))
	for _, table := range allTables {
		status := "available"
		if bookedTableIDs[table.ID] {
			status = "confirmed"
		}

		tables = append(tables, dto.TableAvailabilityResponse{
			ID:       table.ID,
			Name:     table.Name,
			Capacity: table.Capacity,
			Location: table.Location,
			Status:   status,
		})
	}

	response := &dto.AvailabilityResponse{
		Date:      req.Date,
		TimeSlot:  req.TimeSlot,
		PartySize: req.PartySize,
		Tables:    tables,
	}

	// short TTL.... availability changes when someone books
	data, _ := json.Marshal(response)
	s.redisStore.Set(ctx, cacheKey, string(data), 30*time.Second)

	return response, nil
}

func (s *ReservationService) CreateReservation(req *dto.CreateReservationRequest, userID uint) (*dto.ReservationResponse, error) {
	var table models.Table
	var reservation models.Reservation
	const defaultHoldTTL = 2 * time.Minute

	if !utils.IsValidTimeSlots(req.TimeSlot) {
		return nil, domain.ErrInvalidTimeSlot
	}

	if req.Date < time.Now().Format("2006-01-02") {
		return nil, domain.ErrPastDates
	}

	if req.PartySize <= 0 {
		return nil, domain.ErrInvalidTableCapacity
	}

	err := s.db.Where("id = ? AND capacity >= ?", req.TableID, req.PartySize).First(&table).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	lockKey := fmt.Sprintf("lock:reservation:%d:%s:%s",
		req.TableID, req.Date, req.TimeSlot)

	lockValue, err := json.Marshal(userID)
	if err != nil {
		return nil, err
	}

	err = s.redisStore.Hold(ctx, lockKey, string(lockValue), redis.SetArgs{
		Mode: "NX",
		TTL:  defaultHoldTTL,
	})

	if err != nil {
		return nil, domain.ErrTableAlreadyBooked
	}

	defer s.redisStore.Delete(ctx, lockKey)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		time.Sleep(5 * time.Second)
		err := tx.Raw("SELECT * FROM tables WHERE id = ? FOR UPDATE", req.TableID).
			Scan(&table).Error
		if err != nil {
			return err
		}

		var count int64
		tx.Model(&models.Reservation{}).
			Where("table_id = ? AND date = ? AND time_slot = ? AND status NOT IN ?",
				req.TableID, req.Date, req.TimeSlot,
				[]string{"cancelled", "no_show"},
			).Count(&count)

		if count > 0 {
			return domain.ErrTableAlreadyBooked
		}

		parsedDate, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return err
		}
		parsedTimeSlot, err := utils.ParseToDataTypesTime(req.TimeSlot)
		if err != nil {
			return err
		}

		reservation = models.Reservation{
			UserID:          userID,
			TableID:         table.ID,
			Date:            parsedDate,
			TimeSlot:        parsedTimeSlot,
			PartySize:       req.PartySize,
			Status:          models.ReservationStatusConfirmed,
			SpecialRequests: req.SpecialRequests,
		}

		return tx.Create(&reservation).Error
	})
	if err != nil {
		return nil, err
	}

	s.redisStore.FlushByPattern(ctx,
		fmt.Sprintf("availability:%s:%s:*", req.Date, req.TimeSlot),
	)

	holdKey := fmt.Sprintf("hold:table:%d:%s:%s",
		req.TableID, req.Date, req.TimeSlot,
	)
	s.redisStore.Delete(ctx, holdKey)

	if err := s.db.Preload("User").Preload("Table").
		First(&reservation, reservation.ID).Error; err != nil {
		return nil, err
	}

	err = s.publisher.PublishMessage(
		events.ChannelEmailReservationConfirm,
		events.ReservationConfirmPayload{
			Email:     reservation.User.Email,
			FirstName: reservation.User.FirstName,
			Date:      reservation.Date.Format("2006-01-02"),
			TimeSlot:  utils.FormatDataTypesTime(reservation.TimeSlot),
			TableName: reservation.Table.Name,
		},
		map[string]string{"Priority": "Important Mail"},
	)
	if err != nil {
		log.Printf("Failed to put messages in queue: %v", err)
		return nil, err
	}
	return mapper.ReservationResponse(&reservation), nil
}

func (s *ReservationService) GetAllUserReservations(userID uint, req *dto.ReservationFilterRequest) (*dto.ReservationListResponse, error) {
	var reservations []models.Reservation
	var count int64

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	offset := utils.Pagination(req.Page, req.PageSize)

	query := s.db.Preload("User").Preload("Table").
		Where("user_id = ?", userID)

	if req.Date != "" {
		query = query.Where("date = ?", req.Date)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	query = query.Order("date ASC, time_slot ASC")
	query.Model(&models.Reservation{}).Count(&count)

	err := query.Offset(offset).Limit(req.PageSize).Find(&reservations).Error
	if err != nil {
		return nil, err
	}

	response := make([]dto.ReservationResponse, 0, len(reservations))
	for _, rsv := range reservations {
		response = append(response, *mapper.ReservationResponse(&rsv))
	}

	totalPages := int((count + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &dto.ReservationListResponse{
		Reservations: response,
		Total:        count,
		Page:         req.Page,
		PageSize:     req.PageSize,
		TotalPages:   totalPages,
	}, nil
}

func (s *ReservationService) GetUserReservation(userID uint, reservationID uint) (*dto.ReservationResponse, error) {
	var reservation models.Reservation

	err := s.db.Preload("User").Preload("Table").
		Where("id = ? AND user_id = ?", reservationID, userID).
		First(&reservation).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, domain.ErrUserNotFound
	}

	if reservation.UserID != userID {
		return nil, domain.ErrForbidden
	}

	return mapper.ReservationResponse(&reservation), nil
}

func (s *ReservationService) GetAllReservations(req *dto.ReservationFilterRequest) (*dto.ReservationListResponse, error) {
	var reservation []models.Reservation
	var count int64

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	offset := utils.Pagination(req.Page, req.PageSize)

	query := s.db.Preload("User").Preload("Table")

	if req.Date != "" {
		query = query.Where("date = ?", req.Date)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	query = query.Order("time_slot ASC")

	query.Model(models.Reservation{}).Count(&count)

	err := query.Offset(offset).Limit(req.PageSize).Find(&reservation).Error
	if err != nil {
		return nil, err
	}

	response := make([]dto.ReservationResponse, 0, len(reservation))

	for _, rsv := range reservation {
		response = append(response, *mapper.ReservationResponse(&rsv))
	}

	totalPages := int((count + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &dto.ReservationListResponse{
		Reservations: response,
		Total:        count,
		Page:         req.Page,
		PageSize:     req.PageSize,
		TotalPages:   totalPages,
	}, nil

}

func (s *ReservationService) GetTodayReservations(req *dto.ReservationFilterRequest) (*dto.ReservationListResponse, error) {
	var reservation []models.Reservation
	var count int64

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	req.Date = time.Now().Format("2006-01-02")

	offset := utils.Pagination(req.Page, req.PageSize)

	query := s.db.Preload("User").Preload("Table")

	if req.Date != "" {
		query = query.Where("date = ?", req.Date)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	query = query.Order("time_slot ASC")

	query.Model(&models.Reservation{}).Count(&count)

	err := query.Offset(offset).Limit(req.PageSize).Find(&reservation).Error
	if err != nil {
		return nil, err
	}

	response := make([]dto.ReservationResponse, 0, len(reservation))

	for _, rsv := range reservation {
		response = append(response, *mapper.ReservationResponse(&rsv))
	}

	totalPages := int((count + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &dto.ReservationListResponse{
		Reservations: response,
		Total:        count,
		Page:         req.Page,
		PageSize:     req.PageSize,
		TotalPages:   totalPages,
	}, nil

}

func (s *ReservationService) CheckInReservation(reservationID uint, userID uint) (*dto.ReservationResponse, error) {
	var reservation models.Reservation
	var user models.User

	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, domain.ErrForbidden
	}

	err = s.db.Where("id = ? AND status = ? ", reservationID, models.ReservationStatusConfirmed).
		First(&reservation).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if reservation.Status == models.ReservationStatusCheckedIn {
		return nil, domain.ErrAlreadyCheckedIn
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&reservation).Where("id = ?", reservationID).Updates(map[string]interface{}{
			"status":        models.ReservationStatusCheckedIn,
			"checked_in_at": time.Now(),
		}).Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Table{}).Where("id = ?", reservation.TableID).
			Update("status", models.TableStatusOccupied).Error
		if err != nil {
			return err
		}

		tx.Save(&reservation)
		tx.Save(&reservation.Table)

		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := s.db.Preload("User").Preload("Table").
		First(&reservation, reservation.ID).Error; err != nil {
		return nil, err
	}

	err = s.publisher.PublishMessage(
		events.ChannelEmailReservationCheckedIn,
		events.ReservationCheckInPayload{
			Email:     reservation.User.Email,
			FirstName: reservation.User.FirstName,
			Date:      reservation.Date.Format("2006-01-02"),
			TimeSlot:  utils.FormatDataTypesTime(reservation.TimeSlot),
			TableName: reservation.Table.Name,
		},
		map[string]string{"Priority": "Important Mail"},
	)
	if err != nil {
		log.Printf("Failed to put messages in queue: %v", err)
		return nil, err
	}

	return &dto.ReservationResponse{
		ID:      reservation.ID,
		UserID:  reservation.UserID,
		TableID: reservation.TableID,
		Status:  string(models.ReservationStatusCheckedIn),
		Table: dto.TableResponse{
			ID:       reservation.Table.ID,
			Status:   string(models.TableStatusOccupied),
			Capacity: reservation.Table.Capacity,
		},
		CheckedInAt: reservation.CheckedInAt,
	}, nil
}

func (s *ReservationService) CancelRservation(userID uint, reservationID uint) (*dto.ReservationResponse, error) {
	var reservation models.Reservation
	today := time.Now().Truncate(24 * time.Hour)
	reservationDate := reservation.Date.Truncate(24 * time.Hour)
	var waitlist models.Waitlist

	err := s.db.Where("id = ? AND user_id = ?", reservationID, userID).
		First(&reservation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, domain.ErrUserNotFound
	}

	if reservation.UserID != userID {
		return nil, domain.ErrForbidden
	}

	if reservation.Status != models.ReservationStatusPending &&
		reservation.Status != models.ReservationStatusConfirmed {
		return nil, domain.ErrCannotCancel
	}

	if reservation.Status == models.ReservationStatusCancelled {
		return nil, domain.ErrAlreadyCancelled
	}

	if !reservationDate.After(today) {
		return nil, domain.ErrCannotCancel
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&reservation).
			Where("id = ? AND user_id = ?", reservationID, userID).
			Updates(map[string]interface{}{
				"status": models.ReservationStatusCancelled,
			}).Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Table{}).Where("id = ?", reservation.TableID).
			Update("status", models.TableStatusAvailable).Error
		if err != nil {
			return err
		}

		err = tx.Preload("User").
			Where("date = ? AND time_slot = ? AND party_size = ? AND status = ?",
				reservation.Date, reservation.TimeSlot, reservation.PartySize, models.WaitlistStatusWaiting).
			Order("created_at ASC").
			First(&waitlist).Error

		if err != nil {
			return nil
		}

		err = tx.Model(&waitlist).
			Updates(map[string]interface{}{
				"status":      models.WaitlistStatusNotified,
				"notified_at": time.Now(),
				"expires_at":  time.Now().Add(10 * time.Minute),
			}).Error

		if err != nil {
			return nil
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.redisStore.FlushByPattern(ctx,
		fmt.Sprintf("availability:%s:%s:*", reservation.Date, reservation.TimeSlot),
	)

	if err := s.db.Preload("User").Preload("Waitlist").
		First(&reservation, reservation.ID).Error; err != nil {
		return nil, err
	}

	err = s.publisher.PublishMessage(
		events.ChannelEmailWaitlistNotification,
		events.WaitlistNotificationPayload{
			Email:     waitlist.User.Email,
			FirstName: waitlist.User.FirstName,
			Date:      waitlist.Date.Format("2006-01-02"),
			TimeSlot:  utils.FormatDataTypesTime(reservation.TimeSlot),
			TableName: reservation.Table.Name,
		},
		map[string]string{"Priority": "Important Mail"},
	)
	if err != nil {
		log.Printf("Failed to put messages in queue: %v", err)
		return nil, err
	}

	if err := s.db.Preload("User").Preload("Table").
		First(&reservation, reservation.ID).Error; err != nil {
		return nil, err
	}

	err = s.publisher.PublishMessage(
		events.ChannelEmailReservationCancelled,
		events.ReservationCancelledPayload{
			Email:     reservation.User.Email,
			FirstName: reservation.User.FirstName,
			Date:      reservation.Date.Format("2006-01-02"),
			TimeSlot:  utils.FormatDataTypesTime(reservation.TimeSlot),
			TableName: reservation.Table.Name,
		},
		map[string]string{"Priority": "Important Mail"},
	)
	if err != nil {
		log.Printf("Failed to put messages in queue: %v", err)
		return nil, err
	}

	return &dto.ReservationResponse{
		ID:      reservation.ID,
		UserID:  reservation.UserID,
		TableID: reservation.TableID,
		Status:  string(models.ReservationStatusCancelled),
		Table: dto.TableResponse{
			ID:       reservation.Table.ID,
			Status:   string(models.TableStatusOccupied),
			Capacity: reservation.Table.Capacity,
		},
		CheckedInAt: reservation.CheckedInAt,
	}, nil

}
