package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"ruaychatbot/backend/internal/ChatBotCore/channels"
	service "ruaychatbot/backend/internal/ChatBotCore/services"

	"github.com/gin-gonic/gin"
)

// MessageProcessorHandler handles unified message processing across all channels
type MessageProcessorHandler struct {
	processor       *service.MessageProcessor
	facebookChannel channels.Channel
}

// NewMessageProcessorHandler creates a new message processor handler
func NewMessageProcessorHandler(repos *service.Repositories) *MessageProcessorHandler {
	processor := service.NewMessageProcessor(repos)

	// Register Facebook channel
	facebookConfig := channels.ChannelConfig{
		Platform:      "facebook",
		Name:          "facebook",
		Enabled:       true,
		WebhookSecret: "your_webhook_secret", // TODO: Get from config
		Config: map[string]interface{}{
			"access_token": "your_access_token", // TODO: Get from config
		},
	}

	facebookChannel := channels.NewFacebookAdapter(facebookConfig)
	processor.RegisterChannel("facebook", facebookChannel)

	return &MessageProcessorHandler{
		processor:       processor,
		facebookChannel: facebookChannel,
	}
}

// HandleFacebookWebhook processes Facebook webhook requests
func (h *MessageProcessorHandler) HandleFacebookWebhook(c *gin.Context) {
	log.Println("📩 Facebook webhook received via Message Processor")

	// Get raw payload
	rawData, err := c.GetRawData()
	if err != nil {
		log.Printf("❌ Failed to get raw data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// Get signature from header
	signature := c.GetHeader("X-Hub-Signature-256")
	if signature == "" {
		signature = c.GetHeader("X-Hub-Signature") // Fallback for SHA1
	}

	// Process webhook
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := h.processor.ProcessWebhook(ctx, "facebook", rawData, signature); err != nil {
		log.Printf("❌ Failed to process webhook: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "processing failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// HandleFacebookVerification handles Facebook webhook verification
func (h *MessageProcessorHandler) HandleFacebookVerification(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	expectedToken := "your_verify_token" // TODO: Get from config

	if mode == "subscribe" && token == expectedToken {
		log.Println("✅ Facebook webhook verified")
		c.String(http.StatusOK, challenge)
		return
	}

	log.Println("❌ Facebook webhook verification failed")
	c.JSON(http.StatusForbidden, gin.H{"error": "verification failed"})
}

// GetProcessorStats returns message processor statistics
func (h *MessageProcessorHandler) GetProcessorStats(c *gin.Context) {
	stats := h.processor.GetStats()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// GetChannelHealth checks health of all channels
func (h *MessageProcessorHandler) GetChannelHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.processor.HealthCheck(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "All channels healthy",
	})
}

// ProcessTestMessage processes a test message (for debugging)
func (h *MessageProcessorHandler) ProcessTestMessage(c *gin.Context) {
	var req struct {
		Platform string `json:"platform"`
		PageID   string `json:"page_id"`
		UserID   string `json:"user_id"`
		Content  string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Create test message
	testMessage := &channels.IncomingMessage{
		Platform:    req.Platform,
		PageID:      req.PageID,
		UserID:      req.UserID,
		UserName:    stringPtr("Test User"),
		MessageType: "message",
		Content:     req.Content,
		MessageID:   "test_message_id",
		ContentID:   "test_content_id",
		Timestamp:   time.Now(),
	}

	// Get channel
	channel, err := h.processor.GetChannel(req.Platform)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unsupported platform: %s", req.Platform)})
		return
	}

	// Process message
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := h.processor.ProcessMessage(ctx, testMessage, channel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Test message processed successfully",
	})
}

// RegisterRoutes registers all message processor routes
func (h *MessageProcessorHandler) RegisterRoutes(r *gin.RouterGroup) {
	// Facebook webhook routes
	facebook := r.Group("/facebook")
	{
		facebook.GET("/webhook", h.HandleFacebookVerification)
		facebook.POST("/webhook", h.HandleFacebookWebhook)
	}

	// Processor management routes
	processor := r.Group("/processor")
	{
		processor.GET("/stats", h.GetProcessorStats)
		processor.GET("/health", h.GetChannelHealth)
		processor.POST("/test", h.ProcessTestMessage)
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
