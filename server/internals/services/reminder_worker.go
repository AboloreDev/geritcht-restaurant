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

const (
	reminderLeadTime = 20 * time.Minute // 30 minutes before the reservation time
	reminderInterval = 5 * time.Minute
	windowPadding    = 2 * time.Minute // 2 minutes window padding
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
	heartbeat := time.NewTicker(1 * time.Minute)
	defer heartbeat.Stop()

	now := time.Now()
	nextRun := now.Truncate(reminderInterval).Add(reminderInterval)
	initialDelay := time.Until(nextRun)

	log.Info().
		Dur("initial_delay", initialDelay).
		Time("first_run", nextRun).
		Msg("reminder worker aligning to clock")

	timer := time.NewTimer(initialDelay)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("reminder worker stopped")
			return
		case <-timer.C:
			w.processReminder(ctx, log)
			timer.Reset(reminderInterval)
		case <-heartbeat.C:
			log.Info().Msg("reminder worker is alive")
		}
	}
}

func (w *ReminderWorker) processReminder(ctx context.Context, log zerolog.Logger) {
	now := time.Now()

	// target = 30 minutes from now
	target := now.Add(reminderLeadTime)

	// convert to datatypes.Time for DB comparison
	windowStart := target.Add(-windowPadding).Format("15:04:05")
	windowEnd := target.Add(windowPadding).Format("15:04:05")

	reservations, err := w.reminderRepo.GetAllUpcomingReservations(ctx, now, windowStart, windowEnd)
	if err != nil {
		log.Error().Err(err).Msg("failed to get upcoming reservations")
		return
	}

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

	return err
}
