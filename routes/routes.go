package routes

import (
    "net/http"

    "ecommerce/config"
    "ecommerce/handlers"
    "ecommerce/middleware"
    "ecommerce/repositories"
    "ecommerce/services"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func SetupRoutes(
    router *gin.Engine,
    db *gorm.DB,
    cfg *config.Config,
) {

    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "success": true,
            "message": "Ecommerce API is running",
        })
    })

    // Repositories
    userRepository := repositories.NewUserRepository(db)
    categoryRepository := repositories.NewCategoryRepository(db)
    productRepository := repositories.NewProductRepository(db)
	cartRepository := repositories.NewCartRepository(db)
	orderRepository := repositories.NewOrderRepository(db)
	paymentRepository := repositories.NewPaymentRepository(db)

    // Services
    authService := services.NewAuthService(
        userRepository,
        cfg.JWTSecret,
        cfg.JWTExpiresInHours,
    )

    categoryService := services.NewCategoryService(
        categoryRepository,
    )

    productService := services.NewProductService(
        productRepository,
        categoryRepository,
    )

	cartService := services.NewCartService(
		cartRepository,
	)

	orderService := services.NewOrderService(
		orderRepository,
	)

	paymentService := services.NewPaymentService(
		paymentRepository,
	)

    // Handlers
    authHandler := handlers.NewAuthHandler(authService)
    categoryHandler := handlers.NewCategoryHandler(categoryService)
    productHandler := handlers.NewProductHandler(productService)
	cartHandler := handlers.NewCartHandler(cartService)
	orderHandler := handlers.NewOrderHandler(orderService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

    // Middleware
    authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)
    adminMiddleware := middleware.RequireRole("admin")

    // API versioning
    api := router.Group("/api/v1")
    {
        // Authentication
        auth := api.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)

            auth.GET(
                "/profile",
                authMiddleware,
                authHandler.GetProfile,
            )
        }

        // Category
        categories := api.Group("/categories")
        {
            // Public endpoints
            categories.GET("/", categoryHandler.GetCategories)
            categories.GET("/:id", categoryHandler.GetCategoryByID)

            // Admin-only endpoints
            categories.POST(
                "/",
                authMiddleware,
                adminMiddleware,
                categoryHandler.CreateCategory,
            )

            categories.PUT(
                "/:id",
                authMiddleware,
                adminMiddleware,
                categoryHandler.UpdateCategory,
            )

            categories.DELETE(
                "/:id",
                authMiddleware,
                adminMiddleware,
                categoryHandler.DeleteCategory,
            )
        }

        // Product
        products := api.Group("/products")
        {
            // Public endpoints
            products.GET("/", productHandler.GetProducts)
            products.GET("/:id", productHandler.GetProductByID)

            // Admin-only endpoints
            products.POST(
                "/",
                authMiddleware,
                adminMiddleware,
                productHandler.CreateProduct,
            )

            products.PUT(
                "/:id",
                authMiddleware,
                adminMiddleware,
                productHandler.UpdateProduct,
            )

            products.DELETE(
                "/:id",
                authMiddleware,
                adminMiddleware,
                productHandler.DeleteProduct,
            )
        }

		// Cart routes - authenticated users only
		cart := api.Group("/cart")
		cart.Use(authMiddleware)
		{
			cart.GET("", cartHandler.GetCart)

			cart.POST("/items", cartHandler.AddItem)

			cart.PUT("/items/:id", cartHandler.UpdateItem)

			cart.DELETE("/items/:id", cartHandler.RemoveItem)
		}
		// Create payment from an order
		orders := api.Group("/orders")
		orders.Use(authMiddleware)
		{
    		orders.POST("/checkout", orderHandler.Checkout)
    		orders.GET("", orderHandler.GetOrders)
    		orders.GET("/:id", orderHandler.GetOrderByID)
			orders.POST("/:id/payments", paymentHandler.CreatePayment)
			orders.POST("/:id/cancel", orderHandler.CancelOrder)
		}

		// Payment routes
		payments := api.Group("/payments")
		payments.Use(authMiddleware)
		{
    		payments.GET("/:id", paymentHandler.GetPayment)
    		payments.POST("/:id/simulate", paymentHandler.SimulatePayment)
		}

		adminOrders := api.Group("/admin/orders")
		adminOrders.Use(
    		authMiddleware,
    		middleware.RequireRole("admin"),
		)
		{
    		adminOrders.GET("", orderHandler.GetAllOrders)
    		adminOrders.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
		}

	}
}