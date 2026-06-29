package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zaqimaulana/mycatalog-be/handlers"
	"github.com/zaqimaulana/mycatalog-be/middleware"
)

func SetupRouter() *gin.Engine {
	// Gunakan gin.New() agar kita bisa kontrol penuh urutan middleware
	// (tidak pakai gin.Default() yang auto-include logger bawaan gin)
	r := gin.New()
	r.Use(gin.Recovery()) // panic recovery tetap diperlukan
	r.Use(middleware.HTTPLogger())

	// ─── CORS Middleware ───────────────────────────────────────
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// ─── Init handlers ────────────────────────────────────────
	authHandler := handlers.NewAuthHandler()
	productHandler := handlers.NewProductHandler()
	cartHandler := handlers.NewCartHandler()
	orderHandler := handlers.NewOrderHandler()

	// ─── API v1 group ─────────────────────────────────────────
	v1 := r.Group("/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "service": "mycatalog-backend"})
		})

		// ── Auth routes (public) ──────────────────────────────
		auth := v1.Group("/auth")
		{
			auth.POST("/verify-token", authHandler.VerifyToken)
		}

		// ── Protected routes (butuh JWT) ──────────────────────
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// FCM token
			protected.PUT("/auth/fcm-token", authHandler.UpdateFCMToken)
			// Products
			products := protected.Group("/products")
			{
				products.GET("", productHandler.GetAll)
				products.GET("/:id", productHandler.GetByID)

				adminProducts := products.Group("")
				adminProducts.Use(middleware.AdminOnly())
				{
					adminProducts.POST("", productHandler.Create)
					adminProducts.PUT("/:id", productHandler.Update)
					adminProducts.DELETE("/:id", productHandler.Delete)
				}
			}

			// Cart
			cart := protected.Group("/cart")
			{
				cart.GET("", cartHandler.GetCart)               // GET    /v1/cart
				cart.POST("", cartHandler.AddToCart)            // POST   /v1/cart
				cart.PUT("/:id", cartHandler.UpdateCartItem)    // PUT    /v1/cart/:id
				cart.DELETE("/:id", cartHandler.RemoveCartItem) // DELETE /v1/cart/:id
				cart.DELETE("", cartHandler.ClearCart)          // DELETE /v1/cart
			}

			// Orders
			orders := protected.Group("/orders")
			{
				orders.POST("", orderHandler.CreateOrder)        // POST   /v1/orders      ← Flutter
				orders.POST("/checkout", orderHandler.Checkout)  // POST   /v1/orders/checkout ← backend cart
				orders.GET("", orderHandler.GetMyOrders)         // GET    /v1/orders
				orders.GET("/:id", orderHandler.GetOrderByID)    // GET    /v1/orders/:id
			}

			// Admin — order management
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminOnly())
			{
				admin.GET("/orders", orderHandler.GetAllOrders)                 // GET /v1/admin/orders
				admin.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus) // PUT /v1/admin/orders/:id/status
			}
		}
	}

	return r
}
