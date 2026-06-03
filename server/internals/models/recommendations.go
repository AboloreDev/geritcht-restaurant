package models

import (
	"time"

	"gorm.io/gorm"
)

type Recommendation struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       *uint          `json:"user_id" gorm:"index"`
	Allergies    string         `json:"allergies"`
	Dietary      string         `json:"dietary"`
	Mood         string         `json:"mood"`
	Budget       string         `json:"budget"`
	SuggestedIDs string         `json:"suggested_ids"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	// Relationships
	User *User `json:"user,omitempty"`
}
