package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"ruaymanagement/backend/internal/ChatBotCore/models"
	"ruaymanagement/backend/internal/ChatBotCore/services"
)

// MessageHandler handles message processing requests
type MessageHandler struct {
	chatBotService *services.ChatBotService
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(chatBotService *services.ChatBotService) *MessageHandler {
	return &MessageHandler{
		chatBotService: chatBotService,
	}
}

// ProcessMessageRequest represents incoming message processing request
type ProcessMessageRequest struct {
	Content        string            `json:"content"`
	SenderID       string            `json:"sender_id"`
	ConversationID string            `json:"conversation_id"`
	Timestamp      time.Time         `json:"timestamp"`
	MessageType    string            `json:"message_type"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// ProcessMessageResponse represents the response for message processing
type ProcessMessageResponse struct {
	Success        bool   `json:"success"`
	Reply          string `json:"reply,omitempty"`
	MatchedKeyword string `json:"matched_keyword,omitempty"`
	ShouldReply    bool   `json:"should_reply"`
	Error          string `json:"error,omitempty"`
}

// ProcessMessage handles POST /api/messages/process
func (h *MessageHandler) ProcessMessage(c *gin.Context) {
	var req ProcessMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ProcessMessageResponse{
			Success: false,
			Error:   fmt.Sprintf("Invalid request format: %v", err),
		})
		return
	}

	// Validate required fields
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, ProcessMessageResponse{
			Success: false,
			Error:   "Content is required",
		})
		return
	}

	if req.SenderID == "" {
		c.JSON(http.StatusBadRequest, ProcessMessageResponse{
			Success: false,
			Error:   "SenderID is required",
		})
		return
	}

	// Create internal message processing request
	processReq := &models.ProcessMessageRequest{
		FacebookUserID:    req.SenderID,
		UserName:          req.SenderID, // Use sender ID as name for now
		MessageText:       req.Content,
		MessageType:       req.MessageType,
		MessageData: map[string]interface{}{
			"conversation_id": req.ConversationID,
			"timestamp":       req.Timestamp,
		},
	}
	
	// Add external metadata
	for k, v := range req.Metadata {
		processReq.MessageData[k] = v
	}

	ctx := context.Background()

	// Process message through ChatBotService
	response, err := h.chatBotService.ProcessMessage(ctx, processReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ProcessMessageResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to process message: %v", err),
		})
		return
	}

	// Prepare response
	result := ProcessMessageResponse{
		Success:     true,
		ShouldReply: response.ShouldReply,
	}

	if response.ShouldReply {
		result.Reply = response.ResponseContent
		if response.MatchedKeyword != nil {
			result.MatchedKeyword = *response.MatchedKeyword
		}
	}

	c.JSON(http.StatusOK, result)
}
