package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ruaychatbot/backend/libs/monitoring/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	logger.Init("info")
	logger.Info("🚀 Starting FacebookConnect service...")

	// Create Gin router
	router := setupRouter()

	// Configure server
	srv := &http.Server{
		Addr:    ":8091", // Changed port to avoid conflicts
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		logger.Info("✅ FacebookConnect server running on http://localhost:8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("❌ Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("🛑 Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("❌ Server forced to shutdown:", err)
	}

	logger.Info("✅ Server exited")
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "FacebookConnect",
			"version":   "1.0.0",
			"timestamp": time.Now().Unix(),
		})
	})

	// Facebook webhook endpoints
	router.GET("/webhook", verifyWebhook)
	router.POST("/webhook", handleWebhook)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Send message to Facebook
		v1.POST("/messages/send", sendMessage)

		// Page management
		v1.GET("/pages", getPages)
	}

	return router
}

// Webhook verification (for Facebook)
func verifyWebhook(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	// Mock verification - replace with real token validation
	expectedToken := "your_verify_token"

	if mode == "subscribe" && token == expectedToken {
		logger.Info("Webhook verified successfully")
		c.String(http.StatusOK, challenge)
	} else {
		logger.Warn("Webhook verification failed")
		c.JSON(http.StatusForbidden, gin.H{"error": "Verification failed"})
	}
}

// Handle incoming webhook from Facebook
func handleWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	logger.Info("Received webhook payload:", payload)

	// Mock processing - replace with real webhook handling
	// TODO: Process Facebook messages and forward to ChatBotCore

	c.JSON(http.StatusOK, gin.H{"status": "EVENT_RECEIVED"})
}

// Send message to Facebook
func sendMessage(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	// Mock send message - replace with real Facebook API call
	logger.Info("Sending message:", req)

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message_id": "mock_message_123",
		"data":       req,
	})
}

// Get Facebook pages
func getPages(c *gin.Context) {
	// Mock pages data - replace with real Facebook API call
	pages := []map[string]interface{}{
		{
			"id":           "page_123",
			"name":         "Test Page",
			"access_token": "mock_token",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pages,
	})
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
