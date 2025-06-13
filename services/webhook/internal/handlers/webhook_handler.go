package handlers

import (
	"encoding/json"
	"net/http"
	"io"
	
	"github.com/gin-gonic/gin"
	"webhook/internal/models"
	"webhook/internal/services"
	"webhook/integrations/facebook"
	"webhook/integrations/line"
)

// WebhookHandler handles incoming webhook requests
type WebhookHandler struct {
	webhookService *services.WebhookService
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(webhookService *services.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		webhookService: webhookService,
	}
}

// HandleFacebookWebhook handles Facebook webhook events
func (h *WebhookHandler) HandleFacebookWebhook(c *gin.Context) {
	// Handle Facebook webhook verification (GET request)
	if c.Request.Method == "GET" {
		h.handleFacebookVerification(c)
		return
	}
	
	// Handle actual webhook event (POST request)
	signature := c.GetHeader("X-Hub-Signature-256")
	
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	
	// Get Facebook handler and verify signature
	handler, exists := h.webhookService.GetHandler(models.PlatformFacebook)
	if !exists {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Facebook handler not registered"})
		return
	}
	
	if !handler.VerifyWebhook(signature, string(body)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}
	
	// Parse JSON body
	var rawData interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	
	// Create webhook event
	event := &models.WebhookEvent{
		Platform:  models.PlatformFacebook,
		EventType: "comment",
		RawData:   rawData,
		Signature: signature,
	}
	
	// Process webhook
	if err := h.webhookService.ProcessWebhook(c.Request.Context(), models.PlatformFacebook, event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// handleFacebookVerification handles Facebook webhook verification
func (h *WebhookHandler) handleFacebookVerification(c *gin.Context) {
	verifyToken := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")
	
	// TODO: Verify token against configured value
	if verifyToken == "your_verify_token" {
		c.String(http.StatusOK, challenge)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid verify token"})
	}
}

// HandleLineWebhook handles LINE webhook events
func (h *WebhookHandler) HandleLineWebhook(c *gin.Context) {
	signature := c.GetHeader("X-Line-Signature")
	
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	
	// Get LINE handler and verify signature
	handler, exists := h.webhookService.GetHandler(models.PlatformLine)
	if !exists {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "LINE handler not registered"})
		return
	}
	
	if !handler.VerifyWebhook(signature, string(body)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}
	
	// Parse JSON body
	var rawData interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	
	// Create webhook event
	event := &models.WebhookEvent{
		Platform:  models.PlatformLine,
		EventType: "message",
		RawData:   rawData,
		Signature: signature,
	}
	
	// Process webhook
	if err := h.webhookService.ProcessWebhook(c.Request.Context(), models.PlatformLine, event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// HandleTikTokWebhook handles TikTok webhook events
func (h *WebhookHandler) HandleTikTokWebhook(c *gin.Context) {
	// TODO: Implement TikTok webhook handler
	c.JSON(http.StatusNotImplemented, gin.H{"message": "TikTok webhook coming soon"})
}

// HandleInstagramWebhook handles Instagram webhook events  
func (h *WebhookHandler) HandleInstagramWebhook(c *gin.Context) {
	// TODO: Implement Instagram webhook handler
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Instagram webhook coming soon"})
}

// HandleTwitterWebhook handles Twitter webhook events
func (h *WebhookHandler) HandleTwitterWebhook(c *gin.Context) {
	// TODO: Implement Twitter webhook handler
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Twitter webhook coming soon"})
}

// HealthCheck returns service health status
func (h *WebhookHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "webhook-service",
		"platforms": []string{"facebook", "line", "tiktok", "instagram", "twitter"},
	})
}
