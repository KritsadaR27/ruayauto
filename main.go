package main

import (
	"ruayAutoMsg/handler"

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

	// 🚀 Run server ที่ port 8100
	r.Run(":8100")
}
