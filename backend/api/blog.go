package api

import (
	"farmsville/backend/models"
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

func (h *Handler) PostBlog(c *gin.Context) {
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	form := c.Request.MultipartForm
	title := form.Value["title"][0]

	var blocks []models.Block
	index := 0

	for {
		typeKey := fmt.Sprintf("content[%d][type]", index)
		if _, exists := form.Value[typeKey]; !exists {
			break
		}

		blockType := form.Value[typeKey][0]
		block := models.Block{
			Type:  blockType,
			Order: index,
		}

		if blockType == "image" {
			fileKey := fmt.Sprintf("content[%d][file]", index)
			if file, exists := form.File[fileKey]; exists && len(file) > 0 {
				filePath, err := handleBlogPhotoUpload(file[0])
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
					return
				}
				block.Media = filePath
			}
		} else {
			mediaKey := fmt.Sprintf("content[%d][media]", index)
			if media, exists := form.Value[mediaKey]; exists {
				block.Media = media[0]
			}
		}

		blocks = append(blocks, block)
		index++
	}

	blog := models.Blog{
		Title:  title,
		Blocks: blocks,
	}

	if err := h.db.Create(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog posted successfully"})
}

func handleBlogPhotoUpload(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	buff := make([]byte, 512)
	if _, err := src.Read(buff); err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	if _, err := src.Seek(0, 0); err != nil {
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

	filename := fmt.Sprintf("%d_%s", now.Unix(), file.Filename)
	filename = strings.ReplaceAll(filename, " ", "_")
	fullPath := fmt.Sprintf("%s/%s", dirPath, filename)
	photoPath := fmt.Sprintf("/%s/%s", yearMonthDir, filename)

	out, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create photo file: %w", err)
	}
	defer out.Close()

	img, _, err := image.Decode(src)
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

	return photoPath, nil
}
