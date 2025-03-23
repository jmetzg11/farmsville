package models

type ClaimedItemWithItemName struct {
	ClaimedItem
	ItemName string `json:"item_name"`
}
