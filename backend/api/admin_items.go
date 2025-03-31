package api

import (
	"farmsville/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateItem(c *gin.Context) {
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Model(&models.Item{}).Where("id = ?", item.ID).Updates(map[string]interface{}{
		"name":          item.Name,
		"description":   item.Description,
		"quantity":      item.Quantity,
		"remaining_qty": item.RemainingQty,
		"active":        item.Active,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item updated successfully"})
}

func (h *Handler) RemoveItem(c *gin.Context) {
	var item models.ItemRequest
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Model(&models.Item{}).Where("id = ?", item.ID).Update("active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate item"})
		return
	}

	if err := h.db.Model(&models.ClaimedItem{}).Where("item_id = ?", item.ID).Update("active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate claimed items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item and associated claims deactivated successfully"})
}
