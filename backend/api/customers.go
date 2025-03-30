package api

import (
	"farmsville/backend/models"
	"net/http"
	"time"

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

	var claimedItems []models.ClaimedItemWithUserName
	if err := h.db.Raw(`
        SELECT
            claimed_items.*,
            items.name as item_name,
            users.name as user_name,
            users.email as user_email
        FROM
            claimed_items
        JOIN
            items ON claimed_items.item_id = items.id
        JOIN
            users ON claimed_items.user_id = users.id
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
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}
	currentUser, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user information",
		})
		return
	}

	var claimRequest models.ClaimRequest
	if err := c.ShouldBindJSON(&claimRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	var item models.Item
	if err := h.db.First(&item, claimRequest.ItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Item not found",
		})
		return
	}

	if item.RemainingQty < claimRequest.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not enough items in stock",
		})
		return
	}

	item.RemainingQty -= claimRequest.Quantity
	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update item quantity",
		})
		return
	}

	claimedItem := models.ClaimedItem{
		ItemID:    uint(claimRequest.ItemID),
		UserID:    currentUser.ID,
		Quantity:  claimRequest.Quantity,
		CreatedAt: time.Now(),
		Active:    true,
	}

	if err := h.db.Create(&claimedItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create claimed item",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Claim made",
		"claim_id": claimedItem.ID,
	})
}
