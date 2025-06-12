package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"ruaymanagement/backend/external/FacebookConnect/config"
	"ruaymanagement/backend/external/FacebookConnect/handlers"
	"ruaymanagement/backend/external/FacebookConnect/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Initialize services
	facebookAPI := services.NewFacebookAPIClient()
	chatBotClient := services.NewChatBotCoreClient(cfg.ChatBotCoreURL)
	transformer := services.NewMessageTransformer()

	// Initialize handlers
	webhookHandler := handlers.NewWebhookHandler(
		facebookAPI,
		chatBotClient,
		transformer,
		cfg.VerifyToken,
		cfg.PageToken,
	)

	// Setup router
	router := setupRouter(webhookHandler)

	// Setup server
	server := &http.Server{
		Addr:         cfg.GetPort(),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("🚀 FacebookConnect service starting on port %s", cfg.Port)
	log.Printf("📡 ChatBotCore URL: %s", cfg.ChatBotCoreURL)

	// Start server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupRouter configures all routes for the FacebookConnect service
func setupRouter(webhookHandler *handlers.WebhookHandler) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "FacebookConnect",
			"version":   "1.0.0",
			"timestamp": time.Now().Unix(),
		})
	})

	// Facebook webhook endpoints
	router.GET("/webhook", webhookHandler.VerifyWebhook)
	router.POST("/webhook", webhookHandler.HandleWebhook)

	return router
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
