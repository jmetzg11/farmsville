package models

import (
	"time"
)

// Item represents an inventory item
type Item struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Name         string        `gorm:"size:255;not null" json:"name"`
	Description  string        `gorm:"type:text" json:"description"`
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
	User      string    `gorm:"size:255;not null" json:"user"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	Active    bool      `json:"active"`
}

type AuthCode struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"index"`
	Code      string
	ExpiresAt time.Time
	CreatedAt time.Time
}
