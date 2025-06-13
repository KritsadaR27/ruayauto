package handler

import (
	"net/http"
	"strconv"

	"ruaymanagement/backend/internal/ChatBotCore/models"
	service "ruaymanagement/backend/internal/ChatBotCore/services"

	"github.com/gin-gonic/gin"
)

// ChatBotHandler handles HTTP requests for chatbot operations
type ChatBotHandler struct {
	chatbotService *service.ChatBotService
}

// NewChatBotHandler creates a new chatbot handler
func NewChatBotHandler(chatbotService *service.ChatBotService) *ChatBotHandler {
	return &ChatBotHandler{
		chatbotService: chatbotService,
	}
}

// ProcessMessage handles incoming messages for processing
// POST /api/v1/messages/process
func (h *ChatBotHandler) ProcessMessage(c *gin.Context) {
	var req models.ProcessMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request format: " + err.Error(),
			},
		})
		return
	}

	// Validate required fields
	if req.FacebookUserID == "" || req.MessageText == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "MISSING_REQUIRED_FIELDS",
				Message: "facebook_user_id and message_text are required",
			},
		})
		return
	}

	// Process the message
	response, err := h.chatbotService.ProcessMessage(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "PROCESSING_ERROR",
				Message: "Failed to process message: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

// GetConversations retrieves conversations with pagination
// GET /api/v1/conversations
func (h *ChatBotHandler) GetConversations(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	// Get conversations
	conversations, total, err := h.chatbotService.GetConversations(c.Request.Context(), offset, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to retrieve conversations: " + err.Error(),
			},
		})
		return
	}

	// Calculate total pages
	totalPages := (total + perPage - 1) / perPage

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    conversations,
		Meta: &models.Meta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// GetConversation retrieves a specific conversation
// GET /api/v1/conversations/:id
func (h *ChatBotHandler) GetConversation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_ID",
				Message: "Invalid conversation ID",
			},
		})
		return
	}

	conversation, err := h.chatbotService.GetConversation(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "CONVERSATION_NOT_FOUND",
				Message: "Conversation not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    conversation,
	})
}

// UpdateConversationStatus updates conversation status
// PUT /api/v1/conversations/:id/status
func (h *ChatBotHandler) UpdateConversationStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_ID",
				Message: "Invalid conversation ID",
			},
		})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request format: " + err.Error(),
			},
		})
		return
	}

	// Validate status value
	validStatuses := map[string]bool{
		"active":   true,
		"inactive": true,
		"archived": true,
		"blocked":  true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_STATUS",
				Message: "Status must be one of: active, inactive, archived, blocked",
			},
		})
		return
	}

	err = h.chatbotService.UpdateConversationStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "UPDATE_ERROR",
				Message: "Failed to update conversation status: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    gin.H{"message": "Conversation status updated successfully"},
	})
}

// HealthCheck provides a health check endpoint
// GET /api/v1/health
func (h *ChatBotHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"status":  "healthy",
			"service": "ChatBotCore",
			"version": "1.0.0",
		},
	})
}
