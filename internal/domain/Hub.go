package domain

import (
	"bcv/models"
	"sync"
	"github.com/gofiber/contrib/v3/websocket"
)

type Hub struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan []models.BinanceRate
	Mu        sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan []models.BinanceRate),
	}
}

func (h *Hub) Run() {
	for rates := range h.Broadcast {
		h.Mu.Lock()
		for client := range h.Clients {
			err := client.WriteJSON(rates)
			if err != nil {
				client.Close()
				delete(h.Clients, client)
			}
		}
		h.Mu.Unlock()
	}
}
