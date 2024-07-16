package router

import (
	"server/internal/group"
	"server/internal/message"
	"server/internal/user"
	"server/internal/ws"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, userHandler *user.Handler) {
	r.POST("api/register", userHandler.CreateUser)
	r.POST("api/login", userHandler.Login)
	r.GET("api/logout", userHandler.Logout)
	r.GET("api/getAllUsers", userHandler.GetAllUsers)
}

func RegisterGroupRoutes(r *gin.Engine, groupHandler *group.Handler) {
	r.POST("api/createGroup", groupHandler.CreateGroup)
	r.GET("api/getAllGroups", groupHandler.GetAllGroups)
	r.POST("api/joinGroup", groupHandler.JoinGroup)
	// r.POST("api/joinGroup")
}

func RegisterWs(r *gin.Engine, wsHandler *ws.Handler, messageHandler *message.Handler, messageService message.Service) {
	r.POST("api/private-message", messageHandler.GetPrivateMessages)
	r.POST("api/group-message", messageHandler.GetGroupMessages)
	r.GET("/ws/message", func(c *gin.Context) {
		wsHandler.ServeWs(c, messageService)
	})
}
