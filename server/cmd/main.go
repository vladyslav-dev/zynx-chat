package main

import (
	"fmt"
	"log"
	"os"
	"server/db"
	"server/internal/group"
	"server/internal/message"
	"server/internal/user"
	"server/internal/ws"
	"server/middlewares"
	"server/router"

	"github.com/gin-gonic/gin"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database connection: %s", err)
	}

	fmt.Print("Database connected successfully")

	// Initialize ws hub
	hub := ws.NewHub()

	// Initialize repositories
	userRep := user.NewRepository(dbConn.GetDB())
	groupRep := group.NewRepository(dbConn.GetDB())
	messageRep := message.NewRepository(dbConn.GetDB())

	// Initialize services
	userSvc := user.NewService(userRep)
	groupSvc := group.NewService(groupRep, userRep)
	messageSvc := message.NewService(messageRep)

	// Initialize handlers
	userHandler := user.NewHandler(userSvc)
	groupHandler := group.NewHandler(groupSvc)
	messageHandler := message.NewHandler(messageSvc)
	wsHandler := ws.NewHandler(hub)

	// Initialize the router
	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "The server is running",
		})
	})

	// Register routes
	router.RegisterUserRoutes(r, userHandler)
	router.RegisterGroupRoutes(r, groupHandler)
	router.RegisterWs(r, wsHandler, messageHandler, messageSvc)

	// Run ws hub
	go hub.Run()

	// hub := ws.NewHub()
	// wsHandler := ws.NewHandler(hub)
	// go hub.Run()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run("0.0.0.0:" + port)
}
