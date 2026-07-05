package websockets

import "encoding/json"

type OrderStatusMessage struct {
	OrderID uint   `json:"order_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// human readable message
var StatusMessages = map[string]string{
	"confirmed": "Order confirmed ✅",
	"preparing": "Your order is being prepared 👨‍🍳",
	"ready":     "Your order is ready for pickup 🎉",
	"completed": "Order completed. Enjoy your meal!",
	"cancelled": "Your order has been cancelled ❌",
}

func BuildMessageWithStatus(orderID uint, status string) []byte {
	msg := OrderStatusMessage{
		OrderID: orderID,
		Status:  status,
		Message: StatusMessages[status],
	}
	data, _ := json.Marshal(msg)
	return data
}
