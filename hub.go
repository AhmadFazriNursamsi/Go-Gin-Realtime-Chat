package main

import (
	"encoding/json"
	"fmt"
	"myapi/models"
	"sync"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client
}

var hub = &Hub{
	clients: make(map[string]*Client),
}

// gunakan unique key per client, misalnya pointer address
func clientKey(c *Client) string {
	return fmt.Sprintf("%p", c)
}

// register client baru
func (h *Hub) register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[clientKey(c)] = c
}

// unregister client (disconnect)
func (h *Hub) unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, clientKey(c))
	close(c.Send) // ✅ close sekali saja
}

// cek apakah client memiliki room tertentu
func (c *Client) HasRoom(roomID uint) bool {
	for _, r := range c.Rooms {
		if r == roomID {
			return true
		}
	}
	return false
}

// broadcast ke semua client yang ada di room yg sesuai
func (h *Hub) broadcast(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var temp models.Messages
	_ = json.Unmarshal(message, &temp)

	for key, client := range h.clients {
		// jika pesan ada RoomID → kirim hanya ke client yang punya room tersebut
		if temp.RoomID != 0 && !client.HasRoom(temp.RoomID) {
			continue
		}

		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.clients, key)
		}
	}
}
