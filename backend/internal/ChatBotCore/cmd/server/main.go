package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"ruaymanagement/backend/internal/ChatBotCore/config"
	"ruaymanagement/backend/internal/ChatBotCore/handlers"
	"ruaymanagement/backend/internal/ChatBotCore/repository"
	"ruaymanagement/backend/internal/ChatBotCore/services"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup database connection
	db, err := setupDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	repos := initializeRepositories(db)

	// Initialize services
	chatBotService := services.NewChatBotService(repos, cfg)

	// Setup Gin router
	router := setupRouter(cfg, chatBotService, repos)

	// Start server
	log.Printf("🚀 ChatBotCore server starting on port %s", cfg.Port)
	log.Printf("📊 Database connected successfully")

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupDatabase initializes database connection
func setupDatabase(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.GetDatabaseURL())
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(cfg.MaxConnections / 2)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Printf("✅ Database connection established")
	return db, nil
}

// initializeRepositories creates all repository instances
func initializeRepositories(db *sql.DB) *repository.Repositories {
	return &repository.Repositories{
		Conversation:     repository.NewConversationRepository(db),
		Message:          repository.NewMessageRepository(db),
		Keyword:          repository.NewKeywordRepository(db),
		ResponseTemplate: repository.NewResponseTemplateRepository(db),
		Analytics:        repository.NewAnalyticsRepository(db),
	}
}

// setupRouter configures Gin router with all routes
func setupRouter(cfg *config.Config, chatBotService *services.ChatBotService, repos *repository.Repositories) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	router := gin.Default()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "ChatBotCore",
			"version":   "1.0.0",
			"timestamp": time.Now().Unix(),
		})
	})

	// Initialize handlers
	chatHandler := handlers.NewChatBotHandler(chatBotService)
	keywordHandler := handlers.NewKeywordHandler(repos.Keyword)
	messageHandler := handlers.NewMessageHandler(chatBotService) // New handler for external services

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", chatHandler.HealthCheck)

		// Message processing (for external services like FacebookConnect)
		v1.POST("/messages/process", messageHandler.ProcessMessage)

		// Conversations
		conversations := v1.Group("/conversations")
		{
			conversations.GET("", chatHandler.GetConversations)
			conversations.GET("/:id", chatHandler.GetConversation)
			conversations.PUT("/:id/status", chatHandler.UpdateConversationStatus)
		}

		// Keywords management
		keywords := v1.Group("/keywords")
		{
			keywords.POST("", keywordHandler.CreateKeyword)
			keywords.GET("", keywordHandler.GetKeywords)
			keywords.GET("/:id", keywordHandler.GetKeyword)
			keywords.PUT("/:id", keywordHandler.UpdateKeyword)
			keywords.DELETE("/:id", keywordHandler.DeleteKeyword)
			keywords.PATCH("/:id/toggle", keywordHandler.ToggleKeywordStatus)
		}
	}

	return router
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}
