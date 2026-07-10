package services

import (
	"context"
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/events"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type NoShowWorker struct {
	publisher        interfaces.Publisher
	redisStore       interfaces.Cacher
	noShowWorkerRepo repositories.ReservationNoShowInterface
}

func NewNoShowWorker(
	publisher interfaces.Publisher,
	redisStore interfaces.Cacher,
	noShowWorkerRepo repositories.ReservationNoShowInterface) *NoShowWorker {
	return &NoShowWorker{
		publisher:        publisher,
		redisStore:       redisStore,
		noShowWorkerRepo: noShowWorkerRepo,
	}
}

func (w *NoShowWorker) StartMarkNoShowWorker(ctx context.Context, log zerolog.Logger) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Shutting down MarkNoShow worker")
			return
		case <-ticker.C:
			w.processNoShow(ctx)
			log.Info().Msg("MarkNoShow worker executed")

		}
	}
}

func (w *NoShowWorker) processNoShow(ctx context.Context) {
	reservations, err := w.noShowWorkerRepo.GetAllReservations(ctx)
	if err != nil {
		return
	}

	for _, reservation := range reservations {
		w.handleReservation(ctx, reservation)
	}
}

func (w *NoShowWorker) handleReservation(ctx context.Context, reservation models.Reservation) {
	// Business logic (transaction)
	if err := w.noShowWorkerRepo.MarkReservationNoShow(ctx, &reservation); err != nil {
		return
	}
	// Publish Mail
	err := w.publishNoShowEmail(reservation)

	if err != nil {
		log.Error().Err(err).
			Uint("reservation_id", reservation.ID).
			Msg("failed to publish no show email")
		return
	}

	// Redis
	w.invalidateAvailabilityCache(ctx, reservation)
}

// Publish Email
func (w *NoShowWorker) publishNoShowEmail(reservation models.Reservation) error {
	err := w.publisher.PublishMessage(
		events.ChannelEmailReservationNoShow,
		events.ReservationNoShowPayload{
			FirstName: reservation.User.FirstName,
			Email:     reservation.User.Email,
			Date:      reservation.Date.Format("2006-01-02"),
			TimeSlot:  utils.FormatDataTypesTime(reservation.TimeSlot),
			TableName: reservation.Table.Name,
		},
		map[string]string{
			"Priority": "Important Mail",
		},
	)
	return err
}

// Redis Func
func (w *NoShowWorker) invalidateAvailabilityCache(ctx context.Context, reservation models.Reservation) {
	_ = w.redisStore.FlushByPattern(
		ctx,
		fmt.Sprintf(
			"availability:%s:%s:*",
			reservation.Date,
			reservation.TimeSlot,
		),
	)
}
