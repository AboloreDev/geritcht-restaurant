package models

import (
	"time"
)

type OutboxEvent struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	EventType   string     `json:"event_type" gorm:"not null"`
	Payload     string     `json:"payload" gorm:"not null"`
	Status      string     `json:"status" gorm:"default:pending"`
	RetryCount  int        `json:"retry_count" gorm:"default:0"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
