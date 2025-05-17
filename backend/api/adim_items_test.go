package api

import (
	"bytes"
	"encoding/json"
	"farmsville/backend/models"
	"image"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var updatedItemFromDB models.Item
	if err := db.First(&updatedItemFromDB, item.ID).Error; err != nil {
		t.Fatalf("Failed to fetch updated item: %v", err)
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

	tempDir := "data/photos"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Fatalf("Failed to create test photo directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testUser := models.User{
		Name:  "Test User",
		Email: "testuser@example.com",
		Admin: true,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	router := setUpTestRouter(handler)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	testItem := struct {
		title       string
		description string
		quantity    int
	}{
		title:       "Test Item",
		description: "Item for creation test",
		quantity:    100,
	}

	_ = writer.WriteField("title", testItem.title)
	_ = writer.WriteField("description", testItem.description)
	_ = writer.WriteField("quantity", strconv.Itoa(testItem.quantity))

	fileWriter, err := writer.CreateFormFile("photo", "test_image.jpg")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	if err := jpeg.Encode(fileWriter, img, nil); err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}

	writer.Close()

	tokenString, err := getTestUserToken(testUser)
	if err != nil {
		t.Fatalf("Failed to get test user token: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items/create", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	})

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var createdItem models.Item
	if err := db.First(&createdItem, "name = ?", testItem.title).Error; err != nil {
		t.Fatalf("Failed to fetch created item: %v", err)
	}

	if createdItem.Name != testItem.title ||
		createdItem.Description != testItem.description ||
		createdItem.Quantity != testItem.quantity {
		t.Fatalf("Item was not created correctly in the database")
	}

	if _, err := os.Stat(createdItem.PhotoPath); os.IsNotExist(err) {
		t.Fatalf("Photo file does not exist at path: %s", createdItem.PhotoPath)
	}
}

func TestAdminClaimItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	handler := NewHandler(db)

	adminUser := models.User{
		Name:  "Admin User",
		Email: "admin@example.com",
		Admin: true,
	}
	if err := db.Create(&adminUser).Error; err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}

	regularUser := models.User{
		Name:  "Regular User",
		Email: "user@example.com",
		Admin: false,
	}
	if err := db.Create(&regularUser).Error; err != nil {
		t.Fatalf("Failed to create regular user: %v", err)
	}

	router := setUpTestRouter(handler)

	item := models.Item{
		Name:         "Test Item",
		Description:  "Item for admin claiming test",
		Quantity:     100,
		RemainingQty: 100,
		Active:       true,
	}
	db.Create(&item)

	claimRequest := models.AdminClaimItemRequest{
		ItemID: int(item.ID),
		UserID: int(regularUser.ID),
		Amount: 30,
	}
	requestBody, _ := json.Marshal(claimRequest)

	tokenString, err := getTestUserToken(adminUser)
	if err != nil {
		t.Fatalf("Failed to get admin user token: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items/admin-claim", bytes.NewBuffer(requestBody))
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
	if err := db.First(&updatedItem, item.ID).Error; err != nil {
		t.Fatalf("Failed to fetch updated item: %v", err)
	}

	expectedRemainingQty := item.RemainingQty - claimRequest.Amount
	if updatedItem.RemainingQty != expectedRemainingQty {
		t.Fatalf("Expected remaining quantity to be %d, got %d", expectedRemainingQty, updatedItem.RemainingQty)
	}

	var claimedItem models.ClaimedItem
	if err := db.Where("item_id = ? AND user_id = ?", item.ID, regularUser.ID).First(&claimedItem).Error; err != nil {
		t.Fatalf("Failed to fetch claimed item: %v", err)
	}

	if claimedItem.Quantity != claimRequest.Amount || !claimedItem.Active {
		t.Fatalf("Claimed item was not created correctly in the database")
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
