package services

import (
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type OutboxWorker struct {
	db             *gorm.DB
	eventPublisher interfaces.Publisher
}

func NewOutboxWorker(db *gorm.DB, eventPublisher interfaces.Publisher) *OutboxWorker {
	return &OutboxWorker{
		db:             db,
		eventPublisher: eventPublisher,
	}
}

func (w *OutboxWorker) StartOutboxWorker(ctx context.Context, log zerolog.Logger) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.ProcessOutbox(log)
			log.Info().Msg("Outbox worker executed")
		case <-ctx.Done():
			log.Info().Msg("Shutting down outbox worker")
			return
		}
	}
}

func (w *OutboxWorker) ProcessOutbox(log zerolog.Logger) {
	var events []models.OutboxEvent

	w.db.Where("status = ? AND retry_count < ?", "pending", 5).
		Order("created_at ASC").
		Limit(100).
		Find(&events)

	for _, event := range events {
		// Publish message
		err := w.eventPublisher.PublishMessage(
			event.EventType,
			event.Payload,
			map[string]string{"Priority": "Important Mail"},
		)
		// If there is an error, increment the retry count by 1
		if err != nil {
			w.db.Model(&event).Update("retry_count", event.RetryCount+1)
			log.Error().Err(err).
				Uint("event_id", event.ID).
				Msg("outbox publish failed — will retry")
			continue
		}

		now := time.Now()
		w.db.Model(&event).Updates(map[string]interface{}{
			"status":       "published",
			"processed_at": &now,
		})
	}
}
