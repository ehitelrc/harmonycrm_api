package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn    *websocket.Conn
	Channel string // p.ej. "case:123" o "agent:45"
	Send    chan []byte
}

type BroadcastMsg struct {
	Channel string
	Payload []byte
}

type Hub struct {
	mu         sync.RWMutex
	clients    map[string]map[*Client]bool // channel -> set de clients
	register   chan *Client
	unregister chan *Client
	broadcast  chan BroadcastMsg
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan BroadcastMsg),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			if h.clients[c.Channel] == nil {
				h.clients[c.Channel] = make(map[*Client]bool)
			}
			h.clients[c.Channel][c] = true
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			if set, ok := h.clients[c.Channel]; ok {
				if _, ok := set[c]; ok {
					delete(set, c)
					close(c.Send)
					_ = c.Conn.Close()
					if len(set) == 0 {
						delete(h.clients, c.Channel)
					}
				}
			}
			h.mu.Unlock()

		case m := <-h.broadcast:
			h.mu.RLock()
			set := h.clients[m.Channel]
			for cli := range set {
				select {
				case cli.Send <- m.Payload:
				default:
					// cliente congestionado -> desconectar
					h.mu.RUnlock()
					h.mu.Lock()
					delete(set, cli)
					close(cli.Send)
					_ = cli.Conn.Close()
					h.mu.Unlock()
					h.mu.RLock()
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastJSON(channel string, payload []byte) {
	h.broadcast <- BroadcastMsg{Channel: channel, Payload: payload}
}
