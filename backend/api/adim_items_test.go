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
		Description:  "Item for removal test",
		Quantity:     100,
		RemainingQty: 50,
		Active:       true,
	}
	db.Create(&item)

	claimedItem := models.ClaimedItem{
		ItemID:   item.ID,
		UserID:   testUser.ID,
		Quantity: 50,
		Active:   true,
	}
	db.Create(&claimedItem)

	requestBody, _ := json.Marshal(models.ItemRequest{
		ID: int(item.ID),
	})

	tokenString, err := getTestUserToken(testUser)
	if err != nil {
		t.Fatalf("Failed to get test user token: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items/remove", bytes.NewBuffer(requestBody))
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
	if updatedItem.Active {
		t.Fatalf("Item was not deactivated in the database")
	}

	var updatedClaimedItem models.ClaimedItem
	if err := db.First(&updatedClaimedItem, claimedItem.ID).Error; err != nil {
		t.Fatalf("Failed to fetch updated claimed item: %v", err)
	}
	if updatedClaimedItem.Active {
		t.Fatalf("Claimed item was not deactivated in the database")
	}
}

func TestCreateItem(t *testing.T) {
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

	newItem := models.CreateItemRequest{
		Name:        "Test Item",
		Description: "Item for creation test",
		Quantity:    100,
	}
	requestBody, _ := json.Marshal(newItem)

	tokenString, err := getTestUserToken(testUser)
	if err != nil {
		t.Fatalf("Failed to get test user token: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items/create", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	})

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var createdItem models.Item
	if err := db.First(&createdItem, "name = ?", newItem.Name).Error; err != nil {
		t.Fatalf("Failed to fetch created item: %v", err)
	}

	if createdItem.Name != newItem.Name ||
		createdItem.Description != newItem.Description ||
		createdItem.Quantity != newItem.Quantity {
		t.Fatalf("Item was not created correctly in the database")
	}
}

func TestRemoveClaimedItem(t *testing.T) {
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

	testItem := models.Item{
		Name:         "Test Item",
		Description:  "Test Description",
		Quantity:     10,
		RemainingQty: 5,
		Active:       true,
	}
	if err := db.Create(&testItem).Error; err != nil {
		t.Fatalf("Failed to create test item: %v", err)
	}

	testClaimedItem := models.ClaimedItem{
		ItemID:   testItem.ID,
		UserID:   testUser.ID,
		Quantity: 3,
		Active:   true,
	}
	if err := db.Create(&testClaimedItem).Error; err != nil {
		t.Fatalf("Failed to create test claimed item: %v", err)
	}

	claimedItemRequest := models.ItemRequest{
		ID: int(testClaimedItem.ID),
	}
	requestBody, _ := json.Marshal(claimedItemRequest)

	tokenString, err := getTestUserToken(testUser)
	if err != nil {
		t.Fatalf("Failed to get test user token: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/claimed-item/remove", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	})

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var updatedItem models.Item
	if err := db.First(&updatedItem, testItem.ID).Error; err != nil {
		t.Fatalf("Failed to fetch updated item: %v", err)
	}

	expectedQuantity := testItem.RemainingQty + testClaimedItem.Quantity
	if updatedItem.RemainingQty != expectedQuantity {
		t.Fatalf("Expected item quantity to be %d, got %d", expectedQuantity, updatedItem.RemainingQty)
	}

	var deletedItem models.ClaimedItem
	result := db.First(&deletedItem, testClaimedItem.ID)
	if result.Error == nil {
		t.Fatalf("Expected claimed item to be deleted, but it still exists")
	}
}
