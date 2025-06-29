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
	"gorm.io/gorm"
)

func (h *Handler) PostBlog(c *gin.Context) {
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	form := c.Request.MultipartForm
	title := form.Value["title"][0]

	blocks, err := parseBlockForm(form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse blocks"})
		return
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

func parseBlockForm(form *multipart.Form) ([]models.Block, error) {
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
					return nil, err
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

	return blocks, nil
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

func (h *Handler) GetBlogs(c *gin.Context) {
	var blogs []models.Blog
	if err := h.db.Preload("Blocks", func(db *gorm.DB) *gorm.DB {
		return db.Order("`order` ASC")
	}).Order("created_at desc").Find(&blogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get blogs"})
		return
	}

	var content []gin.H
	for _, blog := range blogs {
		content = append(content, gin.H{
			"id":     blog.ID,
			"title":  blog.Title,
			"blocks": blog.Blocks,
		})
	}

	response := gin.H{
		"blogs": content,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetBlogTitles(c *gin.Context) {
	var blogs []models.BlogTitles

	if err := h.db.Model(&models.Blog{}).
		Select("id, title, created_at").
		Order("created_at desc").
		Find(&blogs).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch blog titles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"titles": blogs})
}

func (h *Handler) GetBlogById(c *gin.Context) {
	var blog models.Blog
	id := c.Param("id")

	if err := h.db.Preload("Blocks", func(db *gorm.DB) *gorm.DB {
		return db.Order("`order` ASC")
	}).First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blog": blog})
}

func (h *Handler) EditBlog(c *gin.Context) {
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	form := c.Request.MultipartForm
	title := form.Value["title"][0]
	id := form.Value["id"][0]

	var blog models.Blog
	if err := h.db.Preload("Blocks").First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	blocks, err := parseBlockForEdit(form, blog.Blocks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse blocks"})
		return
	}

	if err := h.db.Where("blog_id = ?", blog.ID).Delete(&models.Block{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blocks"})
		return
	}

	blog.Blocks = blocks
	if blog.Title != title {
		blog.Title = title
	}

	if err := h.db.Save(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog updated successfully"})
}

func parseBlockForEdit(form *multipart.Form, existingBlocks []models.Block) ([]models.Block, error) {
	var blocks []models.Block

	existingImages := make(map[string]bool)
	for _, block := range existingBlocks {
		if block.Type == "image" && block.Media != "" {
			existingImages[block.Media] = true
		}
	}

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
			isNewKey := fmt.Sprintf("content[%d][isNew]", index)
			isNew := form.Value[isNewKey][0] == "true"

			if isNew {
				fileKey := fmt.Sprintf("content[%d][file]", index)
				if file, exists := form.File[fileKey]; exists && len(file) > 0 {
					filePath, err := handleBlogPhotoUpload(file[0])
					if err != nil {
						return nil, err
					}
					block.Media = filePath
				}
			} else {
				mediaKey := fmt.Sprintf("content[%d][media]", index)
				if media, exists := form.Value[mediaKey]; exists {
					block.Media = media[0]
				}
			}

			if block.Media != "" && !isNew {
				delete(existingImages, block.Media)
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

	for imagePath := range existingImages {
		fullPath := fmt.Sprintf("data/photos/%s", imagePath)
		os.Remove(fullPath)
	}

	return blocks, nil
}
