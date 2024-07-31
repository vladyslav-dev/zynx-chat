package router

import (
	"server/internal/group"
	"server/internal/message"
	"server/internal/user"
	"server/internal/ws"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, userHandler *user.Handler) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.GET("/logout", userHandler.Logout)
		auth.GET("/refresh-token", userHandler.RefreshToken)
	}
	// r.GET("api/getAllUsers", userHandler.GetAllUsers)
}

func RegisterGroupRoutes(r *gin.Engine, groupHandler *group.Handler) {
	group := r.Group("/group")
	group.Use(middlewares.AuthMiddleware())
	{
		group.POST("/create", groupHandler.CreateGroup)
		group.GET("/get-all", groupHandler.GetAllGroups)
		group.POST("/join", groupHandler.JoinGroup)
	}
}

func RegisterWs(r *gin.Engine, wsHandler *ws.Handler, messageHandler *message.Handler, messageService message.Service) {
	message := r.Group("/message")
	{
		message.POST("/private", messageHandler.GetPrivateMessages)
		message.POST("/group", messageHandler.GetGroupMessages)
	}
	r.GET("/ws/message", func(c *gin.Context) {
		wsHandler.ServeWs(c, messageService)
	})
}
