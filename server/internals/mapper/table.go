package mapper

import (
	"github.com/AboloreDev/geritcht-restaurant/internals/dto"
	"github.com/AboloreDev/geritcht-restaurant/internals/models"
)

func TableResponse(table *models.Table) *dto.TableResponse {
	return &dto.TableResponse{
		ID:       table.ID,
		Name:     table.Name,
		Capacity: table.Capacity,
		Location: table.Location,
		Status:   string(table.Status),
	}
}

func TableDetailResponse(table *models.Table) *dto.TableDetailResponse {
	response := &dto.TableDetailResponse{
		ID:       table.ID,
		Name:     table.Name,
		Capacity: table.Capacity,
		Location: table.Location,
		Status:   string(table.Status),
	}

	if r := findActiveReservation(table.Reservations); r != nil {
		response.CurrentReservation = ReservationResponse(r)
	}

	if o := findActiveOrder(table.Orders); o != nil {
		response.CurrentOrder = OrderResponse(o)
	}

	return response
}

func findActiveReservation(reservations []models.Reservation) *models.Reservation {
	for i := range reservations {
		s := reservations[i].Status
		if s == models.ReservationStatusConfirmed ||
			s == models.ReservationStatusPending {
			return &reservations[i]
		}
	}
	return nil
}

func findActiveOrder(orders []models.Order) *models.Order {
	for i := range orders {
		s := orders[i].Status
		if s == models.OrderStatusPending ||
			s == models.OrderStatusConfirmed ||
			s == models.OrderStatusPreparing {
			return &orders[i]
		}
	}
	return nil
}
