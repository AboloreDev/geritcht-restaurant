package dto

type CreateTableRequest struct {
	Name     string `json:"name" binding:"required"`
	Capacity int    `json:"capacity" binding:"required,gt=0"`
	Location string `json:"location"`
}

type UpdateTableRequest struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	Location string `json:"location"`
	Status   string `json:"status"`
}

type UpdateTableStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=available occupied reserved"`
}

type TableDetailResponse struct {
	ID                 uint                 `json:"id"`
	Name               string               `json:"name"`
	Capacity           int                  `json:"capacity"`
	Location           string               `json:"location"`
	Status             string               `json:"status"`
	QRCode             string               `json:"qr_code,omitempty"`
	CurrentReservation *ReservationResponse `json:"current_reservation,omitempty"`
	CurrentOrder       *OrderResponse       `json:"current_order,omitempty"`
}
