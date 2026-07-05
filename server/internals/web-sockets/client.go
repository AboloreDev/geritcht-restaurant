package websockets

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait       = 60 * time.Second
	maxMessageSize = 512
	pingPeriod     = pongWait * 9 / 10
	writeWait      = 10 * time.Second
)

type Client struct {
	OrderID uint
	Hub     *Hub
	Conn    *websocket.Conn
	Send    chan []byte
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		c.Conn.Close()
		ticker.Stop()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println(err)
				return
			}

			w.Write(msg)

			for i := 0; i < len(c.Send); i++ {
				w.Write(msg)
			}

			err = w.Close()
			if err != nil {
				log.Println(err)
				return
			}
		case <-ticker.C:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Conn.Close()
		c.Hub.Unregister <- c
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}
