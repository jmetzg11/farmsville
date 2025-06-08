package models

// auth and users
type Email struct {
	Email string `json:"email" binding:"required"`
}

type AuthRequest struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateAccountRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required"`
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

// items
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

type AdminClaimItemRequest struct {
	UserID int `json:"userId" binding:"required"`
	ItemID int `json:"itemId" binding:"required"`
	Amount int `json:"amount" binding:"required"`
}

// message
type PostMessageRequest struct {
	Title   string `json:"title" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type SendEmailRequest struct {
	Emails  []string `json:"emails" binding:"required"`
	Title   string   `json:"title" binding:"required"`
	Message string   `json:"message" binding:"required"`
}
