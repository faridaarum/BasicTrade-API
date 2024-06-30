package main

import (
	"BasicTradeApp/config"
	"BasicTradeApp/routes"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Initialize database
	config.ConnectDB()

	// Initialize routes
	routes.InitializeRoutes(r)

	// Simple health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "80" // Default port if not specified
	}

	// Run the server
	r.Run(fmt.Sprintf(":%s", port))
}
