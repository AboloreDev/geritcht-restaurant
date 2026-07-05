package repositories

import (
	"context"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type OutboxRepository struct {
	db *gorm.DB
}

func NewOutboxRepository(db *gorm.DB) *OutboxRepository {
	return &OutboxRepository{
		db: db,
	}
}

func (o *OutboxRepository) GetPendingEvents(ctx context.Context) ([]models.OutboxEvent, error) {
	var events []models.OutboxEvent

	err := o.db.WithContext(ctx).Where("status = ? AND retry_count < ?", "pending", 5).
		Order("created_at ASC").
		Limit(100).
		Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (o *OutboxRepository) UpdateRetryCount(ctx context.Context, event *models.OutboxEvent) error {
	var err error

	o.db.Model(&event).WithContext(ctx).Update("retry_count", event.RetryCount+1)
	log.Error().Err(err).
		Uint("event_id", event.ID).
		Msg("outbox publish failed — will retry")

	return nil
}

func (o *OutboxRepository) MarkAsPublished(ctx context.Context, event *models.OutboxEvent) error {

	now := time.Now()

	err := o.db.Model(&event).Updates(map[string]interface{}{
		"status":       "published",
		"processed_at": &now,
	}).Error
	if err != nil {
		return err
	}

	return nil
}
