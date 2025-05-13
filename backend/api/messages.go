package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TextMessageRequest struct {
	Numbers []string `json:"numbers" binding:"required"`
	Message string   `json:"message" binding:"required"`
}

func (h *Handler) SendTextMessage(c *gin.Context) {
	var request TextMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(request)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Text message sent",
	})
}
