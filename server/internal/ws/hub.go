package ws

import (
	"fmt"
	"server/internal/message"
	"sync"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan BroadcastMessage
	register   chan *Client
	unregister chan *Client
	channels   map[string]map[*Client]bool // channelID as key
	mu         sync.RWMutex
}

type BroadcastMessage struct {
	Message   *message.MessageWrapper
	ChannelID string
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		channels:   make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

			h.RemoveFromAllChannels(client)
			h.mu.Unlock()
		case broadcastMessage := <-h.broadcast:
			h.mu.Lock()

			logs := fmt.Sprintf("Existing channels: %v\nIncoming ChannelID: %s\nClients in Channel: %v", h.channels, broadcastMessage.ChannelID, h.channels[broadcastMessage.ChannelID])
			fmt.Println(logs)

			clients := h.channels[broadcastMessage.ChannelID]
			h.mu.Unlock()
			for client := range clients {
				select {
				case client.send <- broadcastMessage.Message:

				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) RemoveFromAllChannels(client *Client) {
	for channelID := range client.channels {
		if clients, ok := h.channels[channelID]; ok {
			delete(clients, client)
		}
	}
}

func (h *Hub) AddClientToChannel(client *Client, channelID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Ensure the channels map is initialized
	if h.channels == nil {
		h.channels = make(map[string]map[*Client]bool)
	}

	// Initialize the map for the specific channel if it doesn't exist
	if h.channels[channelID] == nil {
		h.channels[channelID] = make(map[*Client]bool)
	}

	// Add the client to the channel
	h.channels[channelID][client] = true

	// Optionally, you might want to maintain a reverse mapping
	client.channels[channelID] = true
}

func (h *Hub) BroadcastMessage(message BroadcastMessage) {
	h.broadcast <- message
}
