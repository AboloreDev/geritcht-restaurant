package models

import (
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
	ID              uint              `json:"id" gorm:"primaryKey"`
	UserID          uint              `json:"user_id" gorm:"not null;index"`
	TableID         uint              `json:"table_id" gorm:"not null;index"`
	Date            time.Time         `json:"date" gorm:"not null"`
	TimeSlot        time.Time         `json:"time_slot" gorm:"not null"`
	PartySize       int               `json:"party_size" gorm:"not null"`
	Status          ReservationStatus `json:"status" gorm:"default:pending"`
	SpecialRequests string            `json:"special_requests"`
	CheckedInAt     *time.Time        `json:"checked_in_at"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `json:"-" gorm:"index"`
	// Relationships
	User   User    `json:"user,omitempty"`
	Table  Table   `json:"table,omitempty"`
	Orders []Order `json:"-"`
}

type Waitlist struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null;index"`
	Date       time.Time      `json:"date" gorm:"not null"`
	TimeSlot   time.Time      `json:"time_slot" gorm:"not null"`
	PartySize  int            `json:"party_size" gorm:"not null"`
	Status     WaitlistStatus `json:"status" gorm:"default:waiting"`
	NotifiedAt *time.Time     `json:"notified_at"`
	ExpiresAt  *time.Time     `json:"expires_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	// Relationships
	User User `json:"user,omitempty"`
}

type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "pending"
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusCheckedIn ReservationStatus = "checked_in"
	ReservationStatusNoShow    ReservationStatus = "no_show"
	ReservationStatusCancelled ReservationStatus = "cancelled"
)

type WaitlistStatus string

const (
	WaitlistStatusWaiting   WaitlistStatus = "waiting"
	WaitlistStatusNotified  WaitlistStatus = "notified"
	WaitlistStatusConfirmed WaitlistStatus = "confirmed"
	WaitlistStatusExpired   WaitlistStatus = "expired"
)
