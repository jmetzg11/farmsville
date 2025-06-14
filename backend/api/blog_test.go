package api

import (
	"bytes"
	"encoding/json"
	"farmsville/backend/models"
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
	t.Logf("Generated token: %s", token)

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
