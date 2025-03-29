package api

import (
	"errors"
	"farmsville/backend/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) SendAuth(c *gin.Context) {
	var req models.Email
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	authCode, err := h.authService.GenerateRandomCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	err = h.updateOrCreateUser(req.Email, authCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.authService.SendEmailWithAuthCode(req.Email, authCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Authentication email sent",
	})

}

func (h *Handler) updateOrCreateUser(email, code string) error {
	expiresAt := time.Now().Add(15 * time.Minute)
	var existingUser models.User
	result := h.db.Where("email = ?", email).First(&existingUser)

	if result.Error == nil {
		// User exists, update their code and expiration
		existingUser.Code = code
		existingUser.ExpiresAt = expiresAt
		if err := h.db.Save(&existingUser).Error; err != nil {
			return fmt.Errorf("failed to update existing user: %w", err)
		}
		return nil
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// User does not exist, create a new one
		newUser := models.User{
			Email:     email,
			Code:      code,
			ExpiresAt: expiresAt,
			CreatedAt: time.Now(),
			Admin:     false,
		}
		if err := h.db.Create(&newUser).Error; err != nil {
			return fmt.Errorf("failed to create new user: %w", err)
		}
		return nil
	} else {
		return fmt.Errorf("database error: %w", result.Error)
	}
}

func (h *Handler) VerifyAuth(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.getUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Code == req.Code && user.ExpiresAt.After(time.Now()) {
		token, err := h.authService.GenerateJWT(user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
			return
		}

		isProduction := os.Getenv("GIN_MODE") == "release"

		maxAge := 90 * 24 * 60 * 60
		c.SetCookie(
			"auth_token",
			token,
			maxAge,
			"/",
			"",
			isProduction,
			true,
		)
		returnUser := gin.H{
			"name":  user.Name,
			"email": user.Email,
			"admin": user.Admin,
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Authentication successful",
			"user":    returnUser,
		})
		return
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid code or code expired"})
		return
	}
}

func (h *Handler) getUserByEmail(email string) (models.User, error) {
	var user models.User
	result := h.db.Where("email = ?", email).First(&user)
	if result.Error == nil {
		return user, nil
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.User{}, fmt.Errorf("user not found")
	} else {
		return models.User{}, fmt.Errorf("database error: %w", result.Error)
	}
}
