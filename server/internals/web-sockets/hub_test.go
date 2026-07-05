package websockets

import (
	"testing"
	"time"
)

func TestHub_Register(t *testing.T) {
	// Initialise the Hub
	hub := NewHub()

	// start the goroutine
	go hub.Run()

	// Initialise the client (orderID snd send channel)
	client := &Client{
		OrderID: 1,
		Send:    make(chan []byte, 10),
	}

	// If the client receives a signal, send to register channel
	hub.Register <- client

	// Wait for go routine to spin up
	// concurrnecy happens here, use mutex to prevent race condition
	// close when done
	waitForCondition(t, func() bool {
		hub.Mutex.RLock()
		defer hub.Mutex.RUnlock()
		return len(hub.Clients[1]) == 1
	}, 1*time.Second)

	if len(hub.Clients[1]) != 1 {
		t.Errorf("Expected 1 client, got %d", len(hub.Clients[1]))
	}

	if hub.Clients[1][0] != client {
		t.Errorf("expected client %v, got %v", client, hub.Clients[1][0])
	}
}

func TestHub_Unregister(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		OrderID: 1,
		Send:    make(chan []byte, 10),
	}

	hub.Register <- client
	waitForCondition(t, func() bool {
		hub.Mutex.RLock()
		defer hub.Mutex.RUnlock()
		return len(hub.Clients[1]) == 1
	}, 1*time.Second)

	hub.Unregister <- client
	waitForCondition(t, func() bool {
		hub.Mutex.RLock()
		defer hub.Mutex.RUnlock()
		_, exists := hub.Clients[1]
		return !exists
	}, time.Second)

	_, ok := <-client.Send
	if ok {
		t.Error("Expected closed channel, channel still opened!")
	}
}

func TestHub_Broadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create two clients
	client1 := &Client{
		OrderID: 1,
		Send:    make(chan []byte, 10),
	}

	client2 := &Client{
		OrderID: 2,
		Send:    make(chan []byte, 10),
	}

	// Register clients
	hub.Register <- client1
	hub.Register <- client2

	waitForCondition(t, func() bool {
		hub.Mutex.RLock()
		defer hub.Mutex.RUnlock()
		return len(hub.Clients[1]) == 1 && len(hub.Clients[2]) == 1
	}, time.Second)

	// Broadcast message
	message := []byte(`{"status":"preparing"}`)
	hub.Broadcast(1, message)
	hub.Broadcast(2, message)

	// Check if both clients received the message
	select {
	case received := <-client1.Send:
		if string(received) != string(message) {
			t.Errorf("Expected %s, got %s", message, received)
		}
	case <-time.After(time.Second):
		t.Error("Client 1 did not receive the broadcast message")
	}

	select {
	case received := <-client2.Send:
		if string(received) != string(message) {
			t.Errorf("Expected %s, got %s", message, received)
		}
	case <-time.After(time.Second):
		t.Error("Client 2 did not receive the broadcast message")
	}
}

func waitForCondition(t *testing.T, condition func() bool, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("condition not met within timeout")
}

func TestHub_Broadcast_RemovesSlowClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// unbuffered or tiny buffer — fills immediately
	client := &Client{OrderID: 1, Send: make(chan []byte)} // no buffer

	hub.Register <- client
	waitForCondition(t, func() bool {
		hub.Mutex.RLock()
		defer hub.Mutex.RUnlock()
		return len(hub.Clients[1]) == 1
	}, time.Second)

	// nobody reading from client.Send → it's "full" immediately
	hub.Broadcast(1, []byte("test"))

	// hub should detect default case and unregister this client
	waitForCondition(t, func() bool {
		hub.Mutex.RLock()
		defer hub.Mutex.RUnlock()
		_, exists := hub.Clients[1]
		return !exists
	}, time.Second)
}
