package api

import (
	"farmsville/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUsers(c *gin.Context) {
	var dbUsers []models.User
	if err := h.db.Select("id, email, name, phone, admin").Find(&dbUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch users",
		})
		return
	}

	users := make([]models.UserResponse, len(dbUsers))
	for i, user := range dbUsers {
		users[i] = models.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Phone: user.Phone,
			Admin: user.Admin,
		}
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User updated",
	})
}

func (h *Handler) RemoveUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User removed",
	})
}

func (h *Handler) CreateUser(c *gin.Context) {
	var newUser models.CreateUserRequest
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Name:  newUser.Name,
		Email: newUser.Email,
		Phone: newUser.Phone,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User created",
	})
}
