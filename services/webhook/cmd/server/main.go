package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"webhook/integrations/facebook"
	"webhook/integrations/line"
	"webhook/internal/config"
	"webhook/internal/handler"
	"webhook/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting webhook service on %s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Chatbot service URL: %s", cfg.Chatbot.URL)
	log.Printf("Facebook enabled: %v", cfg.Facebook.Enabled)
	log.Printf("LINE enabled: %v", cfg.Line.Enabled)

	// Initialize platform handlers
	var facebookHandler *facebook.Handler
	var lineHandler *line.Handler

	if cfg.Facebook.Enabled {
		facebookHandler = facebook.NewHandler(
			cfg.Facebook.VerifyToken,
			cfg.Facebook.PageToken,
			cfg.Facebook.AppSecret,
		)
		log.Println("Facebook handler initialized")
	}

	if cfg.Line.Enabled {
		lineHandler = line.NewHandler(
			cfg.Line.ChannelSecret,
			cfg.Line.ChannelToken,
		)
		log.Println("LINE handler initialized")
	}

	// Initialize webhook service
	webhookService := service.NewWebhookService(
		cfg.Chatbot.URL,
		facebookHandler,
		lineHandler,
	)

	// Initialize webhook handler
	webhookHandler := handler.NewWebhookHandler(
		cfg.Facebook.VerifyToken,
		cfg.Facebook.PageToken,
		cfg.Facebook.AppSecret,
		cfg.Line.ChannelSecret,
		cfg.Line.ChannelToken,
		webhookService,
	)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", webhookHandler.HealthCheck)

	// Webhook endpoints
	api := router.Group("/webhook")
	{
		if cfg.Facebook.Enabled {
			api.Any("/facebook", webhookHandler.HandleFacebookWebhook)
		}
		if cfg.Line.Enabled {
			api.POST("/line", webhookHandler.HandleLINEWebhook)
		}
		
		// Placeholder endpoints for future platforms
		api.POST("/instagram", webhookHandler.HandleInstagramWebhook)
		api.POST("/tiktok", webhookHandler.HandleTikTokWebhook)
		api.POST("/twitter", webhookHandler.HandleTwitterWebhook)
	}

	// Start server
	address := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("🚀 Webhook service starting on %s", address)
	log.Println("Available endpoints:")
	log.Println("  GET  /health")
	if cfg.Facebook.Enabled {
		log.Println("  GET/POST  /webhook/facebook")
	}
	if cfg.Line.Enabled {
		log.Println("  POST /webhook/line")
	}
	log.Println("  POST /webhook/instagram (placeholder)")
	log.Println("  POST /webhook/tiktok (placeholder)")
	log.Println("  POST /webhook/twitter (placeholder)")

	server := &http.Server{
		Addr:    address,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
