package database

import (
	"farmsville/backend/models"
	"log"

	"gorm.io/gorm"
)

func seedDB(db *gorm.DB) {
	var count int64
	db.Model(&models.Item{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded")
		return
	}

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

	// Create sample claimed items
	claimedItems := []models.ClaimedItem{
		{
			ItemID:   1, // Organic Tomatoes
			User:     "john.doe@example.com",
			Quantity: 5,
			Active:   true,
		},
		{
			ItemID:   2, // Fresh Eggs
			User:     "jane.smith@example.com",
			Quantity: 2,
			Active:   true,
		},
		{
			ItemID:   3, // Honey
			User:     "mark.johnson@example.com",
			Quantity: 1,
			Active:   true,
		},
		{
			ItemID:   1, // Organic Tomatoes
			User:     "sarah.williams@example.com",
			Quantity: 3,
			Active:   true,
		},
		{
			ItemID:   5, // Apples
			User:     "robert.brown@example.com",
			Quantity: 2,
			Active:   true,
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
