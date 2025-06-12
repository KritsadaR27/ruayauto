package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"ruaymanagement/backend/external/FacebookConnect/models"
	"ruaymanagement/backend/external/FacebookConnect/services"
)

// WebhookHandler handles Facebook webhook requests
type WebhookHandler struct {
	facebookAPI   *services.FacebookAPIClient
	chatBotClient *services.ChatBotCoreClient
	transformer   *services.MessageTransformer
	verifyToken   string
	pageToken     string
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(
	facebookAPI *services.FacebookAPIClient,
	chatBotClient *services.ChatBotCoreClient,
	transformer *services.MessageTransformer,
	verifyToken string,
	pageToken string,
) *WebhookHandler {
	return &WebhookHandler{
		facebookAPI:   facebookAPI,
		chatBotClient: chatBotClient,
		transformer:   transformer,
		verifyToken:   verifyToken,
		pageToken:     pageToken,
	}
}

// VerifyWebhook handles GET /webhook - Facebook webhook verification
func (h *WebhookHandler) VerifyWebhook(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	log.Printf("📋 Webhook verification attempt - Mode: %s, Token: %s", mode, token[:min(8, len(token))]+"...")

	challengeResponse, err := h.facebookAPI.VerifyWebhook(mode, token, challenge, h.verifyToken)
	if err != nil {
		log.Printf("❌ Webhook verification failed: %v", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "Verification failed"})
		return
	}

	log.Printf("✅ Webhook verification successful")
	c.String(http.StatusOK, challengeResponse)
}

// HandleWebhook handles POST /webhook - Facebook webhook events
func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	var payload models.FacebookWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("❌ Error binding webhook payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	log.Printf("📩 Facebook webhook received - Object: %s, Entries: %d", payload.Object, len(payload.Entry))

	// Process webhook asynchronously to respond quickly to Facebook
	go h.processWebhookAsync(payload)

	// Respond immediately to Facebook
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

// processWebhookAsync processes webhook payload asynchronously
func (h *WebhookHandler) processWebhookAsync(payload models.FacebookWebhookPayload) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			if err := h.processChange(ctx, change); err != nil {
				log.Printf("⚠️ Error processing change: %v", err)
				continue
			}
		}
	}
}

// processChange processes a single webhook change
func (h *WebhookHandler) processChange(ctx context.Context, change models.FacebookWebhookChange) error {
	log.Printf("🔄 Processing change - Field: %s, Message: %s", change.Field, change.Value.Message)

	// Validate message
	if !h.transformer.ValidateMessage(change) {
		log.Printf("⏭️ Skipping message - validation failed")
		return nil
	}

	// Transform to internal format
	internalMessage := h.transformer.FacebookToInternal(change)

	// Send to ChatBotCore for processing
	response, err := h.chatBotClient.ProcessMessage(ctx, internalMessage)
	if err != nil {
		return fmt.Errorf("failed to process message in ChatBotCore: %w", err)
	}

	// Send reply if needed
	if response.Success && response.ShouldReply && response.Reply != "" {
		commentID := change.Value.CommentID
		if commentID == "" {
			log.Printf("⚠️ No comment ID for reply")
			return nil
		}

		replyReq := models.FacebookReplyRequest{
			CommentID: commentID,
			Message:   response.Reply,
			PageToken: h.pageToken,
		}

		if err := h.facebookAPI.SendCommentReply(ctx, replyReq); err != nil {
			return fmt.Errorf("failed to send Facebook reply: %w", err)
		}

		log.Printf("✅ Reply sent successfully - Keyword: %s, Response: %s",
			response.MatchedKeyword, response.Reply)
	} else {
		log.Printf("🤷 No reply needed - Success: %t, ShouldReply: %t",
			response.Success, response.ShouldReply)
	}

	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
