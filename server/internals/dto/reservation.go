package dto

import "time"

type CreateReservationRequest struct {
	Date            string `json:"date" binding:"required"`
	TimeSlot        string `json:"time_slot" binding:"required"`
	PartySize       int    `json:"party_size" binding:"required,gt=0"`
	SpecialRequests string `json:"special_requests"`
}

type UpdateReservationRequest struct {
	Date            string `json:"date"`
	TimeSlot        string `json:"time_slot"`
	PartySize       int    `json:"party_size"`
	SpecialRequests string `json:"special_requests"`
}

type CheckAvailabilityRequest struct {
	Date      string `form:"date" binding:"required"`
	TimeSlot  string `form:"time_slot" binding:"required"`
	PartySize int    `form:"party_size" binding:"required,gt=0"`
}

type JoinWaitlistRequest struct {
	Date      string `json:"date" binding:"required"`
	TimeSlot  string `json:"time_slot" binding:"required"`
	PartySize int    `json:"party_size" binding:"required,gt=0"`
}

type ReservationFilterRequest struct {
	Date     string `form:"date"`
	Status   string `form:"status"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

type TableResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	Location string `json:"location"`
	Status   string `json:"status"`
	QRCode   string `json:"qr_code,omitempty"`
}

type ReservationResponse struct {
	ID              uint          `json:"id"`
	UserID          uint          `json:"user_id"`
	User            *UserResponse `json:"user,omitempty"`
	TableID         uint          `json:"table_id"`
	Table           TableResponse `json:"table"`
	Date            string        `json:"date"`
	TimeSlot        string        `json:"time_slot"`
	PartySize       int           `json:"party_size"`
	Status          string        `json:"status"`
	SpecialRequests string        `json:"special_requests"`
	CheckedInAt     *time.Time    `json:"checked_in_at"`
	CreatedAt       time.Time     `json:"created_at"`
}

type AvailabilityResponse struct {
	Date      string             `json:"date"`
	TimeSlots []TimeSlotResponse `json:"time_slots"`
}

type TimeSlotResponse struct {
	TimeSlot        string          `json:"time_slot"`
	IsAvailable     bool            `json:"is_available"`
	AvailableTables []TableResponse `json:"available_tables"`
}

type WaitlistResponse struct {
	ID         uint       `json:"id"`
	UserID     uint       `json:"user_id"`
	Date       string     `json:"date"`
	TimeSlot   string     `json:"time_slot"`
	PartySize  int        `json:"party_size"`
	Status     string     `json:"status"`
	Position   int        `json:"position"`
	NotifiedAt *time.Time `json:"notified_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

type ReservationListResponse struct {
	Reservations []ReservationResponse `json:"reservations"`
	Total        int64                 `json:"total"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"page_size"`
}
