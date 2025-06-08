package api

import (
	"farmsville/backend/models"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Create a unique name for this in-memory database instance
	// Each test will get its own isolated database
	dbID := fmt.Sprintf("file:memdb%d?mode=memory&cache=shared", time.Now().UnixNano())

	db, err := gorm.Open(sqlite.Open(dbID), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
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

	err = db.AutoMigrate(&models.Item{}, &models.ClaimedItem{}, &models.User{}, &models.Message{})
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func setUpTestRouter(handler *Handler) *gin.Engine {
	// Try to load env, but don't fail tests if not found
	_ = godotenv.Load("../../.env")

	// Set defaults for required environment variables in tests
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test-jwt-secret")
	}

	router := gin.New()

	authGroup := router.Group("/")
	adminGroup := router.Group("/")

	authGroup.Use(handler.AuthMiddleware())
	adminGroup.Use(handler.AdminMiddleware())

	// admin
	adminGroup.GET("/users", handler.GetUsers)
	adminGroup.POST("/users/update", handler.UpdateUser)
	adminGroup.POST("/users/remove", handler.RemoveUser)
	adminGroup.POST("/users/create", handler.CreateUser)
	adminGroup.POST("/items/update", handler.UpdateItem)
	adminGroup.POST("/items/remove", handler.RemoveItem)
	adminGroup.POST("/items/create", handler.CreateItem)
	adminGroup.POST("/items/admin-claim", handler.AdminClaimItem)
	adminGroup.POST("/claimed-item/remove", handler.RemoveClaimedItem)

	// messages
	router.GET("/messages", handler.GetMessages)
	adminGroup.POST("/post-message", handler.PostMessage)
	adminGroup.DELETE("/messages/:id", handler.DeleteMessage)
	adminGroup.POST("/send-email", handler.SendEmail)

	// customers
	router.GET("/items", handler.GetItems)
	authGroup.POST("/claim", handler.MakeClaim)

	// auth
	router.POST("/auth", handler.SendAuth)
	router.POST("/auth/verify", handler.VerifyAuth)
	router.GET("/auth/me", handler.AuthMe)
	router.POST("/auth/login", handler.LoginWithPassword)
	router.POST("/auth/create", handler.CreateAccount)
	router.POST("/auth/code-to-reset-password", handler.SendCodeToResetPassword)
	router.POST("/auth/reset-password", handler.ResetPassword)
	router.GET("/auth/logout", handler.Logout)

	return router
}

func getTestUserToken(testUser models.User) (string, error) {
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
		return "", err
	}
	return tokenString, nil
}
