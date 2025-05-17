package models

import (
	"time"
)

// Item represents an inventory item
type Item struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Name         string        `gorm:"size:255;not null" json:"name"`
	Description  string        `gorm:"type:text" json:"description"`
	PhotoPath    string        `gorm:"type:text" json:"photo_path"`
	Quantity     int           `json:"quantity"`
	RemainingQty int           `json:"remaining_quantity"`
	ClaimedItems []ClaimedItem `gorm:"foreignKey:ItemID" json:"claimed_items,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	Active       bool          `json:"active"`
}

// ClaimedItem represents an item claimed by a user
type ClaimedItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ItemID    uint      `json:"item_id"`
	UserID    uint      `json:"user_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	Active    bool      `json:"active"`
}

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"index"`
	Name      string
	Phone     string
	Admin     bool
	Code      string
	ExpiresAt time.Time
	CreatedAt time.Time
	Claims    []ClaimedItem `gorm:"foreignKey:UserID" json:"-"`
}
