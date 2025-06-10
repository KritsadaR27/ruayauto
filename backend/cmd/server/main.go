package main

import (
	"context"
	"log"
	"time"

	"ruayautomsg/internal/config"
	"ruayautomsg/internal/database"
	"ruayautomsg/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()

	// Initialize database manager
	mgr, err := database.NewManagerFromURL(cfg.GetDatabaseURL(), "./internal/migrations", true)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		mgr.Shutdown(ctx)
	}()

	// Setup graceful shutdown
	mgr.SetupGracefulShutdown()

	// Wait for database to be healthy
	if err := mgr.WaitForHealth(time.Second * 30); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}

	log.Printf("✅ Database connection established and healthy")

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		status := mgr.GetHealthStatus()
		if status.IsHealthy {
			c.JSON(200, gin.H{
				"status": "healthy",
				"database": status,
			})
		} else {
			c.JSON(503, gin.H{
				"status": "unhealthy",
				"database": status,
			})
		}
	})

	// 👉 สำหรับ Facebook webhook verification
	r.GET("/webhook/facebook", handler.FacebookVerifyHandler)

	// 👉 สำหรับรับข้อความจริง
	r.POST("/webhook/facebook", handler.FacebookWebhookHandler)

	// 👉 REST API สำหรับ frontend
	r.GET("/api/keywords", handler.GetKeywordsHandler)
	r.POST("/api/keywords", handler.SaveKeywordsHandler)

	// 🚀 Run server 
	log.Printf("🚀 Starting server on port %s", cfg.Port)
	log.Printf("📊 Database: %s", cfg.DBHost+":"+cfg.DBPort)
	log.Printf("🔧 Facebook Verify Token configured: %t", cfg.GetFacebookVerifyToken() != "")
	log.Printf("🔑 Facebook Page Token configured: %t", cfg.GetFacebookPageToken() != "")
	
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
