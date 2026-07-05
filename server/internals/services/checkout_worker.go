package services

import (
	"context"
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
	"github.com/rs/zerolog"
	"gorm.io/datatypes"
)

type CheckoutWorker struct {
	redisStore   interfaces.Cacher
	checkoutRepo repositories.ReservationCheckoutInterface
}

func NewCheckoutWorker(

	redisStore interfaces.Cacher,
	checkoutRepo repositories.ReservationCheckoutInterface) *CheckoutWorker {
	return &CheckoutWorker{
		redisStore:   redisStore,
		checkoutRepo: checkoutRepo,
	}
}

var reservationDuration = 70 * time.Minute

func (w *CheckoutWorker) Start(appCtx context.Context, log zerolog.Logger) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.processCheckOut(ctx, log)
			log.Info().Msg("Running checkout worker")
		case <-ctx.Done():
			log.Info().Msg("Shutting down checkout worker")
			return
		}
	}
}

func DataTypesTimeToHourMinute(t datatypes.Time) (int, int) {
	totalSeconds := int64(t) / 1e9
	hours := int(totalSeconds / 3600)
	minutes := int((totalSeconds % 3600) / 60)
	return hours, minutes
}

func (w *CheckoutWorker) processCheckOut(ctx context.Context, log zerolog.Logger) {
	now := time.Now()
	reservations, err := w.checkoutRepo.GetAllRservations(ctx, now)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching checked-in reservations")
	}

	for _, reservation := range reservations {
		hours, minutes := DataTypesTimeToHourMinute(reservation.TimeSlot)

		if err != nil {
			log.Error().Err(err).Msg("Error parsing time slot")
			continue
		}

		slotTime := time.Date(
			now.Year(), now.Month(), now.Day(),
			hours, minutes, 0, 0, now.Location(),
		)

		endTime := slotTime.Add(reservationDuration)

		if now.After(endTime) {
			w.checkoutHandler(ctx, reservation, log)
		}
	}
}

func (w *CheckoutWorker) checkoutHandler(ctx context.Context, reservation models.Reservation, log zerolog.Logger) {
	// Cheeckout
	err := w.checkoutRepo.Checkout(ctx, reservation)
	if err != nil {
		log.Error().Err(err).
			Uint("reservation_id", reservation.ID).
			Msg("failed to checkout reservation")
		return
	}

	// Delete cache
	w.deleteCache(ctx, reservation)

	log.Info().
		Uint("reservation_id", reservation.ID).
		Uint("table_id", reservation.TableID).
		Msg("reservation checked out automatically")
}

// Delete cache function
func (w *CheckoutWorker) deleteCache(ctx context.Context, reservation models.Reservation) {
	w.redisStore.Delete(ctx,
		fmt.Sprintf("table:item:%d", reservation.TableID),
	)
}
