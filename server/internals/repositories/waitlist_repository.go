package repositories

import (
	"context"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type WaitlistRepository struct {
	db *gorm.DB
}

func NewWaitlistRepository(db *gorm.DB) *WaitlistRepository {
	return &WaitlistRepository{db: db}
}

func (r *WaitlistRepository) CountAvailableTables(ctx context.Context, date, timeSlot string, partySize int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Table{}).
		Where("capacity >= ?", partySize).
		Where("status = ?", models.TableStatusAvailable).
		Where("id NOT IN (?)",
			r.db.Model(&models.Reservation{}).
				Select("table_id").
				Where("date = ? AND time_slot = ?", date, timeSlot).
				Where("status NOT IN ?", []string{"cancelled", "no_show"}),
		).
		Count(&count).Error
	return count, err
}

func (r *WaitlistRepository) GetByUserDateSlot(ctx context.Context, userID uint, date, timeSlot string) (*models.Waitlist, error) {
	var waitlist models.Waitlist
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND date = ? AND time_slot = ?", userID, date, timeSlot).
		First(&waitlist).Error
	if err != nil {
		return nil, err
	}
	return &waitlist, nil
}

func (r *WaitlistRepository) Create(ctx context.Context, waitlist *models.Waitlist) error {
	return r.db.WithContext(ctx).Create(waitlist).Error
}

func (r *WaitlistRepository) GetPosition(ctx context.Context, date string, timeSlot datatypes.Time, createdAt time.Time) (int64, error) {
	var position int64
	err := r.db.WithContext(ctx).Model(&models.Waitlist{}).
		Where("date = ? AND time_slot = ? AND status = ? AND created_at < ?",
			date, timeSlot, models.WaitlistStatusWaiting, createdAt).
		Count(&position).Error
	return position, err
}
