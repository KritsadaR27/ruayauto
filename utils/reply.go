package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	keywordPairs = []struct {
		Keywords  []string
		Responses []string
	}{
		{Keywords: []string{"hello"}, Responses: []string{"Hi there!"}},
		{Keywords: []string{"help"}, Responses: []string{"How can I assist you?"}},
		{Keywords: []string{"thanks"}, Responses: []string{"You're welcome!"}},
		{Keywords: []string{"goodbye"}, Responses: []string{"Goodbye! Have a great day!"}},
		{Keywords: []string{"ทดสอบ"}, Responses: []string{"ทดสอบสำเร็จ!"}},
	}
	keywordsMu sync.RWMutex

	// สำหรับตอบกลับแบบไม่มี keyword
	DefaultResponses = &[]string{"ขอบคุณที่แสดงความคิดเห็นค่ะ"}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func MatchKeyword(text string) (string, bool) {
	// แปลงข้อความให้เป็น lower-case ก่อนเปรียบเทียบ
	lowerText := strings.ToLower(text)
	keywordsMu.RLock()
	defer keywordsMu.RUnlock()
	for _, pair := range keywordPairs {
		for _, k := range pair.Keywords {
			if strings.Contains(lowerText, strings.ToLower(k)) {
				// สุ่ม response ถ้ามีหลายอัน
				if len(pair.Responses) > 0 {
					idx := 0
					if len(pair.Responses) > 1 {
						idx = rand.Intn(len(pair.Responses))
					}
					return pair.Responses[idx], true
				}
			}
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

func GetAllKeywordPairs() []map[string]interface{} {
	keywordsMu.RLock()
	defer keywordsMu.RUnlock()
	pairs := make([]map[string]interface{}, 0, len(keywordPairs))
	for _, p := range keywordPairs {
		pairs = append(pairs, map[string]interface{}{
			"keywords":  append([]string{}, p.Keywords...),
			"responses": append([]string{}, p.Responses...),
		})
	}
	return pairs
}

func SetKeywordPairsV2(pairs []struct {
	Keywords  []string `json:"keywords"`
	Responses []string `json:"responses"`
}) {
	keywordsMu.Lock()
	defer keywordsMu.Unlock()
	keywordPairs = make([]struct {
		Keywords  []string
		Responses []string
	}, 0, len(pairs))
	for _, p := range pairs {
		if len(p.Keywords) > 0 && len(p.Responses) > 0 {
			keywordPairs = append(keywordPairs, struct {
				Keywords  []string
				Responses []string
			}{
				Keywords:  append([]string{}, p.Keywords...),
				Responses: append([]string{}, p.Responses...),
			})
		}
	}
}

func GetRandomDefaultResponse() string {
	resp := *DefaultResponses
	if len(resp) == 0 {
		return ""
	}
	if len(resp) == 1 {
		return resp[0]
	}
	return resp[rand.Intn(len(resp))]
}
