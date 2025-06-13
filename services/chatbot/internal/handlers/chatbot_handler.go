package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"chatbot/internal/models"
	"chatbot/internal/services"
)

type ChatbotHandler struct {
	service *services.ChatbotService
}

func NewChatbotHandler(service *services.ChatbotService) *ChatbotHandler {
	return &ChatbotHandler{service: service}
}

func (h *ChatbotHandler) ProcessMessage(c *gin.Context) {
	var req models.AutoReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	resp, err := h.service.ProcessMessage(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

func (h *ChatbotHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "chatbot",
		"port":    "8090",
	})
}
