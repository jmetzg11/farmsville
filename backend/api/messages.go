package api

import (
	"farmsville/backend/models"
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostMessage(c *gin.Context) {
	var request models.PostMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := models.Message{
		Title:   request.Title,
		Message: request.Message,
	}

	if err := h.db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Message posted",
	})
}

func (h *Handler) GetMessages(c *gin.Context) {
	var messages []models.Message
	if err := h.db.Order("created_at DESC").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

func (h *Handler) DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&models.Message{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) SendEmail(c *gin.Context) {
	var request models.SendEmailRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpUsername := os.Getenv("GMAIL_USER")
	smtpPassword := os.Getenv("GMAIL_PASS")

	from := smtpUsername
	subject := request.Title
	body := request.Message

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	for _, email := range request.Emails {
		to := []string{email}
		message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s",
			email, subject, body)
		err := smtp.SendMail(
			fmt.Sprintf("%s:%d", smtpHost, smtpPort),
			auth,
			from,
			to,
			[]byte(message),
		)

		if err != nil {
			fmt.Println(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Text message sent",
	})
}
