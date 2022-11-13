package ws

import (
	"github.com/gorilla/websocket"
	"service.chat/internal/svc"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	svcCtx *svc.ServiceContext
}

func NewHub(svcCtx *svc.ServiceContext) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		svcCtx:     svcCtx,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case message := <-h.svcCtx.ChannelChat:
			for client := range h.clients {
				if (message.JwtId != "" && message.JwtId == client.JwtId) || (message.UserId == client.UserId) {
					client.conn.WriteMessage(websocket.TextMessage, message.Content)
					//select {
					//case client.send <- message.Content:
					//default:
					//	close(client.send)
					//	delete(h.clients, client)
					//}
				}
			}
		}

	}
}
