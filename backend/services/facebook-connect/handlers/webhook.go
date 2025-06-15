package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"facebook-connect/config"
	"facebook-connect/models"
	"facebook-connect/services"

	"github.com/gin-gonic/gin"
)

// WebhookHandler handles Facebook webhook requests
type WebhookHandler struct {
	cfg           *config.Config
	facebookAPI   *services.FacebookAPI
	chatbotClient *services.ChatBotClient
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(cfg *config.Config, facebookAPI *services.FacebookAPI, chatbotClient *services.ChatBotClient) *WebhookHandler {
	return &WebhookHandler{
		cfg:           cfg,
		facebookAPI:   facebookAPI,
		chatbotClient: chatbotClient,
	}
}

// VerifyWebhook handles webhook verification from Facebook
// GET /webhook
func (h *WebhookHandler) VerifyWebhook(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	log.Printf("🔍 Webhook verification - Mode: %s, Token: %s", mode, token)

	if mode == "subscribe" && token == h.cfg.VerifyToken {
		log.Println("✅ Webhook verified successfully")
		c.String(http.StatusOK, challenge)
		return
	}

	log.Println("❌ Webhook verification failed")
	c.JSON(http.StatusForbidden, gin.H{"error": "Verification failed"})
}

// HandleWebhook processes incoming webhook events from Facebook
// POST /webhook
func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	var payload models.FacebookWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("❌ Error parsing webhook payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	log.Printf("📥 Received webhook: Object=%s, Entries=%d", payload.Object, len(payload.Entry))

	// Process each entry
	for _, entry := range payload.Entry {
		h.processEntry(entry)
	}

	c.Status(http.StatusOK)
}

// processEntry processes a single webhook entry
func (h *WebhookHandler) processEntry(entry models.FacebookWebhookEntry) {
	log.Printf("🔄 Processing entry ID: %s", entry.ID)

	for _, change := range entry.Changes {
		if change.Field == "feed" {
			h.processFeedChange(entry.ID, change.Value)
		}
	}
}

// processFeedChange processes feed changes (comments)
func (h *WebhookHandler) processFeedChange(pageID string, value interface{}) {
	// Convert value to JSON for easier parsing
	valueBytes, err := json.Marshal(value)
	if err != nil {
		log.Printf("❌ Error marshaling feed change value: %v", err)
		return
	}

	var feedData map[string]interface{}
	if err := json.Unmarshal(valueBytes, &feedData); err != nil {
		log.Printf("❌ Error unmarshaling feed change value: %v", err)
		return
	}

	// Check if this is a comment
	if feedData["verb"] == "add" && feedData["item"] == "comment" {
		h.processComment(pageID, feedData)
	}
}

// processComment processes a new comment
func (h *WebhookHandler) processComment(pageID string, commentData map[string]interface{}) {
	// Extract comment information
	commentID, _ := commentData["comment_id"].(string)
	postID, _ := commentData["post_id"].(string)
	userID, _ := commentData["from"].(map[string]interface{})["id"].(string)
	message, _ := commentData["message"].(string)

	log.Printf("💬 New comment - Page: %s, Comment: %s, User: %s, Message: %s", pageID, commentID, userID, message)

	if commentID == "" || message == "" || userID == "" {
		log.Println("⚠️ Incomplete comment data, skipping")
		return
	}

	// Send to ChatBot service for processing
	chatBotReq := &models.ChatBotRequest{
		PageID:    pageID,
		UserID:    userID,
		Content:   message,
		CommentID: commentID,
		PostID:    postID,
	}

	chatBotResp, err := h.chatbotClient.ProcessMessage(c.Request.Context(), chatBotReq)
	if err != nil {
		log.Printf("❌ Error processing message with ChatBot: %v", err)
		return
	}

	// If ChatBot says we should reply, send the response via Facebook API
	if chatBotResp.ShouldReply && chatBotResp.Response != "" {
		_, err := h.facebookAPI.ReplyToComment(c.Request.Context(), commentID, chatBotResp.Response)
		if err != nil {
			log.Printf("❌ Error replying to comment: %v", err)
		}
	} else {
		log.Println("💭 No reply needed for this comment")
	}
}

// Health handles health check requests
// GET /health
func (h *WebhookHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "FacebookConnect",
		"version":   "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
