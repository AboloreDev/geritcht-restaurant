package websockets

import (
	"log"
	"sync"
)

type Hub struct {
	Clients map[uint][]*Client
	Register chan *Client
	Unregister chan *Client
	Mutex sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[uint][]*Client),
		Register: make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			h.Clients[client.OrderID] = append(h.Clients[client.OrderID], client)
			h.Mutex.Unlock()
		
		case client := <-h.Unregister:
			h.Mutex.Lock()
			clients := h.Clients[client.OrderID]

			for i, c := range clients {
				if c == client {
					h.Clients[client.OrderID] = append( 
						clients[:i], 
						clients[i+1:]...,
					)
					close(client.Send)
					break
				}
			}

			if len(h.Clients[client.OrderID]) == 0 {
				delete(h.Clients, client.OrderID)
			}
			h.Mutex.Unlock()
		}
	}
}


func (h *Hub) Broadcast(orderID uint, message []byte) {
    h.Mutex.RLock()
    clients := h.Clients[orderID]
    h.Mutex.RUnlock()

    for _, client := range clients {
        select {
        case client.Send <- message:
			log.Printf("Sent message to client")
        default:
            h.Unregister <- client
			log.Printf("Client disconnected")
        }
    }
}