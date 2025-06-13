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
	chatbotService := services.NewChatbotService(keywordRepo)
	chatbotHandler := handlers.NewChatbotHandler(chatbotService)

	// Setup routes
	r := gin.Default()
	r.GET("/health", chatbotHandler.Health)
	r.POST("/process", chatbotHandler.ProcessMessage)

	log.Printf("🎯 Server starting on %s", cfg.GetAddr())
	if err := r.Run(cfg.GetAddr()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
