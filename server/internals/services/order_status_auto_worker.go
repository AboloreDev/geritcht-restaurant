package services

import (
	"context"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	websockets "github.com/AboloreDev/geritcht-restaurant/internals/web-sockets"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

var nextStatus = map[models.OrderStatus]models.OrderStatus{
    models.OrderStatusConfirmed: models.OrderStatusPreparing,
    models.OrderStatusPreparing: models.OrderStatusReady,
    models.OrderStatusReady:     models.OrderStatusCompleted,
}

type OrderAutoWorker struct {
	db *gorm.DB
	hub *websockets.Hub
}

func NewOrderAutoWorker(db *gorm.DB, hub *websockets.Hub) *OrderAutoWorker {
	return &OrderAutoWorker{
		db: db,
		hub: hub,
	}
}

func (s *OrderAutoWorker) StartOrderUpdateWorker(ctx context.Context, log zerolog.Logger) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Stopping order update worker...")
			return
		case <-ticker.C:
			s.ProcessOrderStatus(log)
		}
	}
}

func (s *OrderAutoWorker) ProcessOrderStatus(log zerolog.Logger) {
    var orders []models.Order

    s.db.Where("status IN ? AND type = ? AND updated_at < ?",
        []string{"confirmed", "preparing", "ready"},
        models.OrderTypeTakeout,
        time.Now().Add(-5*time.Minute),
    ).Find(&orders)

    for _, order := range orders {
        next, exists := nextStatus[order.Status]
        if !exists {
            continue
        }

        if err := s.db.Model(&order).Update("status", next).Error; err != nil {
            log.Error().Err(err).
                Uint("order_id", order.ID).
                Msg("failed to update order status")
            continue
        }
        msg := websockets.BuildMessageWithStatus(order.ID, string(next))
        s.hub.Broadcast(order.ID, msg)

        log.Info().
            Uint("order_id", order.ID).
            Str("status", string(next))
    }
}

