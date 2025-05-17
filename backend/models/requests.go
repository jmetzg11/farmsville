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
	Name        string `form:"title" binding:"required"`
	Description string `form:"description" binding:"required"`
	Quantity    int    `form:"quantity" binding:"required"`
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UpdateUserRequest struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Admin bool   `json:"admin"`
}

type AdminClaimItemRequest struct {
	UserID int `json:"userId" binding:"required"`
	ItemID int `json:"itemId" binding:"required"`
	Amount int `json:"amount" binding:"required"`
}
