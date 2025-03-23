package api

import (
	"farmsville/backend/models"
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

	var claimedItems []models.ClaimedItem
	if err := h.db.Where("active = ?", true).Find(&claimedItems).Error; err != nil {
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
