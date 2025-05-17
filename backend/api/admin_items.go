package api

import (
	"farmsville/backend/models"
	"fmt"
	"image"
	"net/http"
	"os"
	"strings"
	"time"

	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
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

func (h *Handler) CreateItem(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	var newItem models.CreateItemRequest
	if err := c.ShouldBind(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photoPath := ""

	file, header, err := c.Request.FormFile("photo")
	if err == nil && file != nil {
		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read photo file"})
			return
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to reset file"})
			return
		}

		fileType := http.DetectContentType(buff)
		if !strings.HasPrefix(fileType, "image/") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is not an image"})
			return
		}

		now := time.Now()
		dirPath := fmt.Sprintf("data/photos/%d%02d", now.Year(), now.Month())
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
			return
		}

		filename := fmt.Sprintf("%d_%s", now.Unix(), header.Filename)
		filename = strings.ReplaceAll(filename, " ", "_")
		photoPath = fmt.Sprintf("%s/%s", dirPath, filename)

		out, err := os.Create(photoPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create photo file"})
			return
		}
		defer out.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode image"})
			return
		}

		var resizedImg image.Image
		bounds := img.Bounds()
		maxWidth := 1200
		if bounds.Dx() > maxWidth {
			newHeight := int(float64(bounds.Dy()) * float64(maxWidth) / float64(bounds.Dx()))
			resizedImg = resize.Resize(uint(maxWidth), uint(newHeight), img, resize.Lanczos3)
		} else {
			resizedImg = img
		}

		err = jpeg.Encode(out, resizedImg, &jpeg.Options{Quality: 85})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode image"})
			return
		}
	}

	item := models.Item{
		Name:         newItem.Name,
		Description:  newItem.Description,
		PhotoPath:    photoPath,
		Quantity:     newItem.Quantity,
		RemainingQty: newItem.Quantity,
		Active:       true,
	}

	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item created successfully"})
}

func (h *Handler) AdminClaimItem(c *gin.Context) {
	var claimRequest models.AdminClaimItemRequest
	if err := c.ShouldBindJSON(&claimRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var item models.Item
	if err := h.db.First(&item, claimRequest.ItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Item not found",
		})
		return
	}

	if item.RemainingQty < claimRequest.Amount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not enough items in stock",
		})
		return
	}

	item.RemainingQty -= claimRequest.Amount
	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update item quantity",
		})
		return
	}

	claimedItem := models.ClaimedItem{
		ItemID:    uint(claimRequest.ItemID),
		UserID:    uint(claimRequest.UserID),
		Quantity:  claimRequest.Amount,
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

func (h *Handler) RemoveClaimedItem(c *gin.Context) {
	var claimedItemRequest models.ItemRequest
	if err := c.ShouldBindJSON(&claimedItemRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var claimedItem models.ClaimedItem
	if result := h.db.First(&claimedItem, claimedItemRequest.ID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Claimed item not found"})
		return
	}

	var item models.Item
	if result := h.db.First(&item, claimedItem.ItemID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parent item not found"})
		return
	}

	item.RemainingQty += claimedItem.Quantity
	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update parent item"})
		return
	}

	if err := h.db.Delete(&claimedItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete claimed item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Claimed item deactivated successfully"})
}
