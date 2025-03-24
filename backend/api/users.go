package api

import (
	"farmsville/backend/models"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SendAuth(c *gin.Context) {
	var req models.Email
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	authCode, err := h.generateRandomCode(req.Email, 6)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	err = sendEmailWithAuthCode(req.Email, authCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Authentication email sent",
	})

}

func (h *Handler) generateRandomCode(email string, length int) (string, error) {
	const digits = "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}

	expiresAt := time.Now().Add(5 * time.Minute)
	authCode := models.AuthCode{
		Email:     email,
		Code:      string(b),
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err := h.db.Create(&authCode).Error; err != nil {
		return "", err
	}

	return string(b), nil
}

func sendEmailWithAuthCode(toEmail, code string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpUsername := os.Getenv("GMAIL_USER")
	smtpPassword := os.Getenv("GMAIL_PASS")

	fmt.Print(smtpUsername)
	fmt.Print(smtpPassword)

	from := smtpUsername
	to := []string{toEmail}
	subject := "Authentication code for Farmsville"
	body := fmt.Sprintf("Your authentication code is: %s", code)

	message := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n\r\n"+
		"%s", from, toEmail, subject, body))

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", smtpHost, smtpPort),
		auth,
		from,
		to,
		message,
	)
}
