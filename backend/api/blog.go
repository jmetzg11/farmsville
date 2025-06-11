package api

import (
	"farmsville/backend/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostBlog(c *gin.Context) {
	var postBlogRequest struct {
		Title   string `json:"title" binding:"required"`
		Content []struct {
		}
	}
}
