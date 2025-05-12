package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SendTextMessage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Text message sent",
	})
}
