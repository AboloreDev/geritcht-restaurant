package mapper

import (
	"testing"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"gorm.io/datatypes"
)

func Test_CartResponse(t *testing.T) {
	cart := &models.Cart{
		ID:     1,
		UserID: 10,
		CartItems: []models.CartItem{
			{
				ID:       1,
				MenuID:   100,
				Quantity: 2,
				Menu: models.Menu{
					ID:    100,
					Name:  "Jollof Rice",
					Price: 2500,
				},
			},
		},
	}

	response := ConvertToCartResponse(cart)

	if response.ID != cart.ID {
		t.Errorf("Expected ID %d, got %d", cart.ID, response.ID)
	}

	if response.UserID != cart.UserID {
		t.Errorf("Expected UserID %d, got %d", cart.UserID, response.UserID)
	}

	if response.Total != 5000 {
		t.Errorf("Expected Total %v, got %v", 5000, response.Total)
	}

	if response.CartItems[0].MenuItem.Name != "Jollof Rice" {
		t.Errorf("Expected MenuItem Name %s, got %s", "", response.CartItems[0].MenuItem.Name)
	}

	if response.ItemCount != 1 {
		t.Errorf("Expected ItemCount %d, got %d", 1, response.ItemCount)
	}
}

func Test_OrderResponse(t *testing.T) {
	userID := uint(10)

	order := &models.Order{
		ID:     1,
		UserID: &userID,
		User: &models.User{
			ID: userID,
		},
		OrderItems: []models.OrderItem{
			{
				ID:       1,
				MenuID:   100,
				Quantity: 2,
				Menu: models.Menu{
					ID:    100,
					Name:  "Jollof Rice",
					Price: 2500,
				},
			},
		},
		Payment: &models.Payment{
			ID:     1,
			Amount: 5000,
			Status: "completed",
			UserID: userID,
		},
	}

	response := OrderResponse(order)

	if response.ID != order.ID {
		t.Errorf("Expected ID %d, got %d", order.ID, response.ID)
	}

	if response.UserID != order.UserID {
		t.Errorf("Expected UserID %d, got %d", userID, response.UserID)
	}

	if response.UserID == nil {
		t.Error("Expected UserID to not be nil")
	}

	if response.Payment == nil {
		t.Error("Expected Payment to not be nil")
	}
}

func Test_ReservationResponse(t *testing.T) {
	timeSlot := datatypes.NewTime(20, 00, 0, 0)

	reservation := &models.Reservation{
		ID:       1,
		UserID:   10,
		TableID:  100,
		Date:     time.Now(),
		TimeSlot: timeSlot,
		Status:   "confirmed",
		User: models.User{
			ID:        10,
			FirstName: "John",
			LastName:  "Doe",
		},
		Table: models.Table{
			ID: 100,
		},
	}

	response := ReservationResponse(reservation)

	if response.ID != reservation.ID {
		t.Errorf("Expected ID %d, got %d", reservation.ID, response.ID)
	}

	if response.UserID != reservation.UserID {
		t.Errorf("Expected UserID %d, got %d", reservation.UserID, response.UserID)
	}

	if response.TimeSlot != reservation.TimeSlot.String() {
		t.Errorf("Expected TimeSlot %s, got %s", reservation.TimeSlot, response.TimeSlot)
	}

	if response.User.FirstName != "John" {
		t.Errorf("Expected User FirstName %s, got %s", "John", response.User.FirstName)
	}

	if response.Table.ID != reservation.TableID {
		t.Errorf("Expected Table ID %d, got %d", 100, response.Table.ID)
	}

	if response.Status != "confirmed" {
		t.Errorf("Expected Status %s, got %s", "confirmed", response.Status)
	}

	if response.Date != reservation.Date.Format("2006-01-02") {
		t.Error("Expected Date to not be zero")
	}
}
