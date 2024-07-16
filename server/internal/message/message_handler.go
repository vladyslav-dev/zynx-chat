package message

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

// func (h *Handler) SendMessage(c *gin.Context, hub *ws.Hub) {
// 	var req SendMessageReq
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	res, err := h.Service.SendMessage(c.Request.Context(), &req)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	hub.BroadcastMessage(ws.BroadcastMessage{GroupID: res.GroupID, Message: []byte(res.Content)})

// 	c.JSON(http.StatusOK, res)
// }

func (h *Handler) GetPrivateMessages(c *gin.Context) {
	var req GetPrivateMessagesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.GetPrivateMessages(c, &req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetGroupMessages(c *gin.Context) {
	var req GetGroupMessagesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.GetGroupMessages(c, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
