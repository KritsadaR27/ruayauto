package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"ruaymanagement/libs/monitoring/logger"
)

func main() {
	// Initialize logger
	logger.Init("info")
	logger.Info("🚀 Starting ChatBotCore service...")

	// Create Gin router
	router := setupRouter()

	// Configure server
	srv := &http.Server{
		Addr:    ":8090", // Changed port to avoid conflicts
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		logger.Info("✅ ChatBotCore server running on http://localhost:8080")
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

// setupRouter configures Gin router with all routes
func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "ChatBotCore",
			"version":   "1.0.0",
			"timestamp": time.Now().Unix(),
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Keywords endpoints
		v1.GET("/keywords", getKeywords)
		v1.POST("/keywords", createKeyword)
		v1.PUT("/keywords/:id", updateKeyword)
		v1.DELETE("/keywords/:id", deleteKeyword)

		// Messages endpoints
		v1.POST("/messages/process", processMessage)
	}

	return router
}

// Mock handlers - replace with real business logic later
func getKeywords(c *gin.Context) {
	keywords := []map[string]interface{}{
		{
			"id":       1,
			"keyword":  "สวัสดี",
			"response": "สวัสดีครับ! มีอะไรให้ช่วยไหมครับ?",
			"active":   true,
		},
		{
			"id":       2,
			"keyword":  "ขอบคุณ",
			"response": "ยินดีครับ! มีอะไรให้ช่วยอีกไหมครับ?",
			"active":   true,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    keywords,
	})
}

func createKeyword(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	// Mock creation
	req["id"] = 3
	req["active"] = true

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    req,
	})
}

func updateKeyword(c *gin.Context) {
	id := c.Param("id")

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	// Mock update
	req["id"] = id

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    req,
	})
}

func deleteKeyword(c *gin.Context) {
	id := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Keyword " + id + " deleted successfully",
	})
}

func processMessage(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	message, ok := req["message"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Message is required",
		})
		return
	}

	// Simple keyword matching
	var response string
	switch message {
	case "สวัสดี", "hello", "hi":
		response = "สวัสดีครับ! มีอะไรให้ช่วยไหมครับ?"
	case "ขอบคุณ", "thank you", "thanks":
		response = "ยินดีครับ! มีอะไรให้ช่วยอีกไหมครับ?"
	default:
		response = "ขออภัยครับ ผมไม่เข้าใจคำถามของคุณ กรุณาลองถามใหม่ครับ"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"original_message": message,
			"response":         response,
			"matched":          response != "ขออภัยครับ ผมไม่เข้าใจคำถามของคุณ กรุณาลองถามใหม่ครับ",
		},
	})
}
