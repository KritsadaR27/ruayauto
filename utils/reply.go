package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var (
	keywords = map[string]string{
		"hello":   "Hi there!",
		"help":    "How can I assist you?",
		"thanks":  "You're welcome!",
		"goodbye": "Goodbye! Have a great day!",
		"ทดสอบ":   "ทดสอบสำเร็จ!",
	}
	keywordsMu sync.RWMutex
)

func MatchKeyword(text string) (string, bool) {
	// แปลงข้อความให้เป็น lower-case ก่อนเปรียบเทียบ
	lowerText := strings.ToLower(text)
	for k, v := range keywords {
		if strings.Contains(lowerText, strings.ToLower(k)) {
			return v, true
		}
	}
	return "", false
}

func PostCommentReply(commentID, reply, pageToken string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v19.0/%s/comments", commentID)
	body := map[string]string{
		"message":      reply,
		"access_token": pageToken,
	}
	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// เพิ่มการเช็ก response
	if resp.StatusCode >= 300 {
		return fmt.Errorf("graph API error: status %d", resp.StatusCode)
	}

	return nil
}

func GetAllKeywordPairs() []map[string]string {
	keywordsMu.RLock()
	defer keywordsMu.RUnlock()
	pairs := make([]map[string]string, 0, len(keywords))
	for k, v := range keywords {
		pairs = append(pairs, map[string]string{"keyword": k, "response": v})
	}
	return pairs
}

func SetKeywordPairs(pairs []map[string]string) {
	keywordsMu.Lock()
	defer keywordsMu.Unlock()
	keywords = make(map[string]string)
	for _, p := range pairs {
		k, ok1 := p["keyword"]
		r, ok2 := p["response"]
		if ok1 && ok2 {
			keywords[k] = r
		}
	}
}
