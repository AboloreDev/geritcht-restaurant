package services

import (
	"context"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/events"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/rs/zerolog"
)

type ReminderWorker struct {
	reminderRepo repositories.ReservationReminderInterface
	redisStore   interfaces.Cacher
	publisher    interfaces.Publisher
}

func NewReminderWorker(
	reminderRepo repositories.ReservationReminderInterface,
	redisStore interfaces.Cacher,
	publisher interfaces.Publisher) *ReminderWorker {
	return &ReminderWorker{
		reminderRepo: reminderRepo,
		redisStore:   redisStore,
		publisher:    publisher,
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
			w.processReminder(ctx, log)
		case <-heartbeat.C:
			log.Info().Msg("reminder worker is alive")
		}
	}
}

func (w *ReminderWorker) processReminder(ctx context.Context, log zerolog.Logger) {
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

	reservations, err := w.reminderRepo.GetAllUpcomingReservations(ctx, now, targetSlot)

	for _, reservation := range reservations {
		// Publish message
		err := w.publish(reservation)
		if err != nil {
			log.Error().Err(err).
				Uint("reservation_id", reservation.ID).
				Msg("failed to publish reminder email")
			continue
		}

		// Send reminder
		err = w.reminderRepo.UpdateReminderValue(ctx, &reservation)
		if err != nil {
			log.Error().Err(err).
				Uint("reservation_id", reservation.ID).
				Msg("failed to mark reminder as sent")
		}
	}
}

func (w *ReminderWorker) publish(reservation models.Reservation) error {
	w.publisher.PublishMessage(
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

	return nil
}
