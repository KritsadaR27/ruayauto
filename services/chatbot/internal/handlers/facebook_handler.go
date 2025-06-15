package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"chatbot/internal/models"
	"chatbot/internal/services"
)

type FacebookHandler struct {
	facebookService services.FacebookService
}

func NewFacebookHandler(facebookService services.FacebookService) *FacebookHandler {
	return &FacebookHandler{
		facebookService: facebookService,
	}
}

// OAuth endpoints
func (h *FacebookHandler) AuthLogin(c *gin.Context) {
	var req models.FacebookOAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.RedirectURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "redirect_uri is required"})
		return
	}

	authURL := h.facebookService.GenerateAuthURL(req.RedirectURI, req.State)
	
	c.JSON(http.StatusOK, models.FacebookOAuthResponse{
		AuthURL: authURL,
		State:   req.State,
	})
}

func (h *FacebookHandler) AuthCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is required"})
		return
	}

	// For now, use a default redirect URI - this should be configurable
	redirectURI := "http://localhost:3008/auth/facebook/callback"
	
	session, err := h.facebookService.HandleCallback(code, state, redirectURI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process callback: " + err.Error()})
		return
	}

	// In a real implementation, you would set a session cookie here
	// For now, just return the user ID
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user_id": session.FacebookUserID,
	})
}

func (h *FacebookHandler) AuthStatus(c *gin.Context) {
	// In a real implementation, you would get the user ID from session/cookie
	// For now, accept it as a query parameter
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusOK, models.FacebookAuthStatus{
			IsAuthenticated: false,
		})
		return
	}

	status, err := h.facebookService.GetAuthStatus(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get auth status: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *FacebookHandler) AuthLogout(c *gin.Context) {
	// In a real implementation, you would get the user ID from session/cookie
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	err := h.facebookService.Logout(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Page management endpoints
func (h *FacebookHandler) GetPages(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	pages, err := h.facebookService.GetUserPages(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pages: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pages": pages})
}

func (h *FacebookHandler) ConnectPage(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id"`
		PageID string `json:"page_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.UserID == "" || req.PageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and page_id are required"})
		return
	}

	err := h.facebookService.ConnectPage(req.UserID, req.PageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect page: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *FacebookHandler) DisconnectPage(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id"`
		PageID string `json:"page_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.UserID == "" || req.PageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and page_id are required"})
		return
	}

	err := h.facebookService.DisconnectPage(req.UserID, req.PageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disconnect page: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *FacebookHandler) GetConnectedPages(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	pages, err := h.facebookService.GetConnectedPages(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get connected pages: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pages": pages})
}

// Webhook endpoints
func (h *FacebookHandler) WebhookVerify(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" {
		// For now, use a simple token verification
		// In production, this should be more secure
		if token == "your_verify_token" {
			c.String(http.StatusOK, challenge)
			return
		}
	}

	c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
}

func (h *FacebookHandler) WebhookReceive(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Log the webhook payload for debugging
	payloadJSON, _ := json.MarshalIndent(payload, "", "  ")
	c.Header("Content-Type", "application/json")
	
	// For now, just acknowledge receipt
	// In production, this would process the webhook and trigger responses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Webhook received",
		"payload": string(payloadJSON),
	})
}
