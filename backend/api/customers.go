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
