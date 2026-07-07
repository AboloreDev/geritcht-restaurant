package repositories

import (
	"context"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReservationNoShowRepository struct {
	db *gorm.DB
}

func NewReservationNoShowRepository(db *gorm.DB) *ReservationNoShowRepository {
	return &ReservationNoShowRepository{
		db: db,
	}
}

func (r *ReservationNoShowRepository) GetAllReservations(ctx context.Context) ([]models.Reservation, error) {
	var reservations []models.Reservation
	threshold := time.Now().Add(-5*time.Minute).Format("15:04:00")
	now :=	time.Now().Format("2006-01-02")

	err := r.db.Preload("Table").
		WithContext(ctx).
		Where("date = ? AND time_slot <= ? AND status = ?",
		now,
		threshold,
			models.ReservationStatusConfirmed).
		Find(&reservations).Error
	if err != nil {
		return nil, domain.ErrReservationNotFound
	}

	return reservations, nil
}

func (r *ReservationNoShowRepository) UpdateReservationStatus(ctx context.Context, tx *gorm.DB, reservation *models.Reservation, status models.ReservationStatus) error {
	err := tx.Model(reservation).
		WithContext(ctx).
		Update("status", status).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ReservationNoShowRepository) UpdateReservedTableStatus(ctx context.Context, tx *gorm.DB, reservation *models.Reservation, status models.TableStatus) error {
	var table models.Table

	err := tx.Model(&table).
		WithContext(ctx).
		Where("id = ?", reservation.TableID).
		Update("status", status).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ReservationNoShowRepository) GetUsersWaitlist(ctx context.Context, tx *gorm.DB, status models.WaitlistStatus, date interface{}, timeSlot datatypes.Time, partySize int) (*models.Waitlist, error) {
	var waitlist models.Waitlist

	err := tx.Preload("User").
		WithContext(ctx).
		Where("date = ? AND time_slot = ? AND party_size = ? AND status = ?",
			date, timeSlot, partySize, status).
		Order("created_at ASC").
		First(&waitlist).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &waitlist, nil
}
func (r *ReservationNoShowRepository) UpdateWaitlistStatus(ctx context.Context, tx *gorm.DB, waitlist *models.Waitlist, status models.WaitlistStatus, notifiedAt, expiresAt time.Time) error {
	tx.Model(waitlist).
		WithContext(ctx).
		Updates(map[string]interface{}{
			"status":      status,
			"notified_at": notifiedAt,
			"expires_at":  expiresAt,
		})

	return nil
}

func (r *ReservationNoShowRepository) MarkReservationNoShow(ctx context.Context, reservation *models.Reservation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.UpdateReservationStatus(ctx, tx, reservation, models.ReservationStatusNoShow); err != nil {
			return err
		}

		if err := r.UpdateReservedTableStatus(
			ctx,
			tx,
			reservation,
			models.TableStatusAvailable,
		); err != nil {
			return err
		}

		waitlist, err := r.GetUsersWaitlist(
			ctx,
			tx,
			models.WaitlistStatusWaiting,
			reservation.Date,
			reservation.TimeSlot,
			reservation.PartySize,
		)
		if err != nil {
			return err
		}

		now := time.Now()

		return r.UpdateWaitlistStatus(
			ctx,
			tx,
			waitlist,
			models.WaitlistStatusNotified,
			now,
			now.Add(10*time.Minute),
		)
	})
}
