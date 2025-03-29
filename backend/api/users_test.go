package api

import (
	"bytes"
	"encoding/json"
	"farmsville/backend/models"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockAuthService struct{}

func (m *MockAuthService) GenerateRandomCode() (string, error) {
	return "123456", nil
}

func (m *MockAuthService) SendEmailWithAuthCode(email, code string) error {
	return nil
}

func (m *MockAuthService) GenerateJWT(user models.User) (string, error) {
	return "jwt_token_123", nil
}

func setupUsersDB(t *testing.T) *gorm.DB {
	// Create a unique name for this in-memory database instance
	dbID := fmt.Sprintf("file:memdb%d?mode=memory&cache=shared", time.Now().UnixNano())

	db, err := gorm.Open(sqlite.Open(dbID), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		sqlDB.Close()
	})

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func setupUsersRouter(handler *Handler) *gin.Engine {
	router := gin.New()

	router.POST("/auth/send", handler.SendAuth)
	router.POST("/auth/verify", handler.VerifyAuth)
	router.GET("/auth/me", handler.AuthMe)
	return router
}

func TestSendAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup database and handler
	db := setupUsersDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setupUsersRouter(handler)

	// Create request body
	reqBody := map[string]string{
		"email": "test@example.com",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/send", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	router.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Errorf("Expected success to be true, got %v", response["success"])
	}

	if message, ok := response["message"].(string); !ok || message != "Authentication email sent" {
		t.Errorf("Expected message to be 'Authentication email sent', got %v", response["message"])
	}

	// Verify user was created in database
	var user models.User
	result := db.Where("email = ?", "test@example.com").First(&user)
	if result.Error != nil {
		t.Fatalf("Failed to find user: %v", result.Error)
	}

	if user.Code != "123456" {
		t.Errorf("Expected user code to be '123456', got '%s'", user.Code)
	}

	if user.Admin {
		t.Errorf("Expected user to not be admin")
	}
}

func TestAuthMe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupUsersDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setupUsersRouter(handler)

	testUser := models.User{
		Name:  "Test User",
		Email: "test@example.com",
		Admin: false,
	}
	db.Create(&testUser)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(testUser.ID),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte("fallback-secret-key"))
	if err != nil {
		t.Fatalf("Failed to sign test token: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/me", nil)

	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	})

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Errorf("Expected success to be true, got %v", response["success"])
	}

	if message, ok := response["message"].(string); !ok || message != "Authentication successful" {
		t.Errorf("Expected message to be 'Authentication successful', got %v", response["message"])
	}

	userData, ok := response["user"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected user data in response")
	}

	if name, ok := userData["name"].(string); !ok || name != testUser.Name {
		t.Errorf("Expected user name to be '%s', got '%v'", testUser.Name, userData["name"])
	}

	if email, ok := userData["email"].(string); !ok || email != testUser.Email {
		t.Errorf("Expected user email to be '%s', got '%v'", testUser.Email, userData["email"])
	}

	if admin, ok := userData["admin"].(bool); !ok || admin != testUser.Admin {
		t.Errorf("Expected user admin status to be %v, got %v", testUser.Admin, userData["admin"])
	}
}

func TestAuthMeInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupUsersDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setupUsersRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/me", nil)

	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: "invalid-token",
	})

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
