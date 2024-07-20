package routes

import (
	"BasicTradeApp/controllers"
	"BasicTradeApp/middlewares"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)
	r.POST("/refresh", controllers.RefreshToken)

	// Public product routes
	r.GET("/products", controllers.GetProducts)
	r.GET("/products/:product_id", controllers.GetProductByID)

	// Public variant routes
	r.GET("/variants", controllers.GetAllVariants)
	r.GET("/products/:product_id/variants", controllers.GetVariants)
	r.GET("/products/:product_id/variants/:variant_id", controllers.GetVariantByID)

	// Protected routes
	authorized := r.Group("/")
	authorized.Use(middlewares.AuthMiddleware())
	{
		authorized.POST("/products", controllers.CreateProduct)
		authorized.PUT("/products/:product_id", controllers.UpdateProduct)
		authorized.DELETE("/products/:product_id", controllers.DeleteProduct)

		// Specific routes for variants within a product
		authorized.POST("/products/:product_id/variants", controllers.CreateVariant)
		authorized.PUT("/products/:product_id/variants/:variant_id", controllers.UpdateVariant)
		authorized.DELETE("/products/:product_id/variants/:variant_id", controllers.DeleteVariant)
	}
}
