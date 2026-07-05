package repositories

import (
	"context"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReservationReminderRepository struct {
	db *gorm.DB
}

func NewReservationReminderRepository(db *gorm.DB) *ReservationReminderRepository {
	return &ReservationReminderRepository{db: db}
}

func (r *ReservationReminderRepository) GetAllUpcomingReservations(ctx context.Context, now time.Time, targetSlot datatypes.Time) ([]models.Reservation, error) {
	var reservations []models.Reservation
	r.db.Preload("Table").Preload("User").WithContext(ctx).
		Where("date = ? AND status = ? AND time_slot = ? AND reminder_sent = ?",
			now.Format("2006-01-02"),
			models.ReservationStatusConfirmed,
			targetSlot,
			false,
		).Find(&reservations)

	return reservations, nil
}

func (r *ReservationReminderRepository) UpdateReminderValue(ctx context.Context, reservation *models.Reservation) error {
	r.db.Model(&reservation).WithContext(ctx).
		Update("reminder_sent", true)
	return nil
}
