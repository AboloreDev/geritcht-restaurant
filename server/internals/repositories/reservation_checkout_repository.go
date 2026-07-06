package repositories

import (
	"context"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"gorm.io/gorm"
)

type ReservationCheckoutRepository struct {
	db *gorm.DB
}

func NewReservationCheckoutRepository(db *gorm.DB) *ReservationCheckoutRepository {
	return &ReservationCheckoutRepository{
		db: db,
	}
}

func (r *ReservationCheckoutRepository) GetAllRservations(ctx context.Context, now time.Time) ([]models.Reservation, error) {
	var reservations []models.Reservation
	today := now.Format("2006-01-02")

	err := r.db.Preload("Table").
		WithContext(ctx).
		Where("date = ? AND status = ?", today, models.ReservationStatusCheckedIn).
		Find(&reservations).Error

	return reservations, err
}

func (r *ReservationCheckoutRepository) UpdateReservationStatus(ctx context.Context, tx *gorm.DB, reservation *models.Reservation, status models.ReservationStatus) error {

	if err := tx.Model(&reservation).WithContext(ctx).
		Update("status", models.ReservationStatusCompleted).Error; err != nil {
		return err
	}

	return nil
}

func (r *ReservationCheckoutRepository) UpdateReservedTableStatus(ctx context.Context, tx *gorm.DB, reservation *models.Reservation, status models.TableStatus) error {

	if err := tx.Model(&models.Table{}).WithContext(ctx).
		Where("id = ?", reservation.TableID).
		Update("status", models.TableStatusAvailable).Error; err != nil {
		return err
	}
	return nil
}

func (r *ReservationCheckoutRepository) Checkout(ctx context.Context, reservation models.Reservation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.UpdateReservationStatus(ctx, tx, &reservation, models.ReservationStatusCompleted); err != nil {
			return err
		}

		if err := r.UpdateReservedTableStatus(ctx, tx, &reservation, models.TableStatusAvailable); err != nil {
			return err
		}

		return nil
	})
}
