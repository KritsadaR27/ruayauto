package handler
package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"ruaychatbot/backend/internal/ChatBotCore/config"
	service "ruaychatbot/backend/internal/ChatBotCore/services"
	"github.com/gin-gonic/gin"
)

// FacebookWebhookRepoHandler handles Facebook webhook requests using repository pattern
type FacebookWebhookRepoHandler struct {
	autoReplyService *service.AutoReplyService
	keywordService   *service.KeywordService
}

// NewFacebookWebhookRepoHandler creates a new Facebook webhook handler
func NewFacebookWebhookRepoHandler(keywordService *service.KeywordService) *FacebookWebhookRepoHandler {
	autoReplyService := service.NewAutoReplyService(keywordService.GetRepository())
	
	return &FacebookWebhookRepoHandler{
		autoReplyService: autoReplyService,
		keywordService:   keywordService,
	}
}

// HandleWebhook processes incoming Facebook webhook requests
func (h *FacebookWebhookRepoHandler) HandleWebhook(c *gin.Context) {
	var payload FacebookWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fmt.Printf("❌ Error binding JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	fmt.Println("📩 Facebook webhook received (repository pattern)")

	ctx := context.Background()
	autoReplyConfig := h.autoReplyService.GetDefaultConfig()

	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			if err := h.processChange(ctx, change, autoReplyConfig); err != nil {
				fmt.Printf("⚠️ Error processing change: %v\n", err)
				continue
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

// processChange processes a single change from the webhook payload
func (h *FacebookWebhookRepoHandler) processChange(ctx context.Context, change struct {
Field string `json:"field"`
Value struct {
From struct {
ID   string `json:"id"`
Name string `json:"name"`
} `json:"from"`
Message   string `json:"message"`
CommentID string `json:"comment_id"`
PostID    string `json:"post_id"`
} `json:"value"`
}, autoReplyConfig *service.AutoReplyConfig) error {

	fmt.Printf("🧠 Field: %s\n", change.Field)
	text := change.Value.Message
	commentID := change.Value.CommentID
	fmt.Printf("🗨️ Message: %s | 💬 CommentID: %s\n", text, commentID)

	// Skip if no text or comment ID
	if strings.TrimSpace(text) == "" || commentID == "" {
		fmt.Println("⏭️ Skipping: empty text or comment ID")
		return nil
	}

	// Check for tags
	if strings.Contains(text, "@") {
		fmt.Println("🚫 Skipping: comment contains tag (@)")
		return nil
	}

	// Process auto-reply
	result, err := h.autoReplyService.ProcessAutoReply(ctx, text, autoReplyConfig)
	if err != nil {
		return fmt.Errorf("failed to process auto-reply: %w", err)
	}

	if result.Matched && result.Response != "" {
		// Post reply to Facebook
		cfg := config.Get()
		
		replyCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		err := h.autoReplyService.PostCommentReply(replyCtx, commentID, result.Response, cfg.GetFacebookPageToken())
		if err != nil {
			return fmt.Errorf("failed to post reply: %w", err)
		}

		fmt.Printf("✅ Replied with %s match: %s\n", result.MatchType, result.Response)
	} else {
		fmt.Println("🤷 No matching response found")
	}

	return nil
}
