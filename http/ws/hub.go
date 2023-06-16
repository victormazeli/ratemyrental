package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"strings"
)

// WebSocketMessage represents a message received via WebSocket
type WebSocketMessage struct {
	Sender    uint   `json:"sender"`
	Recipient uint   `json:"recipient"`
	Content   string `json:"content"`
	ChannelID uint   `json:"channel_id"`
}

// Hub manages the WebSocket connections and broadcasts messages
type Hub struct {
	clients    map[uint]map[*websocket.Conn]bool // Map of channel IDs to client connections
	channels   map[string]uint                   // Map of sender:recipient to channel ID
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	broadcast  chan WebSocketMessage
}

var (
	HubInstance = NewHub() // Exported package-level variable
)

// NewHub initializes a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]map[*websocket.Conn]bool),
		channels:   make(map[string]uint),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		broadcast:  make(chan WebSocketMessage),
	}
}

// Run starts the WebSocket hub
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			// Get the channel ID from the connection's URL path
			// Extract the channel ID from the URL path if needed
			channelID, err := extractChannelIDFromURL(conn)

			if err != nil {
				log.Println("Failed to parse channel ID:", err)
				conn.Close()
				continue
			}

			// Add the connection to the corresponding channel's client list
			if _, ok := h.clients[channelID]; !ok {
				h.clients[channelID] = make(map[*websocket.Conn]bool)
			}
			h.clients[channelID][conn] = true

		case conn := <-h.unregister:
			// Remove the connection from the corresponding channel's client list
			for channelID, clients := range h.clients {
				if _, ok := clients[conn]; ok {
					delete(clients, conn)
					conn.Close()

					// If no clients are left in the channel, delete the channel
					if len(clients) == 0 {
						delete(h.clients, channelID)
						delete(h.channels, string(channelID))
					}
				}
			}

		case message := <-h.broadcast:
			// Broadcast the message to all clients in the corresponding channel
			if clients, ok := h.clients[message.ChannelID]; ok {
				for conn := range clients {
					err := conn.WriteJSON(message)
					if err != nil {
						log.Println("Error broadcasting message:", err)
						conn.Close()
						delete(clients, conn)
					}
				}

				// If no clients are left in the channel, delete the channel
				if len(clients) == 0 {
					delete(h.clients, message.ChannelID)
					delete(h.channels, fmt.Sprintf("%v:%v", message.Sender, message.Recipient))
				}
			}
		}
	}
}

func extractChannelIDFromURL(conn *websocket.Conn) (uint, error) {
	// Get the remote address from the connection
	remoteAddr := conn.RemoteAddr().String()

	// Extract the channel ID from the remote address
	parts := strings.Split(remoteAddr, "/")
	channelIDStr := parts[len(parts)-1]
	log.Printf("channel ID: %v", channelIDStr)
	channelID, err := strconv.ParseUint(channelIDStr, 10, 32)
	if err != nil {
		log.Println("Failed to parse channel ID:", err)
		return 0, err
	}

	return uint(channelID), nil
}
