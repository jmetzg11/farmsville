package models

import (
	"time"
)

// think about adding composite indexes

// Item represents an inventory item
type Item struct {
	ID           uint          `gorm:"primaryKey" json:"id" form:"id"`
	Name         string        `gorm:"size:255;not null" json:"name" form:"name"`
	Description  string        `gorm:"type:text" json:"description" form:"description"`
	PhotoPath    string        `gorm:"type:text" json:"photo_path" form:"photo_path"`
	Quantity     int           `json:"quantity" form:"quantity"`
	RemainingQty int           `json:"remaining_quantity" form:"remaining_quantity"`
	ClaimedItems []ClaimedItem `gorm:"foreignKey:ItemID" json:"claimed_items,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	Active       bool          `gorm:"index"json:"active"`
}

// ClaimedItem represents an item claimed by a user
type ClaimedItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ItemID    uint      `json:"item_id"`
	UserID    uint      `json:"user_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	Active    bool      `gorm:"index" json:"active"`
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
