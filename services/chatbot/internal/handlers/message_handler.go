package handlers

import (
	"net/http"
	"strconv"

	"chatbot/internal/repository"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	repo *repository.MessageRepository
}

func NewMessageHandler(repo *repository.MessageRepository) *MessageHandler {
	return &MessageHandler{repo: repo}
}

// GetMessages returns recent messages
func (h *MessageHandler) GetMessages(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 100
	}

	platform := c.Query("platform")

	var messages interface{}
	if platform != "" {
		messages, err = h.repo.GetByPlatform(c.Request.Context(), platform, limit)
	} else {
		messages, err = h.repo.GetAll(c.Request.Context(), limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      messages,
		"timestamp": "now",
	})
}
