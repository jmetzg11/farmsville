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
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}
	var item models.Item
	if err := c.ShouldBind(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingItem models.Item
	if err := h.db.First(&existingItem, item.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	photoPath, err := h.handlePhotoUpload(c, existingItem.PhotoPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"name":          item.Name,
		"description":   item.Description,
		"quantity":      item.Quantity,
		"remaining_qty": item.RemainingQty,
		"active":        true,
	}

	if photoPath != "" {
		updates["photo_path"] = photoPath
	}

	if err := h.db.Model(&models.Item{}).Where("id = ?", item.ID).Updates(updates).Error; err != nil {
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

	photoPath, err := h.handlePhotoUpload(c, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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

func (h *Handler) handlePhotoUpload(c *gin.Context, existingPhotoPath string) (string, error) {
	file, header, err := c.Request.FormFile("photo")
	if err != nil || file == nil {
		return "", nil
	}
	defer file.Close()

	buff := make([]byte, 512)
	if _, err = file.Read(buff); err != nil {
		return "", fmt.Errorf("failed to read photo file: %w", err)
	}

	if _, err = file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("failed to reset file: %w", err)
	}

	fileType := http.DetectContentType(buff)
	if !strings.HasPrefix(fileType, "image/") {
		return "", fmt.Errorf("file is not an image")
	}

	now := time.Now()
	yearMonthDir := fmt.Sprintf("%d%02d", now.Year(), now.Month())
	dirPath := fmt.Sprintf("data/photos/%s", yearMonthDir)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	filename := fmt.Sprintf("%d_%s", now.Unix(), header.Filename)
	filename = strings.ReplaceAll(filename, " ", "_")
	fullPath := fmt.Sprintf("%s/%s", dirPath, filename)
	photoPath := fmt.Sprintf("/%s/%s", yearMonthDir, filename)

	out, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create photo file: %w", err)
	}
	defer out.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
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

	if err = jpeg.Encode(out, resizedImg, &jpeg.Options{Quality: 85}); err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	if existingPhotoPath != "" && existingPhotoPath != photoPath {
		oldPhotoFullPath := fmt.Sprintf("data/photos%s", existingPhotoPath)
		if err := os.Remove(oldPhotoFullPath); err != nil {
			fmt.Printf("Failed to delete old photo: %s. Error: %v\n", oldPhotoFullPath, err)
		}
	}

	return photoPath, nil
}
