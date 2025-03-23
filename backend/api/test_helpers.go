package api

import (
	"farmsville/backend/models"
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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

	// Keep connection alive for the duration of the test
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		sqlDB.Close()
	})

	err = db.AutoMigrate(&models.Item{}, &models.ClaimedItem{})
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func setupTestRouter(handler *Handler) *gin.Engine {
	router := gin.New()

	// main
	router.GET("/items", handler.GetItems)

	return router
}
