package api

import (
	"farmsville/backend/auth"

	"gorm.io/gorm"
)

type Handler struct {
	db          *gorm.DB
	authService auth.Service
}

func NewHandler(db *gorm.DB, options ...interface{}) *Handler {
	// for production
	h := &Handler{
		db:          db,
		authService: auth.NewService(),
	}

	// for testing
	for _, option := range options {
		if authService, ok := option.(auth.Service); ok {
			h.authService = authService
		}
	}

	return h
}
