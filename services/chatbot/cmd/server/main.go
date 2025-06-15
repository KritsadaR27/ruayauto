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
	ruleRepo := repository.NewRuleRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatbotService := services.NewChatbotService(ruleRepo, messageRepo)
	chatbotHandler := handlers.NewChatbotHandler(chatbotService)
	ruleHandler := handlers.NewRuleHandler(ruleRepo)
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
		// Rules management
		api.GET("/rules", ruleHandler.GetRules)
		api.POST("/rules", ruleHandler.CreateRule)
		api.PUT("/rules/:id", ruleHandler.UpdateRule)
		api.DELETE("/rules/:id", ruleHandler.DeleteRule)
		
		// Rule responses management (Modern Random Response System)
		api.GET("/rules/:id/responses", ruleHandler.GetRuleResponses)
		api.POST("/rules/:id/responses", ruleHandler.CreateRuleResponse)
		api.PUT("/rules/:id/responses/:response_id", ruleHandler.UpdateRuleResponse)
		api.DELETE("/rules/:id/responses/:response_id", ruleHandler.DeleteRuleResponse)
		
		// Messages for AI analysis
		api.GET("/messages", messageHandler.GetMessages)
	}

	log.Printf("🎯 Server starting on %s", cfg.GetAddr())
	if err := r.Run(cfg.GetAddr()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
