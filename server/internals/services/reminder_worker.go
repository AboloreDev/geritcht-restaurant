package services

import (
	"context"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/events"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ReminderWorker struct {
	db         *gorm.DB
	redisStore interfaces.Cacher
	publisher  interfaces.Publisher
}

func NewReminderWorker(db *gorm.DB, redisStore interfaces.Cacher, publisher interfaces.Publisher) *ReminderWorker {
	return &ReminderWorker{
		db:        db,
		redisStore: redisStore,
		publisher: publisher,
	}
}

func (w *ReminderWorker) StartReminderWorker(ctx context.Context, log zerolog.Logger) {
	ticker := time.NewTicker(5 * time.Minute)
	heartbeat := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	defer heartbeat.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("reminder worker stopped")
			return
		case <-ticker.C:
			w.processReminder(log) 
		case <-heartbeat.C:
			log.Info().Msg("reminder worker is alive")
		}
	}
}

func (w *ReminderWorker) processReminder(log zerolog.Logger) {
	var reservations []models.Reservation
	now := time.Now()

	// target = 30 minutes from now
	target := now.Add(30 * time.Minute)

	// convert to datatypes.Time for DB comparison
	targetSlot, err := utils.ParseToDataTypesTime(
		target.Format("15:04"), 
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse target time slot")
		return
	}

	w.db.Preload("Table").Preload("User").
		Where("date = ? AND status = ? AND time_slot = ? AND reminder_sent = ?",
			now.Format("2006-01-02"),
			models.ReservationStatusConfirmed,
			targetSlot,
			false,
		).Find(&reservations)

	for _, reservation := range reservations {
		err := w.publisher.PublishMessage(
			events.ChannelEmailReservationReminder,
			events.ReservationReminderPayload{
				FirstName: reservation.User.FirstName,
				TableName: reservation.Table.Name,
				Date:      reservation.Date.Format("2006-01-02"),
				TimeSlot:  utils.FormatDataTypesTime(reservation.TimeSlot),
				Email:     reservation.User.Email,
			},
			map[string]string{"Priority": "Important Mail"},
		)
		if err != nil {
			log.Error().Err(err).
				Uint("reservation_id", reservation.ID).
				Msg("failed to publish reminder email")
			continue
		}

		if err := w.db.Model(&reservation).
			Update("reminder_sent", true).Error; err != nil {
			log.Error().Err(err).
				Uint("reservation_id", reservation.ID).
				Msg("failed to mark reminder as sent")
		}
	}
}