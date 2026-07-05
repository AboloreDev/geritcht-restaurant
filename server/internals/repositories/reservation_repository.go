package repositories

import (
	"context"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/domain"
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

// Table

func (r *ReservationRepository) GetTableByIDAndCapacity(ctx context.Context, tableID uint, partySize int) (*models.Table, error) {
	var table models.Table
	err := r.db.WithContext(ctx).
		Where("id = ? AND capacity >= ?", tableID, partySize).
		First(&table).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &table, nil
}

func (r *ReservationRepository) GetTablesByCapacity(ctx context.Context, partySize int) ([]models.Table, error) {
	var tables []models.Table
	err := r.db.WithContext(ctx).Where("capacity >= ?", partySize).Find(&tables).Error
	return tables, err
}

func (r *ReservationRepository) UpdateTableStatus(ctx context.Context, tx *gorm.DB, tableID uint, status models.TableStatus) error {
	return r.getDB(tx).WithContext(ctx).Model(&models.Table{}).
		Where("id = ?", tableID).Update("status", status).Error
}

// Reservation

func (r *ReservationRepository) GetReservationsByDateAndSlot(ctx context.Context, date string, timeSlot datatypes.Time) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.WithContext(ctx).Select("table_id").
		Where("date = ? AND time_slot = ? AND status NOT IN ?",
			date, timeSlot, []string{"cancelled", "no_show"}).
		Find(&reservations).Error
	return reservations, err
}

func (r *ReservationRepository) CountByTableDateSlot(ctx context.Context, tableID uint, date string, timeSlot datatypes.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Reservation{}).
		Where("table_id = ? AND date = ? AND time_slot = ? AND status NOT IN ?",
			tableID, date, timeSlot, []string{"cancelled", "no_show"}).
		Count(&count).Error
	return count, err
}

func (r *ReservationRepository) Create(ctx context.Context, tx *gorm.DB, reservation *models.Reservation) error {
	return r.getDB(tx).WithContext(ctx).Create(reservation).Error
}

func (r *ReservationRepository) GetByIDAndUser(ctx context.Context, reservationID, userID uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", reservationID, userID).
		First(&reservation).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &reservation, nil
}

func (r *ReservationRepository) GetByIDWithRelations(ctx context.Context, reservationID uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.WithContext(ctx).Preload("User").Preload("Table").
		First(&reservation, reservationID).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &reservation, nil
}

func (r *ReservationRepository) GetByIDAndStatus(ctx context.Context, reservationID uint, status models.ReservationStatus) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.WithContext(ctx).
		Where("id = ? AND status = ?", reservationID, status).
		First(&reservation).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &reservation, nil
}

func (r *ReservationRepository) UpdateStatus(ctx context.Context, tx *gorm.DB, reservationID uint, updates map[string]interface{}) error {
	return r.getDB(tx).WithContext(ctx).Model(&models.Reservation{}).
		Where("id = ?", reservationID).Updates(updates).Error
}

func (r *ReservationRepository) GetAllByUser(ctx context.Context, userID uint, req *dto.ReservationFilterRequest) ([]models.Reservation, int64, error) {
	var reservations []models.Reservation
	var count int64
	offset := utils.Pagination(req.Page, req.PageSize)

	query := r.db.WithContext(ctx).Preload("User").Preload("Table").
		Where("user_id = ?", userID)

	if req.Date != "" {
		query = query.Where("date = ?", req.Date)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	query.Model(&models.Reservation{}).Count(&count)

	err := query.Order("date ASC, time_slot ASC").
		Offset(offset).Limit(req.PageSize).Find(&reservations).Error
	if err != nil {
		return nil, 0, err
	}

	return reservations, count, nil
}

func (r *ReservationRepository) GetAll(ctx context.Context, req *dto.ReservationFilterRequest) ([]models.Reservation, int64, error) {
	var reservations []models.Reservation
	var count int64
	offset := utils.Pagination(req.Page, req.PageSize)

	query := r.db.WithContext(ctx).Preload("User").Preload("Table")

	if req.Date != "" {
		query = query.Where("date = ?", req.Date)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	query.Model(&models.Reservation{}).Count(&count)

	err := query.Order("time_slot ASC").
		Offset(offset).Limit(req.PageSize).Find(&reservations).Error
	if err != nil {
		return nil, 0, err
	}

	return reservations, count, nil
}

func (r *ReservationRepository) GetTodayReservations(ctx context.Context, req *dto.ReservationFilterRequest) ([]models.Reservation, int64, error) {
	var reservations []models.Reservation
	var count int64
	offset := utils.Pagination(req.Page, req.PageSize)

	query := r.db.WithContext(ctx).Preload("User").Preload("Table").
		Where("date = ?", time.Now().Format("2006-01-02"))

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	query.Model(&models.Reservation{}).Count(&count)

	err := query.Order("time_slot ASC").
		Offset(offset).Limit(req.PageSize).Find(&reservations).Error
	if err != nil {
		return nil, 0, err
	}

	return reservations, count, nil
}

// Waitlist

func (r *ReservationRepository) GetFirstWaitlistByDateSlot(ctx context.Context, tx *gorm.DB, date interface{}, timeSlot datatypes.Time, partySize int) (*models.Waitlist, error) {
	var waitlist models.Waitlist
	err := r.getDB(tx).WithContext(ctx).Preload("User").
		Where("date = ? AND time_slot = ? AND party_size = ? AND status = ?",
			date, timeSlot, partySize, models.WaitlistStatusWaiting).
		Order("created_at ASC").
		First(&waitlist).Error
	if err != nil {
		return nil, err
	}
	return &waitlist, nil
}

func (r *ReservationRepository) UpdateWaitlistStatus(ctx context.Context, tx *gorm.DB, waitlist *models.Waitlist, updates map[string]interface{}) error {
	return r.getDB(tx).WithContext(ctx).Model(waitlist).Updates(updates).Error
}

func (r *ReservationRepository) LockTableForUpdate(ctx context.Context, tx *gorm.DB, tableID uint) (*models.Table, error) {
	var table models.Table
	err := r.getDB(tx).WithContext(ctx).
		Raw("SELECT * FROM tables WHERE id = ? FOR UPDATE", tableID).
		Scan(&table).Error
	return &table, err
}
