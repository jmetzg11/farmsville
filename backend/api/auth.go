package api

import (
	"errors"
	"farmsville/backend/models"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
		adminEmails := strings.Split(os.Getenv("ADMIN_EMAILS"), ",")
		for i, email := range adminEmails {
			adminEmails[i] = strings.TrimSpace(email)
		}

		newUser := models.User{
			Email:     email,
			Code:      code,
			ExpiresAt: expiresAt,
			CreatedAt: time.Now(),
			Admin:     slices.Contains(adminEmails, email),
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
			"name":            user.Name,
			"email":           user.Email,
			"admin":           user.Admin,
			"isAuthenticated": true,
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

func (h *Handler) AuthMe(c *gin.Context) {
	tokenString, err := c.Cookie("auth_token")
	fmt.Println("Auth token in AuthMe:", tokenString, "Error:", err)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Not authenticated",
		})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "fallback-secret-key"
		}
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid token",
		})
		return
	}
	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Invalid user ID in token",
			})
			return
		}
		user, err := h.getUserByID(uint(userID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found",
			})
			return
		}

		returnUser := gin.H{
			"name":            user.Name,
			"email":           user.Email,
			"admin":           user.Admin,
			"isAuthenticated": true,
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Authentication successful",
			"user":    returnUser,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid token",
		})
	}
}

func (h *Handler) Logout(c *gin.Context) {
	fmt.Println("Logout endpoint called")

	// Get the existing cookie first to check if it exists
	existingCookie, err := c.Cookie("auth_token")
	fmt.Println("Existing auth_token:", existingCookie, "Error:", err)

	isProduction := os.Getenv("GIN_MODE") == "release"
	c.SetCookie(
		"auth_token",
		"",
		-1, // negative maxAge immediately expires the cookie
		"/",
		"",
		isProduction,
		true,
	)

	afterCookie, err := c.Cookie("auth_token")
	fmt.Println("After clearing - auth_token:", afterCookie, "Error:", err)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}
