package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"order-service/domain"
)

func InitDB() *gorm.DB {
	var err error
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = "host=localhost user=user password=password dbname=order_db port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	// Auto Migrate
	db.AutoMigrate(&domain.OrderModel{}, &domain.InvoiceModel{})
	
	return db
}
