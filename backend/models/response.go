package models

import "time"

type ClaimedItemWithUserName struct {
	ClaimedItem
	ItemName  string `json:"item_name"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Admin bool   `json:"admin"`
}

type BlogTitles struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}
