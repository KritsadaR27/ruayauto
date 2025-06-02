package handler

import (
	"fmt"
	"net/http"
	"os"
	"ruayAutoMsg/utils"
	"strings"

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

var (
	defaultResponses = []string{"ขอบคุณที่แสดงความคิดเห็นค่ะ"}
	enableDefault    = true
	noTag            = true
	noSticker        = true
)

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

			// เช็ค noTag
			if noTag && strings.Contains(text, "@") {
				fmt.Println("🚫 ไม่ตอบคอมเมนต์ที่มีการแท็ก")
				continue
			}
			// TODO: เช็ค sticker (ถ้ามี field ใน payload)

			if reply, ok := utils.MatchKeyword(text); ok {
				err := utils.PostCommentReply(commentID, reply, os.Getenv("FB_PAGE_TOKEN"))
				if err != nil {
					fmt.Println("⚠️ Failed to reply:", err)
				} else {
					fmt.Println("✅ Replied to comment!")
				}
			} else if enableDefault && len(defaultResponses) > 0 {
				// ตอบแบบ default
				reply := utils.GetRandomDefaultResponse()
				if reply != "" {
					err := utils.PostCommentReply(commentID, reply, os.Getenv("FB_PAGE_TOKEN"))
					if err != nil {
						fmt.Println("⚠️ Failed to reply (default):", err)
					} else {
						fmt.Println("✅ Replied with default response!")
					}
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
	c.JSON(http.StatusOK, gin.H{
		"pairs": pairs,
		"defaultResponses": defaultResponses,
		"enableDefault": enableDefault,
		"noTag": noTag,
		"noSticker": noSticker,
	})
}

func SaveKeywordsHandler(c *gin.Context) {
	var req struct {
		Pairs           []struct {
			Keywords  []string `json:"keywords"`
			Responses []string `json:"responses"`
		} `json:"pairs"`
		DefaultResponses []string `json:"defaultResponses"`
		EnableDefault   bool     `json:"enableDefault"`
		NoTag           bool     `json:"noTag"`
		NoSticker       bool     `json:"noSticker"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data format"})
		return
	}
	utils.SetKeywordPairsV2(req.Pairs)
	defaultResponses = req.DefaultResponses
	enableDefault = req.EnableDefault
	noTag = req.NoTag
	noSticker = req.NoSticker
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
