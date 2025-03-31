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

func TestUpdateItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	handler := NewHandler(db)
	testUser := models.User{
		Name:  "Test User",
		Email: "testuser@example.com",
		Admin: true,
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

	updatedItem := models.Item{
		ID:           item.ID,
		Name:         "Test Item",
		Description:  "Item for claiming",
		Quantity:     100,
		RemainingQty: 100,
		Active:       true,
	}
	requestBody, _ := json.Marshal(updatedItem)

	tokenString, err := getTestUserToken(testUser)
	if err != nil {
		t.Fatalf("Failed to get test user token: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items/update", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	})
	router.ServeHTTP(w, req)

	var updatedItemFromDB models.Item
	if err := db.First(&updatedItemFromDB, item.ID).Error; err != nil {
		t.Fatalf("Failed to fetch updated item: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if updatedItemFromDB.Name != updatedItem.Name ||
		updatedItemFromDB.Description != updatedItem.Description ||
		updatedItemFromDB.Quantity != updatedItem.Quantity ||
		updatedItemFromDB.RemainingQty != updatedItem.RemainingQty ||
		updatedItemFromDB.Active != updatedItem.Active {
		t.Fatalf("Item was not updated correctly in the database")
	}
}

func TestRemoveItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	handler := NewHandler(db)

	// Create test user
	testUser := models.User{
		Name:  "Test User",
		Email: "testuser@example.com",
		Admin: true,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	router := setUpTestRouter(handler)

	// Create test item
	item := models.Item{
		Name:         "Test Item",
		Description:  "Item for removal test",
		Quantity:     100,
		RemainingQty: 50,
		Active:       true,
	}
	db.Create(&item)

	// Create a claimed item associated with the test item
	claimedItem := models.ClaimedItem{
		ItemID:   item.ID,
		UserID:   testUser.ID,
		Quantity: 50,
		Active:   true,
	}
	db.Create(&claimedItem)

	// Create request body
	requestBody, _ := json.Marshal(models.ItemRequest{
		ID: int(item.ID),
	})

	// Get test user token
	tokenString, err := getTestUserToken(testUser)
	if err != nil {
		t.Fatalf("Failed to get test user token: %v", err)
	}

	// Make the request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items/remove", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	})

	router.ServeHTTP(w, req)

	// Check response status
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Verify the item was deactivated in the database
	var updatedItem models.Item
	if err := db.First(&updatedItem, item.ID).Error; err != nil {
		t.Fatalf("Failed to fetch updated item: %v", err)
	}
	if updatedItem.Active {
		t.Fatalf("Item was not deactivated in the database")
	}

	// Verify the claimed item was deactivated in the database
	var updatedClaimedItem models.ClaimedItem
	if err := db.First(&updatedClaimedItem, claimedItem.ID).Error; err != nil {
		t.Fatalf("Failed to fetch updated claimed item: %v", err)
	}
	if updatedClaimedItem.Active {
		t.Fatalf("Claimed item was not deactivated in the database")
	}
}
