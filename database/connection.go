
package database

import (
	"fmt"
	"log"

	"ecommerce/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    TranslateError: true,
})
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}

	// Verify database connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("MySQL ping failed:", err)
	}

	log.Println("MySQL database connected successfully!")

	return db
}
