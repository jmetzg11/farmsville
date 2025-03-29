package models

type Email struct {
	Email string `json:"email"`
}

type AuthRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
