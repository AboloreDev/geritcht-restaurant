package server

import (
	"net/http"
	"strconv"

	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	websockets "github.com/AboloreDev/geritcht-restaurant/internals/web-sockets"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// @Summary Subscribe to order status updates
// @Description Upgrade the connection to a WebSocket to receive real-time updates for an order's status. The authenticated user must own the order.
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 101 {string} string "Switching Protocols (WebSocket connection established)"
// @Failure 400 {object} utils.Response "Invalid order ID or order not found"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Failed to upgrade WebSocket connection"
// @Router /ws/orders/{id} [get]

// @Description Upgrade the connection to a WebSocket to receive real-time updates for an order's status.
// @Description
// @Description Example messages:
// @Description {"order_id":12,"status":"pending"}
// @Description {"order_id":12,"status":"preparing"}
// @Description {"order_id":12,"status":"ready"}
// @Description {"order_id":12,"status":"completed"}
func (s *Server) WebSocketHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid OrderID", err)
		return
	}
	orderID := uint(id)

	var order models.Order

	err = s.orderService.VerifyUserOrder(ctx.Request.Context(), userID, orderID)
	if err != nil {
		utils.BadRequest(ctx, "Order not found", err)
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to upgrade to WebSocket connection", err)
		return
	}

	client := &websockets.Client{
		OrderID: orderID,
		Hub:     s.hub,
		Conn:    conn,
		Send:    make(chan []byte, 256),
	}

	s.hub.Register <- client

	data := websockets.BuildMessageWithStatus(orderID, string(order.Status))
	client.Send <- data

	// Launch the read and write goroutines for the client
	go client.ReadPump()
	go client.WritePump()
}
