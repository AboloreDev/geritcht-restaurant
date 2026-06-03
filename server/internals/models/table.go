package models

import (
	"time"

	"gorm.io/gorm"
)

type Table struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Capacity  int            `json:"capacity" gorm:"not null"`
	Location  string         `json:"location"`
	Status    TableStatus    `json:"status" gorm:"default:available"`
	QRCode    string         `json:"qr_code"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	// Relationships
	Reservations []Reservation `json:"-"`
	Orders       []Order       `json:"-"`
}

type TableStatus string

const (
	TableStatusAvailable TableStatus = "available"
	TableStatusOccupied  TableStatus = "occupied"
	TableStatusReserved  TableStatus = "reserved"
)
