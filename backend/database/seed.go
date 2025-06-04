package database

import (
	"farmsville/backend/models"
	"log"
	"time"

	"gorm.io/gorm"
)

func seedDB(db *gorm.DB) {
	var count int64
	db.Model(&models.Item{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded")
		return
	}

	// Create sample users first
	users := []models.User{
		{
			Email:     "john.doe@example.com",
			Name:      "John Doe",
			Admin:     false,
			Code:      "user1code",
			ExpiresAt: time.Now().AddDate(1, 0, 0),
			CreatedAt: time.Now(),
		},
		{
			Email:     "jane.smith@example.com",
			Name:      "Jane Smith",
			Admin:     false,
			Code:      "user2code",
			ExpiresAt: time.Now().AddDate(1, 0, 0),
			CreatedAt: time.Now(),
		},
		{
			Email:     "mark.johnson@example.com",
			Name:      "Mark Johnson",
			Admin:     false,
			Code:      "user3code",
			ExpiresAt: time.Now().AddDate(1, 0, 0),
			CreatedAt: time.Now(),
		},
		{
			Email:     "sarah.williams@example.com",
			Name:      "Sarah Williams",
			Admin:     false,
			Code:      "user4code",
			ExpiresAt: time.Now().AddDate(1, 0, 0),
			CreatedAt: time.Now(),
		},
		{
			Email:     "robert.brown@example.com",
			Name:      "Robert Brown",
			Admin:     true, // Admin user
			Code:      "admincode",
			ExpiresAt: time.Now().AddDate(1, 0, 0),
			CreatedAt: time.Now(),
		},
	}

	// Insert users
	for i := range users {
		if i == 4 {
			users[i].SetPassword("admin")
		} else {
			users[i].SetPassword("password")
		}
		result := db.Create(&users[i])
		if result.Error != nil {
			log.Printf("Error creating user: %v", result.Error)
		}
	}

	// Create items (same as before)
	items := []models.Item{
		{
			Name:         "Organic Tomatoes",
			Description:  "Fresh organic tomatoes from local farms",
			Quantity:     100,
			RemainingQty: 100,
			Active:       true,
		},
		{
			Name:         "Fresh Eggs",
			Description:  "Free-range eggs from local chickens",
			Quantity:     50,
			RemainingQty: 50,
			Active:       true,
		},
		{
			Name:         "Honey",
			Description:  "Raw unfiltered honey, 16oz jars",
			Quantity:     30,
			RemainingQty: 30,
			Active:       true,
		},
		{
			Name:         "Lettuce",
			Description:  "Organic lettuce heads",
			Quantity:     80,
			RemainingQty: 80,
			Active:       true,
		},
		{
			Name:         "Apples",
			Description:  "Gala apples, 5lb bags",
			Quantity:     40,
			RemainingQty: 40,
			Active:       true,
		},
	}

	// Insert items
	for i := range items {
		result := db.Create(&items[i])
		if result.Error != nil {
			log.Printf("Error creating item: %v", result.Error)
		}
	}

	// Create sample claimed items with UserID instead of User email
	claimedItems := []models.ClaimedItem{
		{
			ItemID:    1, // Organic Tomatoes
			UserID:    1, // John Doe
			Quantity:  5,
			CreatedAt: time.Now(),
			Active:    true,
		},
		{
			ItemID:    2, // Fresh Eggs
			UserID:    2, // Jane Smith
			Quantity:  2,
			CreatedAt: time.Now(),
			Active:    true,
		},
		{
			ItemID:    3, // Honey
			UserID:    3, // Mark Johnson
			Quantity:  1,
			CreatedAt: time.Now(),
			Active:    true,
		},
		{
			ItemID:    1, // Organic Tomatoes
			UserID:    4, // Sarah Williams
			Quantity:  3,
			CreatedAt: time.Now(),
			Active:    true,
		},
		{
			ItemID:    5, // Apples
			UserID:    5, // Robert Brown
			Quantity:  2,
			CreatedAt: time.Now(),
			Active:    true,
		},
	}

	// Insert claimed items
	for i := range claimedItems {
		result := db.Create(&claimedItems[i])
		if result.Error != nil {
			log.Printf("Error creating claimed item: %v", result.Error)
		}

		// Update the remaining quantity of the associated item
		var item models.Item
		db.First(&item, claimedItems[i].ItemID)
		item.RemainingQty -= claimedItems[i].Quantity
		db.Save(&item)
	}

	log.Println("Database seeded successfully!")
}
