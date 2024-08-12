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
	auth := r.Group("api/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.POST("/logout", userHandler.Logout)
		auth.GET("/refreshToken", userHandler.RefreshToken)
	}

	user := r.Group("api/user").Use(middlewares.AuthMiddleware())
	user.GET("/getAll", userHandler.GetAllUsers)
	user.GET("/getByGroupId", userHandler.GetUsersByGroupID)
	user.POST("/getByIds", userHandler.GetUsersByIDs)

}

func RegisterGroupRoutes(r *gin.Engine, groupHandler *group.Handler) {
	group := r.Group("api/group")
	group.Use(middlewares.AuthMiddleware())
	{
		group.POST("/create", groupHandler.CreateGroup)
		group.GET("/get", groupHandler.GetGroupByID)
		group.GET("/getAll", groupHandler.GetAllGroups)
		group.POST("/join", groupHandler.JoinGroup)
	}
}

func RegisterWs(r *gin.Engine, wsHandler *ws.Handler, messageHandler *message.Handler, messageService message.Service) {
	message := r.Group("api/message")
	message.Use(middlewares.AuthMiddleware())
	{
		message.POST("/private", messageHandler.GetPrivateMessages)
		message.POST("/group", messageHandler.GetGroupMessages)
	}
	r.GET("/ws/message", func(c *gin.Context) {
		wsHandler.ServeWs(c, messageService)
	})
}
