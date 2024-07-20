package main

import (
	"fmt"
	"os"

	"BasicTradeApp/config"
	"BasicTradeApp/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

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
		port = "8080" // Default port if not specified
	}

	// Run the server
	r.Run(fmt.Sprintf(":%s", port))
}
