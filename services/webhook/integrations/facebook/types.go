package facebook

import "time"

// FacebookWebhookPayload represents the main webhook payload from Facebook
type FacebookWebhookPayload struct {
	Object string          `json:"object"`
	Entry  []FacebookEntry `json:"entry"`
}

// FacebookEntry represents each entry in the webhook
type FacebookEntry struct {
	ID      string                `json:"id"`
	Time    int64                 `json:"time"`
	Changes []FacebookChange      `json:"changes,omitempty"`
	Messaging []FacebookMessaging `json:"messaging,omitempty"`
}

// FacebookChange represents changes in posts/comments
type FacebookChange struct {
	Field string                 `json:"field"`
	Value FacebookChangeValue    `json:"value"`
}

// FacebookChangeValue represents the actual change data
type FacebookChangeValue struct {
	From         FacebookUser    `json:"from"`
	Post         FacebookPost    `json:"post,omitempty"`
	Comment      FacebookComment `json:"comment,omitempty"`
	CommentID    string          `json:"comment_id,omitempty"`
	PostID       string          `json:"post_id,omitempty"`
	ParentID     string          `json:"parent_id,omitempty"`
	CreatedTime  int64           `json:"created_time"`
	Message      string          `json:"message"`
	Verb         string          `json:"verb"` // add, edit, remove
}

// FacebookMessaging represents direct messages
type FacebookMessaging struct {
	Sender    FacebookUser    `json:"sender"`
	Recipient FacebookUser    `json:"recipient"`
	Timestamp int64           `json:"timestamp"`
	Message   FacebookMessage `json:"message,omitempty"`
}

// FacebookUser represents a Facebook user
type FacebookUser struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// FacebookPost represents a Facebook post
type FacebookPost struct {
	StatusType    string    `json:"status_type"`
	IsPublished   bool      `json:"is_published"`
	UpdatedTime   time.Time `json:"updated_time"`
	PermalinkURL  string    `json:"permalink_url"`
	PromotionStatus string  `json:"promotion_status"`
	ID            string    `json:"id"`
}

// FacebookComment represents a Facebook comment
type FacebookComment struct {
	ID          string    `json:"id"`
	Message     string    `json:"message"`
	CreatedTime time.Time `json:"created_time"`
	From        FacebookUser `json:"from"`
	ParentID    string    `json:"parent_id,omitempty"`
}

// FacebookMessage represents a direct message
type FacebookMessage struct {
	MID  string `json:"mid"`
	Text string `json:"text"`
}

// FacebookAPIResponse represents response from Facebook Graph API
type FacebookAPIResponse struct {
	Success bool   `json:"success,omitempty"`
	Error   *FacebookError `json:"error,omitempty"`
	ID      string `json:"id,omitempty"`
}

// FacebookError represents an error from Facebook API
type FacebookError struct {
	Message      string `json:"message"`
	Type         string `json:"type"`
	Code         int    `json:"code"`
	ErrorSubcode int    `json:"error_subcode,omitempty"`
	FBTraceID    string `json:"fbtrace_id,omitempty"`
}
