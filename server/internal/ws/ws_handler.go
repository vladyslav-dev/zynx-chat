package ws

import (
	"log"
	"server/internal/message"
	"sort"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

func GetChannelID(t, senderID, recipientID, groupID *string) string {
	if *t == "group" {
		return "group-" + *groupID
	}
	// Sort the IDs to ensure a consistent channel ID regardless of the order
	ids := []string{*senderID, *recipientID}
	sort.Strings(ids)

	return "private-" + ids[0] + "-" + ids[1]
}

func (h *Handler) ServeWs(c *gin.Context, messageService message.Service) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection %v", err)
		return
	}

	client := &Client{
		hub:            h.hub,
		conn:           conn,
		send:           make(chan *message.MessageWrapper, 256),
		messageService: &messageService,
		channels:       make(map[string]bool),
	}

	client.hub.register <- client

	t := c.Query("type")
	senderID := c.Query("sender_id")
	recipientID := c.Query("recipient_id")
	groupID := c.Query("group_id")

	channelID := GetChannelID(&t, &senderID, &recipientID, &groupID)

	h.hub.AddClientToChannel(client, channelID)

	go client.writePump()
	go client.readPump()
}
