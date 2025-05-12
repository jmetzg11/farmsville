package api

import (
	"errors"
	"farmsville/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	var updateUser models.UpdateUserRequest
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := h.db.First(&user, updateUser.ID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	user.Name = updateUser.Name
	user.Email = updateUser.Email
	user.Phone = updateUser.Phone
	user.Admin = updateUser.Admin

	if result := h.db.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
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
