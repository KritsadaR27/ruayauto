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
	facebookRepo := repository.NewFacebookRepository(db)
	
	chatbotService := services.NewChatbotService(ruleRepo, messageRepo)
	facebookService := services.NewFacebookService(facebookRepo, cfg.FacebookAppID, cfg.FacebookAppSecret, "v19.0")
	
	chatbotHandler := handlers.NewChatbotHandler(chatbotService)
	ruleHandler := handlers.NewRuleHandler(ruleRepo)
	messageHandler := handlers.NewMessageHandler(messageRepo)
	facebookHandler := handlers.NewFacebookHandler(facebookService)

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
		
		// Facebook integration
		facebook := api.Group("/facebook")
		{
			// OAuth flow
			facebook.POST("/auth/login", facebookHandler.AuthLogin)
			facebook.GET("/auth/callback", facebookHandler.AuthCallback)
			facebook.GET("/auth/status", facebookHandler.AuthStatus)
			facebook.POST("/auth/logout", facebookHandler.AuthLogout)
			
			// Page management
			facebook.GET("/pages", facebookHandler.GetPages)
			facebook.POST("/pages/connect", facebookHandler.ConnectPage)
			facebook.POST("/pages/disconnect", facebookHandler.DisconnectPage)
			facebook.GET("/pages/connected", facebookHandler.GetConnectedPages)
		}
	}
	
	// Facebook webhook endpoints (outside of /api group)
	r.GET("/webhook/facebook", facebookHandler.WebhookVerify)
	r.POST("/webhook/facebook", facebookHandler.WebhookReceive)

	log.Printf("🎯 Server starting on %s", cfg.GetAddr())
	if err := r.Run(cfg.GetAddr()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
