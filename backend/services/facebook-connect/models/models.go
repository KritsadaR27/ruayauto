package models

// FacebookWebhookEntry represents a Facebook webhook entry
type FacebookWebhookEntry struct {
	ID      string                  `json:"id"`
	Time    int64                   `json:"time"`
	Changes []FacebookWebhookChange `json:"changes,omitempty"`
}

// FacebookWebhookChange represents a change in Facebook webhook
type FacebookWebhookChange struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

// FacebookWebhookPayload represents the full webhook payload from Facebook
type FacebookWebhookPayload struct {
	Object string                 `json:"object"`
	Entry  []FacebookWebhookEntry `json:"entry"`
}

// CommentData represents Facebook comment data
type CommentData struct {
	CommentID   string `json:"comment_id"`
	PostID      string `json:"post_id"`
	PageID      string `json:"page_id"`
	UserID      string `json:"user_id"`
	Message     string `json:"message"`
	CreatedTime string `json:"created_time"`
}

// ChatBotRequest represents request to ChatBot service
type ChatBotRequest struct {
	PageID    string `json:"page_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CommentID string `json:"comment_id,omitempty"`
	PostID    string `json:"post_id,omitempty"`
}

// ChatBotResponse represents response from ChatBot service
type ChatBotResponse struct {
	ShouldReply    bool   `json:"should_reply"`
	Response       string `json:"response,omitempty"`
	MatchedKeyword string `json:"matched_keyword,omitempty"`
	MatchType      string `json:"match_type,omitempty"`
}

// FacebookAPIResponse represents Facebook API response
type FacebookAPIResponse struct {
	Success bool              `json:"success"`
	ID      string            `json:"id,omitempty"`
	Error   *FacebookAPIError `json:"error,omitempty"`
}

// FacebookAPIError represents Facebook API error
type FacebookAPIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    int    `json:"code"`
	Subcode int    `json:"error_subcode,omitempty"`
}
