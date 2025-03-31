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

func TestGetItems(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	handler := NewHandler(db)

	router := setUpTestRouter(handler)

	user1 := models.User{
		Name:  "User One",
		Email: "user1@example.com",
	}

	user2 := models.User{
		Name:  "User Two",
		Email: "user2@example.com",
	}

	db.Create(&user1)
	db.Create(&user2)

	item1 := models.Item{
		Name:         "Test Item 1",
		Description:  "First test item",
		Quantity:     100,
		RemainingQty: 80,
		Active:       true,
	}

	item2 := models.Item{
		Name:         "Test Item 2",
		Description:  "Second test item",
		Quantity:     50,
		RemainingQty: 45,
		Active:       true,
	}

	item3 := models.Item{
		Name:         "Test Item 3",
		Description:  "third test item",
		Quantity:     50,
		RemainingQty: 45,
		Active:       false,
	}

	db.Create(&item1)
	db.Create(&item2)
	db.Create(&item3)

	claimedItem1 := models.ClaimedItem{
		ItemID:   item1.ID,
		UserID:   user1.ID,
		Quantity: 10,
		Active:   true,
	}

	claimedItem2 := models.ClaimedItem{
		ItemID:   item1.ID,
		UserID:   user2.ID,
		Quantity: 10,
		Active:   false,
	}

	db.Create(&claimedItem1)
	db.Create(&claimedItem2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	items, ok := response["items"].([]interface{})
	if !ok {
		t.Fatalf("Expected items to be an array")
	}

	claimedItems, ok := response["claimedItems"].([]interface{})
	if !ok {
		t.Fatalf("Expected claimedItems to be an array")
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items (only active ones), got %d", len(items))
	}

	if len(claimedItems) != 1 {
		t.Errorf("Expected 1 claimed items (only active ones), got %d", len(claimedItems))
	}

	// Verify that the claimed item contains all the required fields
	if len(claimedItems) > 0 {
		claimedItem, ok := claimedItems[0].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected claimed item to be a map")
		}

		// Check for item_name field
		itemName, exists := claimedItem["item_name"]
		if !exists {
			t.Errorf("Expected claimed item to have 'item_name' field")
		} else if itemName != "Test Item 1" {
			t.Errorf("Expected item_name to be 'Test Item 1', got '%v'", itemName)
		}

		// Check for user_name field
		userName, exists := claimedItem["user_name"]
		if !exists {
			t.Errorf("Expected claimed item to have 'user_name' field")
		} else if userName != "User One" {
			t.Errorf("Expected user_name to be 'User One', got '%v'", userName)
		}

		// Check for user_email field
		userEmail, exists := claimedItem["user_email"]
		if !exists {
			t.Errorf("Expected claimed item to have 'user_email' field")
		} else if userEmail != "user1@example.com" {
			t.Errorf("Expected user_email to be 'user1@example.com', got '%v'", userEmail)
		}
	}

}

func TestMakeClaim(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	handler := NewHandler(db)

	// Create a test user in the database
	testUser := models.User{
		Name:  "Test User",
		Email: "testuser@example.com",
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	router := setUpTestRouter(handler)

	item := models.Item{
		Name:         "Test Item",
		Description:  "Item for claiming",
		Quantity:     100,
		RemainingQty: 100,
		Active:       true,
	}
	db.Create(&item)

	initialQty := item.RemainingQty

	claimQty := 10
	requestBody, _ := json.Marshal(models.ClaimRequest{
		ItemID:   int(item.ID),
		Quantity: claimQty,
	})

	tokenString, err := getTestUserToken(testUser)
	if err != nil {
		t.Fatalf("Failed to get test user token: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/claim", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	})

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var updatedItem models.Item
	if err := db.First(&updatedItem, item.ID).Error; err != nil {
		t.Fatalf("Failed to fetch updated item: %v", err)
	}
	expectedQty := initialQty - claimQty
	if updatedItem.RemainingQty != expectedQty {
		t.Errorf("Expected remaining quantity to be %d, got %d", expectedQty, updatedItem.RemainingQty)
	}

	var claimedItems []models.ClaimedItem
	if err := db.Where("item_id = ?", item.ID).Find(&claimedItems).Error; err != nil {
		t.Fatalf("Failed to fetch claimed items: %v", err)
	}
	if len(claimedItems) != 1 {
		t.Errorf("Expected 1 claimed item, got %d", len(claimedItems))
	}
}
