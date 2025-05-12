package models

type Email struct {
	Email string `json:"email"`
}

type AuthRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ClaimRequest struct {
	ItemID   int `json:"itemId"`
	Quantity int `json:"quantity"`
}

type ItemRequest struct {
	ID int `json:"id"`
}

type CreateItemRequest struct {
	Name        string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required"`
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
