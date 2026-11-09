
package database

import (
	"log"

	"ecommerce/models"

	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
	)

	if err != nil {
		log.Fatal("Database migration failed:", err)
	}

	log.Println("Database migration completed successfully!")
}
