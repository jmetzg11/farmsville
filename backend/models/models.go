package models

import (
	"time"

	"gorm.io/gorm"
)

// Item represents an inventory item
type Item struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Name         string        `gorm:"size:255;not null" json:"name"`
	Description  string        `gorm:"type:text" json:"description"`
	DateCreated  time.Time     `json:"date_created"`
	Quantity     int           `json:"quantity"`
	RemainingQty int           `json:"remaining_quantity"`
	ClaimedItems []ClaimedItem `gorm:"foreignKey:ItemID" json:"claimed_items,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Active       bool          `json:"active"`
}

// ClaimedItem represents an item claimed by a user
type ClaimedItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ItemID      uint      `json:"item_id"`
	User        string    `gorm:"size:255;not null" json:"user"`
	Quantity    int       `json:"quantity"`
	DateClaimed time.Time `json:"date_claimed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Active      bool      `json:"active"`
}

// BeforeCreate is a GORM hook that sets the DateCreated field before creating a new Item
func (i *Item) BeforeCreate(tx *gorm.DB) error {
	i.DateCreated = time.Now()
	i.RemainingQty = i.Quantity
	return nil
}

// BeforeCreate is a GORM hook that sets the DateClaimed field before creating a new ClaimedItem
func (ci *ClaimedItem) BeforeCreate(tx *gorm.DB) error {
	ci.DateClaimed = time.Now()
	return nil
}
