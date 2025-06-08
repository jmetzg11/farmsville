package api

import (
	"bytes"
	"encoding/json"
	"farmsville/backend/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPostMessage(t *testing.T) {
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

	messageRequest := models.PostMessageRequest{
		Title:   "Test Message",
		Message: "This is a test message",
	}

	requestBody, _ := json.Marshal(messageRequest)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post-message", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	})

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Errorf("Expected success to be true")
	}

	if message, ok := response["message"].(string); !ok || message != "Message posted" {
		t.Errorf("Expected message to be 'Message posted', got '%v'", message)
	}

	var createdMessage models.Message
	if err := db.Where("title = ?", messageRequest.Title).First(&createdMessage).Error; err != nil {
		t.Errorf("Failed to find created message in database: %v", err)
	}

	if createdMessage.Title != messageRequest.Title {
		t.Errorf("Expected title %s, got %s", messageRequest.Title, createdMessage.Title)
	}
	if createdMessage.Message != messageRequest.Message {
		t.Errorf("Expected message %s, got %s", messageRequest.Message, createdMessage.Message)
	}
}

func TestGetMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	handler := NewHandler(db)
	router := setUpTestRouter(handler)

	testMessages := []models.Message{
		{
			Title:   "First Message",
			Message: "This is the first message",
		},
		{
			Title:   "Second Message",
			Message: "This is the second message",
		},
	}

	for _, msg := range testMessages {
		if err := db.Create(&msg).Error; err != nil {
			t.Fatalf("Failed to create test message: %v", err)
		}
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	messages, ok := response["messages"].([]interface{})
	if !ok {
		t.Fatalf("Expected messages to be an array")
	}

	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
	}
}

func TestDeleteMessage(t *testing.T) {
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

	testMessage := models.Message{
		Title:   "Test Message",
		Message: "This message will be deleted",
	}

	if err := db.Create(&testMessage).Error; err != nil {
		t.Fatalf("Failed to create test message: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/messages/"+strconv.Itoa(int(testMessage.ID)), nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	})

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Errorf("Expected success to be true")
	}

	var deletedMessage models.Message
	result := db.First(&deletedMessage, testMessage.ID)
	if result.Error == nil {
		t.Errorf("Expected message to be deleted, but it still exists in the database")
	}
}

func TestSendEmail(t *testing.T) {
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

	emailRequest := models.SendEmailRequest{
		Title:   "Test Email",
		Message: "This is a test email",
		Emails:  []string{"test@example.com"},
	}

	requestBody, _ := json.Marshal(emailRequest)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/send-email", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	})

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Errorf("Expected success to be true")
	}

	if message, ok := response["message"].(string); !ok || message != "Text message sent" {
		t.Errorf("Expected message to be 'Text message sent', got '%v'", message)
	}
}
