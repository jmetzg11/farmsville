package auth

import (
	"farmsville/backend/models"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestDefaultService(t *testing.T) {
	_ = godotenv.Load("../../.env")

	service := NewService()

	t.Run("GenerateRandomCode", func(t *testing.T) {
		code, err := service.GenerateRandomCode()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(code) != 6 {
			t.Fatalf("Expected code to be 6 digits long, got %d", len(code))
		}

		for _, c := range code {
			if c < '0' || c > '9' {
				t.Fatalf("Expected code to contain only digits, found %c", c)
			}
		}

		code2, _ := service.GenerateRandomCode()
		if code == code2 {
			t.Fatalf("Generated the same code twice in a row, which is highly unlikely")
		}
	})

	t.Run("SendEmailWithAuthCode", func(t *testing.T) {
		if os.Getenv("GMAIL_USER") == "" || os.Getenv("GMAIL_PASS") == "" {
			t.Skip("Skipping email test: Environment variables are not set")
		}

		err := service.SendEmailWithAuthCode(os.Getenv("GMAIL_USER"), "123456")
		if err != nil {
			t.Fatalf("Failed to send email: %v", err)
		}
	})

	t.Run("GenerateJWT", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret-key")
		defer os.Unsetenv("JWT_SECRET")

		user := models.User{
			ID:        123,
			Email:     "test@example.com",
			Admin:     true,
			CreatedAt: time.Now(),
		}

		token, err := service.GenerateJWT(user)
		if err != nil {
			t.Fatalf("Failed to generate JWT: %v", err)
		}

		if token == "" {
			t.Fatal("Generated token is empty")
		}

		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			t.Fatalf("Expected JWT format with 3 parts, got %d parts", len(parts))
		}
	})
}
