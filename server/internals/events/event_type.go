package events

const (
	ChannelEmailVerification         = "email:verification"
	ChannelEmailPasswordReset        = "email:password_reset"
	ChannelEmailPasswordChanged      = "email:password_changed"
	ChannelEmailOrderReceipt         = "email:order_receipt"
	ChannelEmailReservationConfirm   = "email:reservation_confirm"
	ChannelEmailReservationReminder  = "email:reservation_reminder"
	ChannelEmailLowStockAlert        = "email:low_stock_alert"
	ChannelEmailReservationCheckedIn = "email:checked_in"
	ChannelEmailWaitlistNotification = "email:waitlist_notification"
	ChannelEmailReservationCancelled = "email:cancelled"
	ChannelEmailReservationNoShow    = "email:noshow"
)

type VerificationEmailPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	Token     string `json:"token"`
}

type PasswordResetEmailPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	Token     string `json:"token"`
}

type PasswordChangedEmailPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
}

type OrderReceiptEmailPayload struct {
	Email       string             `json:"email"`
	FirstName   string             `json:"first_name"`
	OrderID     uint               `json:"order_id"`
	Items       []OrderItemPayload `json:"items"`
	TotalAmount float64            `json:"total_amount"`
	Reference   string             `json:"reference"`
}

type OrderItemPayload struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type ReservationConfirmPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	Date      string `json:"date"`
	TimeSlot  string `json:"time_slot"`
	PartySize int    `json:"party_size"`
	TableName string `json:"table_name"`
}
type ReservationCheckInPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	Date      string `json:"date"`
	TimeSlot  string `json:"time_slot"`
	PartySize int    `json:"party_size"`
	TableName string `json:"table_name"`
}

type LowStockAlertPayload struct {
	AdminEmail string            `json:"admin_email"`
	AdminName  string            `json:"admin_name"`
	Items      []LowStockPayload `json:"items"`
}

type LowStockPayload struct {
	Name         string  `json:"name"`
	Unit         string  `json:"unit"`
	CurrentStock float64 `json:"current_stock"`
	MinThreshold float64 `json:"min_threshold"`
}

type WaitlistNotificationPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	Date      string `json:"date"`
	TimeSlot  string `json:"time_slot"`
	PartySize int    `json:"party_size"`
	TableName string `json:"table_name"`
}

type ReservationCancelledPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	Date      string `json:"date"`
	TimeSlot  string `json:"time_slot"`
	PartySize int    `json:"party_size"`
	TableName string `json:"table_name"`
}

type ReservationReminderPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	Date      string `json:"date"`
	TimeSlot  string `json:"time_slot"`
	PartySize int    `json:"party_size"`
	TableName string `json:"table_name"`
}
type ReservationNoShowPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	Date      string `json:"date"`
	TimeSlot  string `json:"time_slot"`
	PartySize int    `json:"party_size"`
	TableName string `json:"table_name"`
}
