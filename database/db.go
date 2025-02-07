package database

import (
	"log"
	"notification-service/config"
	"notification-service/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection and performs migrations.
func InitDB() {
	// Build DSN (Data Source Name)
	dsn := config.AppConfig.DBUser + ":" +
		config.AppConfig.DBPassword + "@tcp(" +
		config.AppConfig.DBHost + ":" +
		config.AppConfig.DBPort + ")/" +
		config.AppConfig.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"

	// Open a connection to the database
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Automatically migrate the database schema
	err = DB.AutoMigrate(&models.SMSRequest{}, &models.BlacklistedNumbers{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	log.Println("Database initialized successfully")
}

// GetDB returns the database instance.
func GetDB() *gorm.DB {
	return DB
}
