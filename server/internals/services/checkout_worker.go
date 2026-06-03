package services

import (
	"context"
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/rs/zerolog"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CheckoutWorker struct {
	db         *gorm.DB
	redisStore interfaces.Cacher
}

func NewCheckoutWorker(db *gorm.DB, redisStore interfaces.Cacher) *CheckoutWorker {
	return &CheckoutWorker{
		db:         db,
		redisStore: redisStore,
	}
}

var reservationDuration = 70 * time.Minute

func (w *CheckoutWorker) Start(appCtx context.Context, log zerolog.Logger) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.processCheckOut(log)
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

func (w *CheckoutWorker) processCheckOut(log zerolog.Logger) {
	now := time.Now()
	var reservations []models.Reservation
	today := now.Format("2006-01-02")

	err := w.db.Preload("Table").
		Where("date = ? AND status = ?", today, models.ReservationStatusCheckedIn).
		Find(&reservations).Error

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
			w.checkout(reservation, log)
		}
	}
}

func (w *CheckoutWorker) checkout(reservation models.Reservation, log zerolog.Logger) {
	err := w.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&reservation).
			Update("status", models.ReservationStatusCompleted).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Table{}).
			Where("id = ?", reservation.TableID).
			Update("status", models.TableStatusAvailable).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).
			Uint("reservation_id", reservation.ID).
			Msg("failed to checkout reservation")
		return
	}

	w.redisStore.Delete(ctx,
		fmt.Sprintf("table:item:%d", reservation.TableID),
	)

	log.Info().
		Uint("reservation_id", reservation.ID).
		Uint("table_id", reservation.TableID).
		Msg("reservation checked out automatically")
}
