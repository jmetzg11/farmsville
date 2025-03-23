package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"log"

	"farmsville/backend/models"
)

var DB *gorm.DB

func Connect() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("data/farmsville.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	DB.AutoMigrate(&models.Item{}, &models.ClaimedItem{})

	seedDB(DB)
	return nil
}
