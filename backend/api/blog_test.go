package api

import (
	"bytes"
	"encoding/json"
	"farmsville/backend/models"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostBlog(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	handler := NewHandler(db)
	router := setUpTestRouter(handler)

	testUser := models.User{
		Name:  "Test Admin",
		Email: "admin@example.com",
		Admin: true,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	token, err := getTestUserToken(testUser)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	writer.WriteField("title", "Test Blog Post")
	writer.WriteField("content[0][type]", "text")
	writer.WriteField("content[0][media]", "This is test content")
	writer.WriteField("content[1][type]", "text")
	writer.WriteField("content[1][media]", "Second paragraph")

	writer.Close()

	req, _ := http.NewRequest("POST", "/post-blog", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: token,
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Response body: %s", w.Body.String())
		t.Logf("User ID: %d", testUser.ID)
	}

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Blog posted successfully", response["message"])

	var blog models.Blog
	db.Preload("Blocks").First(&blog)
	assert.Equal(t, "Test Blog Post", blog.Title)
	assert.Len(t, blog.Blocks, 2)
	if len(blog.Blocks) > 0 {
		assert.Equal(t, "text", blog.Blocks[0].Type)
		assert.Equal(t, "This is test content", blog.Blocks[0].Media)
	}
}

func TestGetBlogs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	handler := NewHandler(db)
	router := setUpTestRouter(handler)

	blog := models.Blog{
		Title: "Test Blog",
		Blocks: []models.Block{
			{Type: "text", Media: "First block", Order: 0},
			{Type: "text", Media: "Second block", Order: 1},
		},
	}
	db.Create(&blog)

	req, _ := http.NewRequest("GET", "/blogs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Response body: %s", w.Body.String())
	}

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	blogs := response["blogs"].([]interface{})
	assert.Len(t, blogs, 1)

	firstBlog := blogs[0].(map[string]interface{})
	assert.Equal(t, "Test Blog", firstBlog["title"])

	blocks := firstBlog["blocks"].([]interface{})
	assert.Len(t, blocks, 2)
}

func TestGetBlogTitles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	handler := NewHandler(db)
	router := setUpTestRouter(handler)

	testUser := models.User{
		Name:  "Test Admin",
		Email: "admin@example.com",
		Admin: true,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	blog := models.Blog{
		Title: "Test Blog",
	}
	db.Create(&blog)

	tokenString, err := getTestUserToken(testUser)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	req, _ := http.NewRequest("GET", "/get-blog-titles", nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	titles := response["titles"].([]interface{})
	assert.Len(t, titles, 1)

	firstTitle := titles[0].(map[string]interface{})
	assert.Equal(t, "Test Blog", firstTitle["title"])
}

func TestGetBlogById(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	handler := NewHandler(db)
	router := setUpTestRouter(handler)

	testUser := models.User{
		Name:  "Test Admin",
		Email: "admin@example.com",
		Admin: true,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	token, err := getTestUserToken(testUser)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	blog := models.Blog{
		Title: "Test Blog",
		Blocks: []models.Block{
			{Type: "text", Media: "First block", Order: 0},
			{Type: "text", Media: "Second block", Order: 1},
		},
	}
	db.Create(&blog)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/get-blog/%d", blog.ID), nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: token,
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Response body: %s", w.Body.String())
	}

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	blogData := response["blog"].(map[string]interface{})
	assert.Equal(t, "Test Blog", blogData["title"])

	blocks := blogData["blocks"].([]interface{})
	assert.Len(t, blocks, 2)

	firstBlock := blocks[0].(map[string]interface{})
	assert.Equal(t, "text", firstBlock["type"])
	assert.Equal(t, "First block", firstBlock["media"])
	assert.Equal(t, float64(0), firstBlock["order"])

	secondBlock := blocks[1].(map[string]interface{})
	assert.Equal(t, "text", secondBlock["type"])
	assert.Equal(t, "Second block", secondBlock["media"])
	assert.Equal(t, float64(1), secondBlock["order"])
}

func TestEditBlog(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB(t)
	handler := NewHandler(db)
	router := setUpTestRouter(handler)

	testUser := models.User{
		Name:  "Test Admin",
		Email: "admin@example.com",
		Admin: true,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	token, err := getTestUserToken(testUser)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	blog := models.Blog{
		Title: "Original Title",
		Blocks: []models.Block{
			{Type: "text", Media: "Original content", Order: 0},
			{Type: "text", Media: "Second block", Order: 1},
		},
	}
	db.Create(&blog)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	writer.WriteField("title", "Updated Blog Title")
	writer.WriteField("id", fmt.Sprintf("%d", blog.ID))
	writer.WriteField("content[0][type]", "text")
	writer.WriteField("content[0][media]", "Updated content")
	writer.WriteField("content[0][isNew]", "false")
	writer.WriteField("content[1][type]", "text")
	writer.WriteField("content[1][media]", "New second block")
	writer.WriteField("content[1][isNew]", "false")

	writer.Close()

	req, _ := http.NewRequest("POST", "/edit-blog", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: token,
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Response body: %s", w.Body.String())
	}

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Blog updated successfully", response["message"])

	var updatedBlog models.Blog
	db.Preload("Blocks").First(&updatedBlog, blog.ID)
	assert.Equal(t, "Updated Blog Title", updatedBlog.Title)
	assert.Len(t, updatedBlog.Blocks, 2)
	if len(updatedBlog.Blocks) > 0 {
		assert.Equal(t, "text", updatedBlog.Blocks[0].Type)
		assert.Equal(t, "Updated content", updatedBlog.Blocks[0].Media)
		assert.Equal(t, "New second block", updatedBlog.Blocks[1].Media)
	}
}
