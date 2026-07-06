package subscriber

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/AboloreDev/geritcht-restaurant/internals/email"
	"github.com/AboloreDev/geritcht-restaurant/internals/events"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type EventSubscriber struct {
	streamsSubscriber message.Subscriber
	queueName         string
}

func NewEventSubscriber(ctx context.Context, cfg *config.RedisConfig, client *redis.Client) (*EventSubscriber, error) {
	logger := watermill.NewStdLogger(true, true)

	subscriberConfig := redisstream.SubscriberConfig{
		Client:        client,
		ConsumerGroup: "email-service",
		Consumer:      "email-service-1",
		BlockTime:     time.Second * 10,
	}

	streamsSubscriber, err := redisstream.NewSubscriber(subscriberConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("Failed to create publisher %v", err)
	}

	return &EventSubscriber{
		streamsSubscriber: streamsSubscriber,
		queueName:         cfg.QUEUE_NAME,
	}, nil
}

func (s *EventSubscriber) Start(ctx context.Context, emailClient *email.ResendEmailClient) error {
	messages, err := s.streamsSubscriber.Subscribe(ctx, s.queueName)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	go func() {
		for {
			select {
			case msg := <-messages:
				err := s.handleMessage(msg, emailClient)
				if err != nil {
					log.Printf("Error processing message: %s", err)
					msg.Nack()
				} else {
					msg.Ack()
				}
			case <-ctx.Done():
				s.streamsSubscriber.Close()
				return
			}
		}
	}()

	return nil
}

func (s *EventSubscriber) handleMessage(msg *message.Message, emailClient *email.ResendEmailClient) error {
	eventType := msg.Metadata.Get("event_type")

	switch eventType {
	case events.ChannelEmailVerification:
		return s.HandleSendVerificationMail(msg, emailClient)
	case events.ChannelEmailPasswordReset:
		return s.HandleSendPasswordReset(msg, emailClient)
	case events.ChannelEmailPasswordChanged:
		return s.HandleSendPasswordChangedMail(msg, emailClient)
	case events.ChannelEmailReservationConfirm:
		return s.HandleReservationConfirmationMail(msg, emailClient)
	case events.ChannelEmailReservationCancelled:
		return s.HandleReservationCancellationMail(msg, emailClient)
	case events.ChannelEmailReservationReminder:
		return s.HandleReservationReminderMail(msg, emailClient)
	case events.ChannelEmailReservationCheckedIn:
		return s.HandleReservationCheckInMail(msg, emailClient)
	case events.ChannelEmailReservationNoShow:
		return s.HandleReservationNoShowMail(msg, emailClient)
	case events.ChannelOrderConfirmation:
		return s.HandleOrderConfirmationMail(msg, emailClient)
	case events.ChannelOrderRefunded:
		return s.HandleOrderRefundPayload(msg, emailClient)
	case events.ChannelEmailLowStockAlert:
		return s.HandleLowStockAlert(msg, emailClient)
	case events.ChannelEmailWaitlistNotification:
		return s.HandleWaitlistNotifier(msg, emailClient)
	default:
		log.Printf("Unknown event type %s", eventType)
		return nil
	}
}
