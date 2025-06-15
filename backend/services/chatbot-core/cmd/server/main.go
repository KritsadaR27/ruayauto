package main

import (
	"log"
	"net/http"
	"time"

	"chatbot-core/config"
	"chatbot-core/database"
	"chatbot-core/handlers"
	"chatbot-core/repository"
	"chatbot-core/services"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("🚀 Starting ChatBotCore service...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.GetDSN())
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	keywordRepo := repository.NewKeywordRepository(db)

	// Initialize services
	chatbotService := services.NewChatBotService(keywordRepo)

	// Initialize handlers
	chatbotHandler := handlers.NewChatBotHandler(chatbotService)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Middleware
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

	// Routes
	api := router.Group("/api")
	{
		api.POST("/process-message", chatbotHandler.ProcessMessage)
		api.GET("/health", chatbotHandler.Health)
	}

	// Health check route
	router.GET("/health", chatbotHandler.Health)

	// Start server
	server := &http.Server{
		Addr:           cfg.GetServerAddr(),
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("🌟 ChatBotCore service running on %s", cfg.GetServerAddr())
	log.Printf("📋 Health check: http://%s/health", cfg.GetServerAddr())
	log.Printf("🤖 Process message: POST http://%s/api/process-message", cfg.GetServerAddr())

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
