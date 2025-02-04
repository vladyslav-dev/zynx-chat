package group

import (
	"net/http"
	"strconv"

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

func (h *Handler) CreateGroup(c *gin.Context) {
	var g CreateGroupReq
	if err := c.ShouldBindJSON(&g); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	res, err := h.Service.CreateGroup(c.Request.Context(), &g)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetGroupByID(c *gin.Context) {
	groupID := c.Query("id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Missing required query parameter id"})
		return
	}

	id, err := strconv.Atoi(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "id is not a number"})
		return
	}

	g, err := h.Service.GetGroupById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, g)
}

func (h *Handler) GetAllGroups(c *gin.Context) {
	gs, err := h.Service.GetAllGroups(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
	}

	c.JSON(http.StatusOK, gs)
}

func (h *Handler) JoinGroup(c *gin.Context) {
	var req JoinGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.JoinGroup(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
	}

	c.JSON(http.StatusOK, res)
}
