package main

import (
    "log"

    "ecommerce/config"
    "ecommerce/database"
    "ecommerce/routes"

    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.LoadConfig()

    if len(cfg.JWTSecret) < 32 {
        log.Fatal("JWT_SECRET must contain at least 32 characters")
    }

    gin.SetMode(cfg.GinMode)

    db := database.ConnectDB(cfg)

    database.MigrateDB(db)

    router := gin.Default()

    routes.SetupRoutes(router, db, cfg)

    log.Println("Ecommerce API running on port:", cfg.AppPort)

    if err := router.Run(":" + cfg.AppPort); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}