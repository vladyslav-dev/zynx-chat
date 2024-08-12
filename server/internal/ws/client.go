package ws

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"server/internal/message"
	m "server/internal/message"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub            *Hub
	conn           *websocket.Conn
	send           chan *m.MessageWrapper
	messageService *m.Service
	channels       map[string]bool // channelID as key, can be group or private
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, mes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		mes = bytes.TrimSpace(mes)

		var msg m.SendMessageReq
		if err := json.Unmarshal(mes, &msg); err != nil {
			log.Printf("Invalid message format: %v", err)

			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("syntax error at byte offset %d", e.Offset)
			}
		}

		c.handleMessage(msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			jsonMessage, err := m.MarshalMessageJSON(message)
			if err != nil {
				log.Printf("Error marshalling message: %v", err)
				return
			}
			w.Write(jsonMessage)

			n := len(c.send)

			for i := 0; i < n; i++ {
				jsonMessage, err := m.MarshalMessageJSON(<-c.send)
				if err != nil {
					log.Printf("Error marshalling message: %v", err)
					return
				}

				w.Write(jsonMessage)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(msg message.SendMessageReq) {
	ctx := context.Background()

	if msg.Type == "group" {
		sendMessageReq := &m.SendMessageReq{
			Type:     msg.Type,
			SenderID: msg.SenderID,
			GroupID:  msg.GroupID,
			Content:  msg.Content,
		}

		insertedMessage, err := (*c.messageService).SendMessage(ctx, sendMessageReq)
		if err != nil {
			log.Printf("Error saving message to the database: %v", err)
			return
		}

		strGroupID := strconv.Itoa(msg.GroupID)

		channelID := GetChannelID(&msg.Type, nil, nil, &strGroupID) // Get the correct channel ID based on the message
		c.hub.BroadcastMessage(BroadcastMessage{
			Message: &m.MessageWrapper{
				GroupMsg: &m.GroupMessageRes{
					ID:        insertedMessage.ID,
					Type:      insertedMessage.Type,
					SenderID:  insertedMessage.SenderID,
					GroupID:   *insertedMessage.GroupID,
					Content:   insertedMessage.Content,
					CreatedAt: insertedMessage.CreatedAt,
				},
			},
			ChannelID: channelID,
		})
	} else if msg.Type == "private" {
		sendMessageReq := &m.SendMessageReq{
			Type:        msg.Type,
			SenderID:    msg.SenderID,
			RecipientID: msg.RecipientID,
			Content:     msg.Content,
		}

		insertedMessage, err := (*c.messageService).SendMessage(ctx, sendMessageReq)
		if err != nil {
			log.Printf("Error saving message to the database: %v", err)
			return
		}
		strSenderID := strconv.Itoa(msg.SenderID)
		strRecipientID := strconv.Itoa(msg.RecipientID)

		channelID := GetChannelID(&msg.Type, &strSenderID, &strRecipientID, nil) // Get the correct channel ID based on the message

		c.hub.BroadcastMessage(BroadcastMessage{
			Message: &m.MessageWrapper{
				PrivateMsg: &m.PrivateMessageRes{
					ID:          insertedMessage.ID,
					Type:        insertedMessage.Type,
					SenderID:    insertedMessage.SenderID,
					RecipientID: *insertedMessage.RecipientID,
					Content:     insertedMessage.Content,
					CreatedAt:   insertedMessage.CreatedAt,
				},
			},
			ChannelID: channelID,
		})
	}
}
