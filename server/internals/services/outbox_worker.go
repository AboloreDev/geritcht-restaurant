package services

import (
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
)

type OutboxWorker struct {
	outboxRepo     repositories.OutboxRepositoryInterface
	eventPublisher interfaces.Publisher
}

func NewOutboxWorker(
	outboxRepo repositories.OutboxRepositoryInterface,
	eventPublisher interfaces.Publisher) *OutboxWorker {
	return &OutboxWorker{
		outboxRepo:     outboxRepo,
		eventPublisher: eventPublisher,
	}
}

func (w *OutboxWorker) StartOutboxWorker(ctx context.Context, log zerolog.Logger) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.ProcessOutbox(ctx, log)
			log.Info().Msg("Outbox worker executed")
		case <-ctx.Done():
			log.Info().Msg("Shutting down outbox worker")
			return
		}
	}
}

func (w *OutboxWorker) ProcessOutbox(ctx context.Context, log zerolog.Logger) {
	events, err := w.outboxRepo.GetPendingEvents(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch outbox events")
		return
	}

	for _, event := range events {
		// Publish message
		err := w.publish(event)

		// If there is an error, increment the retry count by 1
		if err != nil {
			w.outboxRepo.UpdateRetryCount(ctx, &event)
			continue
		}

		// Mark event as published
		err = w.outboxRepo.MarkAsPublished(ctx, &event)
		if err != nil {
			log.Error().Err(err).Msg("failed to publish events")
			return
		}
	}
}

func (w *OutboxWorker) publish(event models.OutboxEvent) error {
	w.eventPublisher.PublishMessage(
		event.EventType,
		event.Payload,
		map[string]string{"Priority": "Important Mail"},
	)

	return nil
}
