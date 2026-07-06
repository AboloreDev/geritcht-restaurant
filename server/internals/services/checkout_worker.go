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

// reservationDuration is how long a reservation occupies a table before
// it's eligible for automatic checkout.
var reservationDuration = 5 * time.Minute // TODO: bump back to 70 * time.Minute after testing

func (w *CheckoutWorker) Start(appCtx context.Context, log zerolog.Logger) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.processCheckOut(appCtx, log)
			log.Info().Msg("running checkout worker")
		case <-appCtx.Done():
			log.Info().Msg("shutting down checkout worker")
			return
		}
	}
}

// DataTypesTimeToHourMinuteSecond converts a datatypes.Time (nanoseconds
// since midnight) into hour, minute, and second components.
func DataTypesTimeToHourMinuteSecond(t datatypes.Time) (int, int, int) {
	totalSeconds := int64(t) / 1e9
	hours := int(totalSeconds / 3600)
	minutes := int((totalSeconds % 3600) / 60)
	seconds := int(totalSeconds % 60)
	return hours, minutes, seconds
}

func (w *CheckoutWorker) processCheckOut(ctx context.Context, log zerolog.Logger) {
	now := time.Now()

	reservations, err := w.checkoutRepo.GetAllRservations(ctx, now)
	if err != nil {
		log.Error().Err(err).Msg("error fetching checked-in reservations")
		return
	}

	for _, reservation := range reservations {
		hours, minutes, seconds := DataTypesTimeToHourMinuteSecond(reservation.TimeSlot)

		slotTime := time.Date(
			now.Year(), now.Month(), now.Day(),
			hours, minutes, seconds, 0, now.Location(),
		)

		endTime := slotTime.Add(reservationDuration)

		if now.After(endTime) {
			w.checkoutHandler(ctx, reservation, log)
		}
	}
}

func (w *CheckoutWorker) checkoutHandler(ctx context.Context, reservation models.Reservation, log zerolog.Logger) {
	// Checkout
	if err := w.checkoutRepo.Checkout(ctx, reservation); err != nil {
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

// deleteCache invalidates the cached table state after checkout.
func (w *CheckoutWorker) deleteCache(ctx context.Context, reservation models.Reservation) {
	w.redisStore.Delete(ctx,
		fmt.Sprintf("table:item:%d", reservation.TableID),
	)
}