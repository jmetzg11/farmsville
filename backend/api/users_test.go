package api

import (
	"bytes"
	"encoding/json"
	"farmsville/backend/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	handler := NewHandler(db)
	router := setUpTestRouter(handler)

	testUsers := []models.User{
		{
			Name:  "User One",
			Email: "user1@example.com",
			Phone: "123-456-7890",
			Admin: true,
		},
		{
			Name:  "User Two",
			Email: "user2@example.com",
			Phone: "987-654-3210",
			Admin: false,
		},
	}

	for _, user := range testUsers {
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	adminUser := models.User{
		Name:  "Admin User",
		Email: "admin@example.com",
		Admin: true,
	}
	if err := db.Create(&adminUser).Error; err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}

	tokenString, err := getTestUserToken(adminUser)
	if err != nil {
		t.Fatalf("Failed to get test user token: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	})

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var responseUsers []models.UserResponse
	err = json.Unmarshal(w.Body.Bytes(), &responseUsers)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(responseUsers) < 3 {
		t.Fatalf("Expected at least 3 users, got %d", len(responseUsers))
	}

	foundUser1 := false
	foundUser2 := false

	for _, user := range responseUsers {
		if user.Email == "user1@example.com" {
			foundUser1 = true
			if user.Name != "User One" || user.Phone != "123-456-7890" || !user.Admin {
				t.Fatalf("User1 data is incorrect")
			}
		}
		if user.Email == "user2@example.com" {
			foundUser2 = true
			if user.Name != "User Two" || user.Phone != "987-654-3210" || user.Admin {
				t.Fatalf("User2 data is incorrect")
			}
		}
	}

	if !foundUser1 || !foundUser2 {
		t.Fatalf("Not all test users were returned")
	}
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	handler := NewHandler(db)
	router := setUpTestRouter(handler)

	adminUser := models.User{
		Name:  "Admin User",
		Email: "admin@example.com",
		Admin: true,
	}

	if err := db.Create(&adminUser).Error; err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}

	tokenString, err := getTestUserToken(adminUser)
	if err != nil {
		t.Fatalf("Failed to get test user token: %v", err)
	}

	newUser := models.CreateUserRequest{
		Name:  "Test User",
		Email: "testuser@example.com",
		Phone: "555-123-4567",
	}

	requestBody, err := json.Marshal(newUser)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/create", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	})

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Fatalf("Expected success to be true, got %v", response["success"])
	}

	if message, ok := response["message"].(string); !ok || message != "User created" {
		t.Fatalf("Expected message to be 'User created', got %v", response["message"])
	}

	var createdUser models.User
	result := db.Where("email = ?", newUser.Email).First(&createdUser)
	if result.Error != nil {
		t.Fatalf("Failed to find created user in database: %v", result.Error)
	}

	if createdUser.Name != newUser.Name {
		t.Errorf("Expected name %s, got %s", newUser.Name, createdUser.Name)
	}
	if createdUser.Email != newUser.Email {
		t.Errorf("Expected email %s, got %s", newUser.Email, createdUser.Email)
	}
	if createdUser.Phone != newUser.Phone {
		t.Errorf("Expected phone %s, got %s", newUser.Phone, createdUser.Phone)
	}

}
