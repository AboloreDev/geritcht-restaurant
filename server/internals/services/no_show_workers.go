package services

import (
	"context"
	"fmt"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/events"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type NoShowWorker struct {
	db         *gorm.DB
	publisher  interfaces.Publisher
	redisStore interfaces.Cacher
}

func NewNoShowWorker(
	db *gorm.DB,
	publisher interfaces.Publisher,
	redisStore interfaces.Cacher) *NoShowWorker {
	return &NoShowWorker{
		db:         db,
		publisher:  publisher,
		redisStore: redisStore,
	}
}

func (w *NoShowWorker) StartMarkNoShowWorker(ctx context.Context, log zerolog.Logger) {
	ticker := time.NewTicker(5 * time.Minute)
	heartbeat := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	defer heartbeat.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Shutting down MarkNoShow worker")
			return
		case <-ticker.C:
			w.processNoShow()
			log.Info().Msg("MarkNoShow worker executed")
		case <-heartbeat.C:
			log.Info().Msg("MarkNoShow worker is alive")
		}
	}
}

func (w *NoShowWorker) processNoShow() {
	var reservations []models.Reservation

	w.db.Preload("Table").Where("date = ? AND time_slot = ? AND status = ?",
		time.Now().Format("2006-01-02"),
		time.Now().Add(-45*time.Minute).Format("15:04"),
		models.ReservationStatusConfirmed,
	).Find(&reservations)

	for _, reservation := range reservations {
		w.markNoShow(reservation)
	}
}

func (w *NoShowWorker) markNoShow(reservation models.Reservation) {
	var waitlist models.Waitlist
	var table models.Table

	w.db.Transaction(func(tx *gorm.DB) error {
		tx.Model(&reservation).Update("status", models.ReservationStatusNoShow)

		tx.Model(&table).Where("id = ?", reservation.TableID).
			Update("status", models.TableStatusAvailable)

		err := tx.Preload("User").
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

	w.publisher.PublishMessage(
		events.ChannelEmailReservationNoShow,
		events.ReservationNoShowPayload{
			FirstName: reservation.User.FirstName,
			Email:     reservation.User.Email,
			Date:      reservation.Date.Format("2006-01-02"),
			TimeSlot:  utils.FormatDataTypesTime(reservation.TimeSlot),
			TableName: reservation.Table.Name,
		},
		map[string]string{"Priority": "Important Mail"},
	)

	w.redisStore.FlushByPattern(ctx,
		fmt.Sprintf("availability:%s:%s:*", reservation.Date, reservation.TimeSlot),
	)
}
