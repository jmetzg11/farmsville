package api

import (
	"farmsville/backend/models"
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Create a unique name for this in-memory database instance
	// Each test will get its own isolated database
	dbID := fmt.Sprintf("file:memdb%d?mode=memory&cache=shared", time.Now().UnixNano())

	db, err := gorm.Open(sqlite.Open(dbID), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		sqlDB.Close()
	})

	err = db.AutoMigrate(&models.Item{}, &models.ClaimedItem{}, &models.User{})
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func setUpTestRouter(handler *Handler) *gin.Engine {
	router := gin.New()

	authGroup := router.Group("/")
	adminGroup := router.Group("/")

	authGroup.Use(handler.AuthMiddleware())
	adminGroup.Use(handler.AdminMiddleware())

	// admin
	adminGroup.POST("/items/update", handler.UpdateItem)
	adminGroup.POST("/items/remove", handler.RemoveItem)
	adminGroup.POST("/items/create", handler.CreateItem)
	adminGroup.POST("/claimed-item/remove", handler.RemoveClaimedItem)

	// customers
	router.GET("/items", handler.GetItems)
	authGroup.POST("/claim", handler.MakeClaim)

	// users
	router.POST("/auth/send", handler.SendAuth)
	router.POST("/auth/verify", handler.VerifyAuth)
	router.GET("/auth/me", handler.AuthMe)

	return router
}

func getTestUserToken(testUser models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(testUser.ID),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte("fallback-secret-key"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
