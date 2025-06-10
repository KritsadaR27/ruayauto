package main

import (
	"ruayautomsg/internal/handler"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 👉 สำหรับ Facebook webhook verification
	r.GET("/webhook/facebook", handler.FacebookVerifyHandler)

	// 👉 สำหรับรับข้อความจริง
	r.POST("/webhook/facebook", handler.FacebookWebhookHandler)

	// 👉 REST API สำหรับ frontend
	r.GET("/api/keywords", handler.GetKeywordsHandler)
	r.POST("/api/keywords", handler.SaveKeywordsHandler)

	// 🚀 Run server 
	port := os.Getenv("PORT")
	if port == "" {
		port = "3006"
	}
	r.Run(":" + port)
}
