package api

import (
	"bytes"
	"encoding/json"
	"farmsville/backend/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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

func TestSendAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup database and handler
	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	// Create request body
	reqBody := map[string]string{
		"email": "test@example.com",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth", bytes.NewBuffer(jsonBody))
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

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "fallback-secret-key"
	}
	tokenString, err := token.SignedString([]byte(jwtSecret))
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

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

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

func TestLoginWithPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	testUser := models.User{
		Name:  "Test User",
		Email: "test@example.com",
		Admin: false,
	}
	testUser.SetPassword("testpassword123")
	db.Create(&testUser)

	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "testpassword123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

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

	if message, ok := response["message"].(string); !ok || message != "Login successful" {
		t.Errorf("Expected message to be 'Login successful', got %v", response["message"])
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

	if isAuth, ok := userData["isAuthenticated"].(bool); !ok || !isAuth {
		t.Errorf("Expected isAuthenticated to be true, got %v", userData["isAuthenticated"])
	}

	cookies := w.Result().Cookies()
	var authCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			authCookie = cookie
			break
		}
	}
	if authCookie == nil {
		t.Error("Expected auth_token cookie to be set")
	} else if authCookie.Value != "jwt_token_123" {
		t.Errorf("Expected cookie value to be 'jwt_token_123', got '%s'", authCookie.Value)
	}
}

func TestLoginWithPasswordInvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	testUser := models.User{
		Name:  "Test User",
		Email: "test@example.com",
		Admin: false,
	}
	testUser.SetPassword("correctpassword")
	db.Create(&testUser)

	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "workpassword",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestLoginWithPasswordUserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	reqBody := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}
func TestCreateAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	// Test successful account creation
	reqBody := map[string]string{
		"name":     "John Doe",
		"email":    "john@example.com",
		"phone":    "+1234567890",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

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

	if message, ok := response["message"].(string); !ok || message != "Account created successfully" {
		t.Errorf("Expected message to be 'Account created successfully', got %v", response["message"])
	}

	userData, ok := response["user"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected user data in response")
	}

	if name, ok := userData["name"].(string); !ok || name != "John Doe" {
		t.Errorf("Expected user name to be 'John Doe', got '%v'", userData["name"])
	}

	if email, ok := userData["email"].(string); !ok || email != "john@example.com" {
		t.Errorf("Expected user email to be 'john@example.com', got '%v'", userData["email"])
	}

	if admin, ok := userData["admin"].(bool); !ok || admin {
		t.Errorf("Expected user admin status to be false, got %v", userData["admin"])
	}

	if isAuth, ok := userData["isAuthenticated"].(bool); !ok || !isAuth {
		t.Errorf("Expected isAuthenticated to be true, got %v", userData["isAuthenticated"])
	}

	var user models.User
	result := db.Where("email = ?", "john@example.com").First(&user)
	if result.Error != nil {
		t.Fatalf("Failed to find created user: %v", result.Error)
	}

	if user.Name != "John Doe" {
		t.Errorf("Expected user name to be 'John Doe', got '%s'", user.Name)
	}

	if user.Phone != "+1234567890" {
		t.Errorf("Expected user phone to be '+1234567890', got '%s'", user.Phone)
	}

	if user.Admin {
		t.Errorf("Expected user to not be admin")
	}

	cookies := w.Result().Cookies()
	var authCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			authCookie = cookie
			break
		}
	}
	if authCookie == nil {
		t.Error("Expected auth_token cookie to be set")
	} else if authCookie.Value != "jwt_token_123" {
		t.Errorf("Expected cookie value to be 'jwt_token_123', got '%s'", authCookie.Value)
	}
}

func TestCreateAccountAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	existingUser := models.User{
		Name:  "Existing User",
		Email: "existing@example.com",
		Admin: false,
	}
	db.Create(&existingUser)

	reqBody := map[string]string{
		"name":     "New User",
		"email":    "existing@example.com",
		"phone":    "+1234567890",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status code %d, got %d", http.StatusConflict, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || success {
		t.Errorf("Expected success to be false, got %v", response["success"])
	}

	if message, ok := response["message"].(string); !ok || message != "Account already exists" {
		t.Errorf("Expected message to be 'Account already exists', got %v", response["message"])
	}
}

func TestCreateAccountInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/create", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSendCodeToResetPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	testUser := models.User{
		Name:  "Test User",
		Email: "test@example.com",
		Admin: false,
	}
	db.Create(&testUser)

	reqBody := map[string]string{
		"email": "test@example.com",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/code-to-reset-password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

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

	if message, ok := response["message"].(string); !ok || message != "Password reset email sent" {
		t.Errorf("Expected message to be 'Password reset email sent', got %v", response["message"])
	}

	var updatedUser models.User
	result := db.Where("email = ?", "test@example.com").First(&updatedUser)
	if result.Error != nil {
		t.Fatalf("Failed to find updated user: %v", result.Error)
	}

	if updatedUser.Code != "123456" {
		t.Errorf("Expected user code to be '123456', got '%s'", updatedUser.Code)
	}

	if updatedUser.ExpiresAt.IsZero() {
		t.Error("Expected ExpiresAt to be set")
	}

	expectedExpiry := time.Now().Add(15 * time.Minute)
	timeDiff := updatedUser.ExpiresAt.Sub(expectedExpiry).Abs()
	if timeDiff > time.Minute {
		t.Errorf("Expected expiration to be around 15 minutes from now, got difference of %v", timeDiff)
	}
}

func TestSendCodeToResetPasswordUserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	reqBody := map[string]string{
		"email": "nonexistent@example.com",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/code-to-reset-password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || success {
		t.Errorf("Expected success to be false, got %v", response["success"])
	}

	if message, ok := response["message"].(string); !ok || message != "Account does not exist. Please create an account." {
		t.Errorf("Expected message to be 'Account does not exist. Please create an account.', got %v", response["message"])
	}
}

func TestSendCodeToResetPasswordInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/reset-password", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestResetPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	testUser := models.User{
		Name:      "Test User",
		Email:     "test@example.com",
		Admin:     false,
		Code:      "123456",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	db.Create(&testUser)

	reqBody := map[string]string{
		"email":    "test@example.com",
		"code":     "123456",
		"password": "newpassword123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/reset-password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

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

	if message, ok := response["message"].(string); !ok || message != "Password reset successful" {
		t.Errorf("Expected message to be 'Password reset successful', got %v", response["message"])
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

	if isAuth, ok := userData["isAuthenticated"].(bool); !ok || !isAuth {
		t.Errorf("Expected isAuthenticated to be true, got %v", userData["isAuthenticated"])
	}

	cookies := w.Result().Cookies()
	var authCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			authCookie = cookie
			break
		}
	}
	if authCookie == nil {
		t.Error("Expected auth_token cookie to be set")
	} else if authCookie.Value != "jwt_token_123" {
		t.Errorf("Expected cookie value to be 'jwt_token_123', got '%s'", authCookie.Value)
	}

	var updatedUser models.User
	result := db.Where("email = ?", "test@example.com").First(&updatedUser)
	if result.Error != nil {
		t.Fatalf("Failed to find updated user: %v", result.Error)
	}

	if !updatedUser.CheckPassword("newpassword123") {
		t.Error("Expected new password to be set correctly")
	}
}

func TestResetPasswordInvalidCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	testUser := models.User{
		Name:      "Test User",
		Email:     "test@example.com",
		Admin:     false,
		Code:      "123456",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	db.Create(&testUser)

	reqBody := map[string]string{
		"email":    "test@example.com",
		"code":     "wrongcode",
		"password": "newpassword123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/reset-password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || success {
		t.Errorf("Expected success to be false, got %v", response["success"])
	}

	if message, ok := response["message"].(string); !ok || message != "Invalid or expired reset code" {
		t.Errorf("Expected message to be 'Invalid or expired reset code', got %v", response["message"])
	}
}

func TestResetPasswordExpiredCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	// Create user with expired reset code
	testUser := models.User{
		Name:      "Test User",
		Email:     "test@example.com",
		Admin:     false,
		Code:      "123456",
		ExpiresAt: time.Now().Add(-5 * time.Minute),
	}
	db.Create(&testUser)

	reqBody := map[string]string{
		"email":    "test@example.com",
		"code":     "123456",
		"password": "newpassword123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/reset-password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || success {
		t.Errorf("Expected success to be false, got %v", response["success"])
	}

	if message, ok := response["message"].(string); !ok || message != "Invalid or expired reset code" {
		t.Errorf("Expected message to be 'Invalid or expired reset code', got %v", response["message"])
	}
}

func TestResetPasswordUserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	reqBody := map[string]string{
		"email":    "nonexistent@example.com",
		"code":     "123456",
		"password": "newpassword123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/reset-password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || success {
		t.Errorf("Expected success to be false, got %v", response["success"])
	}

	if message, ok := response["message"].(string); !ok || message != "User not found" {
		t.Errorf("Expected message to be 'User not found', got %v", response["message"])
	}
}

func TestResetPasswordInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/reset-password", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestLogout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/logout", nil)

	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: "existing_token_123",
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

	if message, ok := response["message"].(string); !ok || message != "Logged out successfully" {
		t.Errorf("Expected message to be 'Logged out successfully', got %v", response["message"])
	}

	cookies := w.Result().Cookies()
	var authCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			authCookie = cookie
			break
		}
	}

	if authCookie == nil {
		t.Error("Expected auth_token cookie to be set in response")
	} else {
		if authCookie.Value != "" {
			t.Errorf("Expected cookie value to be empty, got '%s'", authCookie.Value)
		}
		if authCookie.MaxAge != -1 {
			t.Errorf("Expected cookie MaxAge to be -1, got %d", authCookie.MaxAge)
		}
	}
}

func TestLogoutWithoutExistingCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	mockAuth := &MockAuthService{}
	handler := NewHandler(db, mockAuth)
	router := setUpTestRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/logout", nil)

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

	if message, ok := response["message"].(string); !ok || message != "Logged out successfully" {
		t.Errorf("Expected message to be 'Logged out successfully', got %v", response["message"])
	}

	cookies := w.Result().Cookies()
	var authCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			authCookie = cookie
			break
		}
	}

	if authCookie == nil {
		t.Error("Expected auth_token cookie to be set in response")
	} else {
		if authCookie.Value != "" {
			t.Errorf("Expected cookie value to be empty, got '%s'", authCookie.Value)
		}
		if authCookie.MaxAge != -1 {
			t.Errorf("Expected cookie MaxAge to be -1, got %d", authCookie.MaxAge)
		}
	}
}
