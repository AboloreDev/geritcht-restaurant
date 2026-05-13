package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type EventPublisher struct {
	streamsPublisher message.Publisher
	queueName        string
}

func NewEventPublisher(ctx context.Context, cfg *config.RedisConfig, client *redis.Client) (*EventPublisher, error) {
	logger := watermill.NewStdLogger(true, true)

	publisherConfig := redisstream.PublisherConfig{
		Client:     client,
		Marshaller: nil,
	}

	streamsPublisher, err := redisstream.NewPublisher(publisherConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("Failed to create publisher %v", err)
	}

	return &EventPublisher{
		streamsPublisher: streamsPublisher,
		queueName:        cfg.QUEUE_NAME,
	}, nil
}

func (p *EventPublisher) PublishMessage(eventType string, payload interface{}, metadata map[string]string) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), data)

	msg.Metadata.Set("event_type", eventType)
	for i, j := range metadata {
		msg.Metadata.Set(i, j)
	}

	return p.streamsPublisher.Publish(p.queueName, msg)
}

func (p *EventPublisher) CloseMessage() error {
	return p.streamsPublisher.Close()
}
