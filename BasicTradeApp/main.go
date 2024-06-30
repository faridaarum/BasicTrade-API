package main

import (
	"BasicTradeApp/config"
	"BasicTradeApp/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Initialize database
	config.ConnectDB()

	// Initialize routes
	routes.InitializeRoutes(r)

	// Run the server
	r.Run(":80")
}
