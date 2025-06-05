package api

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
)

type SendEmailRequest struct {
	Emails  []string `json:"emails" binding:"required"`
	Title   string   `json:"title" binding:"required"`
	Message string   `json:"message" binding:"required"`
}

// after adding user the fields don't clear, on the frontend

func (h *Handler) SendEmail(c *gin.Context) {
	var request SendEmailRequest
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
