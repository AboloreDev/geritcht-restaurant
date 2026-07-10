package services

import (
	"context"
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
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReservationService struct {
	db              *gorm.DB
	reservationRepo repositories.ReservationRepositoryInterface
	redisStore      interfaces.Cacher
	publisher       interfaces.Publisher
}

func NewReservationService(
	db *gorm.DB,
	reservationRepo repositories.ReservationRepositoryInterface,
	redisStore interfaces.Cacher,
	publisher interfaces.Publisher,
) *ReservationService {
	return &ReservationService{
		db:              db,
		reservationRepo: reservationRepo,
		redisStore:      redisStore,
		publisher:       publisher,
	}
}

func (s *ReservationService) CheckTableAvailability(ctx context.Context, req *dto.CheckAvailabilityRequest) (*dto.AvailabilityResponse, error) {
	if !utils.IsValidTimeSlots(req.TimeSlot) {
		return nil, domain.ErrInvalidTimeSlot
	}

	cacheKey := fmt.Sprintf("availability:%s:%s:%d", req.Date, req.TimeSlot, req.PartySize)

	cached, err := s.redisStore.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var response dto.AvailabilityResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}
	}

	allTables, err := s.reservationRepo.GetTablesByCapacity(ctx, req.PartySize)
	if err != nil {
		return nil, err
	}

	timeSlot, _ := utils.ParseToDataTypesTime(req.TimeSlot)
	reservations, err := s.reservationRepo.GetReservationsByDateAndSlot(ctx, req.Date, timeSlot)
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

	data, _ := json.Marshal(response)
	s.redisStore.Set(ctx, cacheKey, string(data), 30*time.Second)

	return response, nil
}

func (s *ReservationService) CreateReservation(ctx context.Context, req *dto.CreateReservationRequest, userID uint) (*dto.ReservationResponse, error) {
	const defaultHoldTTL = 2 * time.Minute

	// validate
	if !utils.IsValidTimeSlots(req.TimeSlot) {
		return nil, domain.ErrInvalidTimeSlot
	}
	if req.Date < time.Now().Format("2006-01-02") {
		return nil, domain.ErrPastDates
	}
	if req.PartySize <= 0 {
		return nil, domain.ErrInvalidTableCapacity
	}

	// check table exists with enough capacity
	_, err := s.reservationRepo.GetTableByIDAndCapacity(ctx, req.TableID, req.PartySize)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	// acquire distributed lock
	lockKey := fmt.Sprintf("lock:reservation:%d:%s:%s", req.TableID, req.Date, req.TimeSlot)
	lockValue, _ := json.Marshal(userID)

	if err := s.redisStore.Hold(ctx, lockKey, string(lockValue), redis.SetArgs{
		Mode: "NX",
		TTL:  defaultHoldTTL,
	}); err != nil {
		return nil, domain.ErrTableAlreadyBooked
	}
	defer s.redisStore.Delete(ctx, lockKey)

	// parse
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, domain.ErrInvalidDate
	}
	parsedTimeSlot, err := utils.ParseToDataTypesTime(req.TimeSlot)
	if err != nil {
		return nil, domain.ErrInvalidTimeSlot
	}

	var reservation models.Reservation

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// FOR UPDATE lock
		table, err := s.reservationRepo.LockTableForUpdate(ctx, tx, req.TableID)

		// double check no booking exists
		count, err := s.reservationRepo.CountByTableDateSlot(ctx, req.TableID, req.Date, parsedTimeSlot)
		if err != nil {
			return err
		}
		if count > 0 {
			return domain.ErrTableAlreadyBooked
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

		return s.reservationRepo.Create(ctx, tx, &reservation)
	})
	if err != nil {
		return nil, err
	}

	// invalidate availability cache
	s.redisStore.FlushByPattern(ctx, fmt.Sprintf("availability:%s:%s:*", req.Date, req.TimeSlot))

	// fetch full reservation with relations
	fullReservation, err := s.reservationRepo.GetByIDWithRelations(ctx, reservation.ID)
	if err != nil {
		return nil, err
	}

	// publish confirmation email
	if err := s.publisher.PublishMessage(
		events.ChannelEmailReservationConfirm,
		events.ReservationConfirmPayload{
			Email:     fullReservation.User.Email,
			FirstName: fullReservation.User.FirstName,
			Date:      fullReservation.Date.Format("2006-01-02"),
			TimeSlot:  utils.FormatDataTypesTime(fullReservation.TimeSlot),
			TableName: fullReservation.Table.Name,
		},
		map[string]string{"Priority": "Important Mail"},
	); err != nil {
		log.Printf("Failed to publish reservation confirmation: %v", err)
	}

	return mapper.ReservationResponse(fullReservation), nil
}

func (s *ReservationService) GetAllUserReservations(ctx context.Context, userID uint, req *dto.ReservationFilterRequest) (*dto.ReservationListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	reservations, count, err := s.reservationRepo.GetAllByUser(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	return s.buildReservationListResponse(reservations, count, req), nil
}

func (s *ReservationService) GetUserReservation(ctx context.Context, userID uint, reservationID uint) (*dto.ReservationResponse, error) {
	reservation, err := s.reservationRepo.GetByIDAndUser(ctx, reservationID, userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if reservation.UserID != userID {
		return nil, domain.ErrForbidden
	}

	fullReservation, err := s.reservationRepo.GetByIDWithRelations(ctx, reservationID)
	if err != nil {
		return nil, err
	}

	return mapper.ReservationResponse(fullReservation), nil
}

func (s *ReservationService) GetAllReservations(ctx context.Context, req *dto.ReservationFilterRequest) (*dto.ReservationListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	reservations, count, err := s.reservationRepo.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}

	return s.buildReservationListResponse(reservations, count, req), nil
}

func (s *ReservationService) GetTodayReservations(ctx context.Context, req *dto.ReservationFilterRequest) (*dto.ReservationListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	reservations, count, err := s.reservationRepo.GetTodayReservations(ctx, req)
	if err != nil {
		return nil, err
	}

	return s.buildReservationListResponse(reservations, count, req), nil
}

func (s *ReservationService) CheckInReservation(ctx context.Context, reservationID uint, userID uint) (*dto.ReservationResponse, error) {
	reservation, err := s.reservationRepo.GetByIDAndStatus(ctx, reservationID, models.ReservationStatusConfirmed)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if reservation.Status == models.ReservationStatusCheckedIn {
		return nil, domain.ErrAlreadyCheckedIn
	}

	// Prevent illegal checkin
	// only allow check in 5 min before the time
	if canCheckIn(reservation.Date, reservation.TimeSlot) != nil {
		return nil, domain.ErrCannotCheckIn
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.reservationRepo.UpdateStatus(ctx, tx, reservationID, map[string]interface{}{
			"status":        models.ReservationStatusCheckedIn,
			"checked_in_at": time.Now(),
		}); err != nil {
			return err
		}
		return s.reservationRepo.UpdateTableStatus(ctx, tx, reservation.TableID, models.TableStatusOccupied)
	})
	if err != nil {
		return nil, err
	}

	fullReservation, err := s.reservationRepo.GetByIDWithRelations(ctx, reservationID)
	if err != nil {
		return nil, err
	}

	if err := s.publisher.PublishMessage(
		events.ChannelEmailReservationCheckedIn,
		events.ReservationCheckInPayload{
			Email:     fullReservation.User.Email,
			FirstName: fullReservation.User.FirstName,
			Date:      fullReservation.Date.Format("2006-01-02"),
			TimeSlot:  utils.FormatDataTypesTime(fullReservation.TimeSlot),
			TableName: fullReservation.Table.Name,
		},
		map[string]string{"Priority": "Important Mail"},
	); err != nil {
		log.Printf("Failed to publish check-in event: %v", err)
	}

	return mapper.ReservationResponse(fullReservation), nil
}

func (s *ReservationService) CancelReservation(ctx context.Context, userID uint, reservationID uint) (*dto.ReservationResponse, error) {
	reservation, err := s.reservationRepo.GetByIDAndUser(ctx, reservationID, userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if reservation.UserID != userID {
		return nil, domain.ErrForbidden
	}

	if reservation.Status == models.ReservationStatusCancelled {
		return nil, domain.ErrAlreadyCancelled
	}

	if reservation.Status != models.ReservationStatusPending &&
		reservation.Status != models.ReservationStatusConfirmed {
		return nil, domain.ErrCannotCancel
	}

	today := time.Now().Truncate(24 * time.Hour)
	reservationDate := reservation.Date.Truncate(24 * time.Hour)
	if !reservationDate.After(today) {
		return nil, domain.ErrCannotCancel
	}

	var waitlist *models.Waitlist

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.reservationRepo.UpdateStatus(ctx, tx, reservationID, map[string]interface{}{
			"status": models.ReservationStatusCancelled,
		}); err != nil {
			return err
		}

		if err := s.reservationRepo.UpdateTableStatus(ctx, tx, reservation.TableID, models.TableStatusAvailable); err != nil {
			return err
		}

		// notify next waitlist person
		wl, err := s.reservationRepo.GetFirstWaitlistByDateSlot(ctx, tx,
			reservation.Date, reservation.TimeSlot, reservation.PartySize)
		if err != nil {
			return nil // no waitlist → ok
		}

		waitlist = wl

		return s.reservationRepo.UpdateWaitlistStatus(ctx, tx, wl, map[string]interface{}{
			"status":      models.WaitlistStatusNotified,
			"notified_at": time.Now(),
			"expires_at":  time.Now().Add(10 * time.Minute),
		})
	})
	if err != nil {
		return nil, err
	}

	s.redisStore.FlushByPattern(ctx,
		fmt.Sprintf("availability:%s:%s:*", reservation.Date.Format("2006-01-02"), reservation.TimeSlot),
	)

	fullReservation, err := s.reservationRepo.GetByIDWithRelations(ctx, reservationID)
	if err != nil {
		return nil, err
	}

	// notify waitlist person
	if waitlist != nil {
		s.publisher.PublishMessage(
			events.ChannelEmailWaitlistNotification,
			events.WaitlistNotificationPayload{
				Email:     waitlist.User.Email,
				FirstName: waitlist.User.FirstName,
				Date:      waitlist.Date.Format("2006-01-02"),
				TimeSlot:  utils.FormatDataTypesTime(reservation.TimeSlot),
				TableName: fullReservation.Table.Name,
			},
			map[string]string{"Priority": "Important Mail"},
		)
	}

	// notify cancelled user
	s.publisher.PublishMessage(
		events.ChannelEmailReservationCancelled,
		events.ReservationCancelledPayload{
			Email:     fullReservation.User.Email,
			FirstName: fullReservation.User.FirstName,
			Date:      fullReservation.Date.Format("2006-01-02"),
			TimeSlot:  utils.FormatDataTypesTime(fullReservation.TimeSlot),
			TableName: fullReservation.Table.Name,
		},
		map[string]string{"Priority": "Important Mail"},
	)

	return mapper.ReservationResponse(fullReservation), nil
}

//  Private helper

func (s *ReservationService) buildReservationListResponse(
	reservations []models.Reservation,
	count int64,
	req *dto.ReservationFilterRequest,
) *dto.ReservationListResponse {
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
	}
}

func canCheckIn(reservationDate time.Time, timeSlot datatypes.Time) error {
	now := time.Now()

	// Reservation must be for today — no early check-in on a future date
	if now.Format("2006-01-02") != reservationDate.Format("2006-01-02") {
		return fmt.Errorf("check-in is only available on the day of your reservation")
	}

	nowAsTime := utils.TimeToDataTypesTime(now)

	const checkInWindow = int64(2 * time.Minute) // in nanoseconds
	earliestAllowed := int64(timeSlot) - checkInWindow

	if int64(nowAsTime) < earliestAllowed {
		return fmt.Errorf("check-in opens 2 minutes before your reservation time")
	}

	return nil
}
