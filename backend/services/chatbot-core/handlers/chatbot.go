package handlers

import (
	"net/http"
	"time"

	"chatbot-core/models"
	"chatbot-core/services"

	"github.com/gin-gonic/gin"
)

// ChatBotHandler handles HTTP requests for chatbot functionality
type ChatBotHandler struct {
	chatbotService *services.ChatBotService
}

// NewChatBotHandler creates a new chatbot handler
func NewChatBotHandler(chatbotService *services.ChatBotService) *ChatBotHandler {
	return &ChatBotHandler{
		chatbotService: chatbotService,
	}
}

// ProcessMessage handles auto-reply processing requests
// POST /api/process-message
func (h *ChatBotHandler) ProcessMessage(c *gin.Context) {
	var req models.AutoReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate required fields
	if req.PageID == "" || req.UserID == "" || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields: page_id, user_id, content"})
		return
	}

	// Process the message
	response, err := h.chatbotService.ProcessMessage(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process message", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Health handles health check requests
// GET /health
func (h *ChatBotHandler) Health(c *gin.Context) {
	response := models.HealthResponse{
		Status:    "healthy",
		Service:   "ChatBotCore",
		Version:   "1.0.0",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}
