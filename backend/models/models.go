package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
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
	Password  string
	Phone     string
	Admin     bool
	Code      string
	ExpiresAt time.Time
	CreatedAt time.Time
	Claims    []ClaimedItem `gorm:"foreignKey:UserID" json:"-"`
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title" form:"title"`
	Message   string    `gorm:"type:text" json:"message" form:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type Blog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Blocks    []Block   `gorm:"foreignKey:BlogID;constraint:OnDelete:CASCADE" json:"blocks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Block struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	BlogID    uint      `json:"blog_id"`
	Type      string    `gorm:"size:50;not null" json:"type"`
	Order     int       `json:"order"`
	Text      string    `gorm:"type:text" json:"text,omitempty"`
	Media     string    `gorm:"type:text" json:"media,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
