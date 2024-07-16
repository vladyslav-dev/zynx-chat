package ws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"server/internal/message"

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
	send           chan []byte
	messageService *message.Service
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
		_, m, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		fmt.Printf("msg %d", m)

		m = bytes.TrimSpace(m)

		var msg message.SendMessageReq
		if err := json.Unmarshal(m, &msg); err != nil {
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
			fmt.Println("message: ", message)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)

			for i := 0; i < n; i++ {
				w.Write(<-c.send)
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
		sendMessageReq := &message.SendMessageReq{
			Type:     msg.Type,
			SenderID: msg.SenderID,
			GroupID:  msg.GroupID,
			Content:  msg.Content,
		}

		res, err := (*c.messageService).SendMessage(ctx, sendMessageReq)
		if err != nil {
			log.Printf("Error saving message to the database: %v", err)
			return
		}

		strGroupID := strconv.Itoa(msg.GroupID)

		channelID := GetChannelID(&msg.Type, nil, nil, &strGroupID) // Get the correct channel ID based on the message

		c.hub.BroadcastMessage(BroadcastMessage{
			Message:   []byte(res.Content),
			ChannelID: channelID,
		})
	} else if msg.Type == "private" {
		sendMessageReq := &message.SendMessageReq{
			Type:        msg.Type,
			SenderID:    msg.SenderID,
			RecipientID: msg.RecipientID,
			Content:     msg.Content,
		}

		res, err := (*c.messageService).SendMessage(ctx, sendMessageReq)
		if err != nil {
			log.Printf("Error saving message to the database: %v", err)
			return
		}
		log.Printf("msg")
		log.Print(msg)
		strSenderID := strconv.Itoa(msg.SenderID)
		strRecipientID := strconv.Itoa(msg.RecipientID)

		log.Println("msg.SenderID")
		log.Println(msg.SenderID)
		log.Println(strSenderID)
		log.Println(strRecipientID)

		channelID := GetChannelID(&msg.Type, &strSenderID, &strRecipientID, nil) // Get the correct channel ID based on the message

		c.hub.BroadcastMessage(BroadcastMessage{
			Message:   []byte(res.Content),
			ChannelID: channelID,
		})
	}
}
