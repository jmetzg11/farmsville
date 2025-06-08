package database

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"log"

	"farmsville/backend/models"
)

var DB *gorm.DB

func Connect() error {

	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatal("Failed to create data directory", err)
	}

	var err error
	DB, err = gorm.Open(sqlite.Open("data/farmsville.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	DB.AutoMigrate(&models.Item{}, &models.ClaimedItem{}, &models.User{}, &models.Message{})

	seedDB(DB)
	return nil
}
