package auth

import (
	"farmsville/backend/models"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type Service interface {
	GenerateRandomCode() (string, error)
	SendEmailWithAuthCode(toEmail, code string) error
	GenerateJWT(user models.User) (string, error)
}

type DefaultService struct{}

func NewService() Service {
	return &DefaultService{}
}

func (s *DefaultService) GenerateRandomCode() (string, error) {
	return GenerateRandomCode()
}

func (s *DefaultService) SendEmailWithAuthCode(email, code string) error {
	return SendEmailWithAuthCode(email, code)
}

func (s *DefaultService) GenerateJWT(user models.User) (string, error) {
	return GenerateJWT(user)
}

func GenerateRandomCode() (string, error) {
	const digits = "0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}

	return string(b), nil
}

func SendEmailWithAuthCode(toEmail, code string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpUsername := os.Getenv("GMAIL_USER")
	smtpPassword := os.Getenv("GMAIL_PASS")

	from := smtpUsername
	to := []string{toEmail}
	subject := "Farmsville Authentication"
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

func GenerateJWT(user models.User) (string, error) {

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"admin":   user.Admin,
		"exp":     time.Now().Add(time.Hour * 24 * 90).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "fallback-secret-key" // Only use in development
	}
	tokenString, err := token.SignedString([]byte(jwtSecret))
	return tokenString, err
}
