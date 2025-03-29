package api

import (
	"farmsville/backend/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetItems(c *gin.Context) {
	var items []models.Item
	if err := h.db.Where("active = ?", true).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch items",
		})
		return
	}

	var claimedItems []models.ClaimedItemWithItemName
	if err := h.db.Raw(`
		SELECT
			claimed_items.*,
			items.name as item_name
		FROM
			claimed_items
		JOIN
			items ON claimed_items.item_id = items.id
		WHERE
			claimed_items.active = ?
	`, true).Scan(&claimedItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch claimed items",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":        items,
		"claimedItems": claimedItems,
	})
}

func (h *Handler) MakeClaim(c *gin.Context) {
	fmt.Println("MakeClaim")
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}
	fmt.Println(user)
	currentUser, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user information",
		})
		return
	}

	fmt.Println(currentUser)

	var claim models.Claim
	if err := c.ShouldBindJSON(&claim); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	fmt.Println(claim.ItemID, claim.Quantity)

	c.JSON(http.StatusOK, gin.H{
		"message": "Claim made",
	})
}
