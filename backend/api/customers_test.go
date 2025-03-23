package api

import (
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
	router := setupTestRouter(handler)

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
		User:     "user1@example.com",
		Quantity: 10,
		Active:   true,
	}

	claimedItem2 := models.ClaimedItem{
		ItemID:   item1.ID,
		User:     "user2@example.com",
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

	// Verify that the claimed item contains the item_name field
	if len(claimedItems) > 0 {
		claimedItem, ok := claimedItems[0].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected claimed item to be a map")
		}

		itemName, exists := claimedItem["item_name"]
		if !exists {
			t.Errorf("Expected claimed item to have 'item_name' field")
		}

		// Check that the item_name matches the expected value
		expectedName := "Test Item 1"
		if itemName != expectedName {
			t.Errorf("Expected item_name to be '%s', got '%v'", expectedName, itemName)
		}
	}

}
