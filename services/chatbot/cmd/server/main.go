package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"chatbot/internal/config"
	"chatbot/internal/database"
	"chatbot/internal/handlers"
	"chatbot/internal/repository"
	"chatbot/internal/services"
)

func main() {
	log.Println("🚀 Starting Chatbot Service...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize layers
	keywordRepo := repository.NewKeywordRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatbotService := services.NewChatbotService(keywordRepo, messageRepo)
	chatbotHandler := handlers.NewChatbotHandler(chatbotService)
	keywordHandler := handlers.NewKeywordHandler(keywordRepo)
	messageHandler := handlers.NewMessageHandler(messageRepo)

	// Setup routes
	r := gin.Default()
	
	// Health check - support both GET and HEAD methods
	r.GET("/health", chatbotHandler.Health)
	r.HEAD("/health", chatbotHandler.Health)
	
	// Core chatbot functionality
	r.POST("/process", chatbotHandler.ProcessMessage)
	
	// API routes for web interface
	api := r.Group("/api")
	{
		// Keywords management
		api.GET("/keywords", keywordHandler.GetKeywords)
		api.POST("/keywords", keywordHandler.BulkCreateKeywords)
		api.PUT("/keywords/:id", keywordHandler.UpdateKeyword)
		api.DELETE("/keywords/:id", keywordHandler.DeleteKeyword)
		
		// Messages for AI analysis
		api.GET("/messages", messageHandler.GetMessages)
	}

	log.Printf("🎯 Server starting on %s", cfg.GetAddr())
	if err := r.Run(cfg.GetAddr()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
