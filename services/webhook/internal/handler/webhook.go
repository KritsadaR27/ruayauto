package handler

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"webhook/integrations/facebook"
	"webhook/integrations/line"
	"webhook/internal/service"
)

// WebhookHandler handles all platform webhooks
type WebhookHandler struct {
	facebook       *facebook.Handler
	line           *line.Handler
	webhookService *service.WebhookService
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(
	fbVerifyToken, fbPageToken, fbAppSecret string,
	lineChannelSecret, lineChannelToken string,
	webhookService *service.WebhookService,
) *WebhookHandler {
	return &WebhookHandler{
		facebook:       facebook.NewHandler(fbVerifyToken, fbPageToken, fbAppSecret),
		line:           line.NewHandler(lineChannelSecret, lineChannelToken),
		webhookService: webhookService,
	}
}

// HandleFacebookWebhook handles Facebook webhook verification and events
func (h *WebhookHandler) HandleFacebookWebhook(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		h.handleFacebookVerification(c)
	case "POST":
		h.handleFacebookEvent(c)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	}
}

// handleFacebookVerification handles Facebook webhook verification challenge
func (h *WebhookHandler) handleFacebookVerification(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" {
		if err := h.facebook.VerifyWebhook(token); err != nil {
			log.Printf("Facebook verification failed: %v", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "Verification failed"})
			return
		}

		// Return the challenge as plain text
		challengeInt, err := strconv.Atoi(challenge)
		if err != nil {
			c.String(http.StatusOK, challenge)
		} else {
			c.JSON(http.StatusOK, challengeInt)
		}
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification request"})
}

// handleFacebookEvent handles Facebook webhook events
func (h *WebhookHandler) handleFacebookEvent(c *gin.Context) {
	// Read the raw body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read Facebook webhook body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Reset body for potential re-reading
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	// Get signature from header
	signature := c.GetHeader("X-Hub-Signature-256")

	// Process webhook
	messages, err := h.facebook.HandleWebhook(body, signature)
	if err != nil {
		log.Printf("Failed to process Facebook webhook: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to process webhook"})
		return
	}

	// Process each message
	for _, msg := range messages {
		if err := h.webhookService.ProcessMessage(msg); err != nil {
			log.Printf("Failed to process Facebook message: %v", err)
			// Continue processing other messages
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// HandleLINEWebhook handles LINE webhook events
func (h *WebhookHandler) HandleLINEWebhook(c *gin.Context) {
	// Read the raw body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read LINE webhook body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get signature from header
	signature := c.GetHeader("X-Line-Signature")

	// Process webhook
	messages, err := h.line.HandleWebhook(body, signature)
	if err != nil {
		log.Printf("Failed to process LINE webhook: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to process webhook"})
		return
	}

	// Process each message
	for _, msg := range messages {
		if err := h.webhookService.ProcessMessage(msg); err != nil {
			log.Printf("Failed to process LINE message: %v", err)
			// Continue processing other messages
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// HandleInstagramWebhook handles Instagram webhook events
func (h *WebhookHandler) HandleInstagramWebhook(c *gin.Context) {
	// TODO: Implement Instagram webhook handling
	log.Println("Instagram webhook received - not implemented yet")
	c.JSON(http.StatusOK, gin.H{"status": "not_implemented"})
}

// HandleTikTokWebhook handles TikTok webhook events
func (h *WebhookHandler) HandleTikTokWebhook(c *gin.Context) {
	// TODO: Implement TikTok webhook handling
	log.Println("TikTok webhook received - not implemented yet")
	c.JSON(http.StatusOK, gin.H{"status": "not_implemented"})
}

// HandleTwitterWebhook handles Twitter webhook events
func (h *WebhookHandler) HandleTwitterWebhook(c *gin.Context) {
	// TODO: Implement Twitter webhook handling
	log.Println("Twitter webhook received - not implemented yet")
	c.JSON(http.StatusOK, gin.H{"status": "not_implemented"})
}

// HealthCheck provides health check endpoint
func (h *WebhookHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "webhook",
		"platforms": []string{"facebook", "line", "instagram", "tiktok", "twitter"},
	})
}
