package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Email         string         `json:"email" gorm:"uniqueIndex;not null"`
	Password      string         `json:"-" gorm:"not null"`
	FirstName     string         `json:"first_name" gorm:"not null"`
	LastName      string         `json:"last_name" gorm:"not null"`
	PhoneNumber   string         `json:"phone_number" gorm:"uniqueIndex;not null"`
	Role          UserRole       `json:"role" gorm:"default:customer"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	EmailVerified bool           `json:"email_verified" gorm:"default:false"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	// Relationships
	RefreshTokens  []RefreshToken  `json:"-"`
	Cart           []Cart          `json:"-"`
	Orders         []Order         `json:"-"`
	Reservations   []Reservation   `json:"-"`
	StockMovements []StockMovement `json:"-" gorm:"foreignKey:CreatedBy"`
}

type RefreshToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	TokenHash string         `json:"-" gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	// Relationships
	User User `json:"-"`
}

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleStaff    UserRole = "staff"
	RoleAdmin    UserRole = "admin"
)
