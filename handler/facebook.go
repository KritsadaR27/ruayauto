package handler

import (
	"fmt"
	"net/http"
	"os"
	"ruayAutoMsg/utils"

	"github.com/gin-gonic/gin"
)

type FacebookWebhookPayload struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Field string `json:"field"`
			Value struct {
				From struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"from"`
				Message   string `json:"message"`
				CommentID string `json:"comment_id"`
				PostID    string `json:"post_id"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

func FacebookWebhookHandler(c *gin.Context) {
	var payload FacebookWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fmt.Println("❌ Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	fmt.Println("📩 Webhook received")

	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			fmt.Printf("🧠 Field: %s\n", change.Field)
			text := change.Value.Message
			commentID := change.Value.CommentID
			fmt.Printf("🗨️ Message: %s | 💬 CommentID: %s\n", text, commentID)

			if reply, ok := utils.MatchKeyword(text); ok {
				err := utils.PostCommentReply(commentID, reply, os.Getenv("FB_PAGE_TOKEN"))
				if err != nil {
					fmt.Println("⚠️ Failed to reply:", err)
				} else {
					fmt.Println("✅ Replied to comment!")
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func FacebookVerifyHandler(c *gin.Context) {
	verifyToken := os.Getenv("FB_VERIFY_TOKEN")
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == verifyToken {
		c.String(http.StatusOK, challenge)
	} else {
		c.String(http.StatusForbidden, "Verification failed")
	}
}

func GetKeywordsHandler(c *gin.Context) {
	pairs := utils.GetAllKeywordPairs()
	c.JSON(http.StatusOK, pairs)
}

func SaveKeywordsHandler(c *gin.Context) {
	var pairs []map[string]string
	if err := c.ShouldBindJSON(&pairs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data format"})
		return
	}
	utils.SetKeywordPairs(pairs)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
