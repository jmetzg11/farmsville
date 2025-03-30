package models

type ClaimedItemWithUserName struct {
	ClaimedItem
	ItemName  string `json:"item_name"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}
