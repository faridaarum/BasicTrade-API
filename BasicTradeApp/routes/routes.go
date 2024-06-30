package routes

import (
	"BasicTradeApp/controllers"
	"BasicTradeApp/middlewares"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Protected routes
	authorized := r.Group("/")
	authorized.Use(middlewares.AuthMiddleware())
	{
		authorized.POST("/products", controllers.CreateProduct)
		authorized.PUT("/products/:id", controllers.UpdateProduct)
		authorized.DELETE("/products/:id", controllers.DeleteProduct)

		authorized.POST("/variants", controllers.CreateVariant)
		authorized.PUT("/variants/:id", controllers.UpdateVariant)
		authorized.DELETE("/variants/:id", controllers.DeleteVariant)
	}

	// Public product and variant routes
	r.GET("/products", controllers.GetProducts)
	r.GET("/products/:id", controllers.GetProductByID)
	r.GET("/variants", controllers.GetVariants)
	r.GET("/variants/:id", controllers.GetVariantByID)
}
