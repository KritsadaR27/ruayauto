package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ruaychatbot/backend/internal/ChatBotCore/config"
	"ruaychatbot/backend/internal/ChatBotCore/database"
	"ruaychatbot/backend/internal/ChatBotCore/handlers"
	"ruaychatbot/backend/internal/ChatBotCore/repository"
	service "ruaychatbot/backend/internal/ChatBotCore/services"
	"ruaychatbot/backend/libs/monitoring/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	logger.Init("info")
	log.Println("🚀 Starting ChatBotCore service with Channel Adapter Pattern...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Initialize database
	dbManager, err := database.NewManagerFromURL(cfg.GetDatabaseURL(), "./migrations", true)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer dbManager.Shutdown(context.Background())

	// Initialize repositories
	repos := repository.NewRepositories(dbManager.GetDB())

	// Initialize service layer repositories
	serviceRepos := service.NewServiceRepositories(repos)

	// Create Gin router with dependencies
	router := setupRouter(serviceRepos)

	// Configure server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Printf("✅ ChatBotCore server running on http://localhost:%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited")
}

// setupRouter configures Gin router with all routes and Channel Adapter Pattern
func setupRouter(repos *service.Repositories) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Hub-Signature, X-Hub-Signature-256")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":       "healthy",
			"service":      "ChatBotCore",
			"version":      "1.0.0",
			"architecture": "Channel Adapter Pattern",
			"timestamp":    time.Now().Unix(),
		})
	})

	// Initialize handlers
	messageProcessorHandler := handler.NewMessageProcessorHandler(repos)
	ruleHandler := handler.NewRuleHandler(repos)
	pageHandler := handler.NewPageHandler(repos)
	keywordHandler := handler.NewKeywordHandler(repos.Keyword)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Message Processing (Channel Adapter Pattern)
		messageProcessorHandler.RegisterRoutes(v1)

		// Rule Management (Multi-page support)
		ruleHandler.RegisterRoutes(v1)

		// Page Management
		pageHandler.RegisterRoutes(v1)

		// Legacy Keywords endpoints (backward compatibility)
		keywords := v1.Group("/keywords")
		{
			keywords.GET("", keywordHandler.GetKeywords)
			keywords.POST("", keywordHandler.CreateKeyword)
			keywords.PUT("/:id", keywordHandler.UpdateKeyword)
			keywords.DELETE("/:id", keywordHandler.DeleteKeyword)
		}
	}

	// Legacy webhook endpoint (redirect to new processor)
	router.POST("/webhook/facebook", messageProcessorHandler.HandleFacebookWebhook)
	router.GET("/webhook/facebook", messageProcessorHandler.HandleFacebookVerification)

	log.Println("🔧 Router configured with Channel Adapter Pattern")
	log.Println("📱 Registered channels: Facebook")
	log.Println("📄 Multi-page support: Enabled")
	log.Println("🔧 Rule Engine: Active")

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
