package model

import "time"

// Platform represents different social media platforms
type Platform string

const (
	PlatformFacebook  Platform = "facebook"
	PlatformLine      Platform = "line"
	PlatformTikTok    Platform = "tiktok"
	PlatformInstagram Platform = "instagram"
	PlatformTwitter   Platform = "twitter"
)

// UnifiedMessage represents a standardized message across all platforms
type UnifiedMessage struct {
	ID          string            `json:"id"`
	Platform    Platform          `json:"platform"`
	UserID      string            `json:"user_id"`
	PageID      string            `json:"page_id"`
	Content     string            `json:"content"`
	MessageType MessageType       `json:"message_type"`
	Timestamp   time.Time         `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	
	// Platform-specific IDs
	CommentID   string `json:"comment_id,omitempty"`   // Facebook
	PostID      string `json:"post_id,omitempty"`      // Facebook/Instagram
	ReplyToken  string `json:"reply_token,omitempty"`  // LINE
	TweetID     string `json:"tweet_id,omitempty"`     // Twitter
	VideoID     string `json:"video_id,omitempty"`     // TikTok
}

// MessageType represents different types of messages
type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeVideo    MessageType = "video"
	MessageTypeSticker  MessageType = "sticker"
	MessageTypeComment  MessageType = "comment"
	MessageTypeReply    MessageType = "reply"
)

// WebhookEvent represents incoming webhook events
type WebhookEvent struct {
	Platform    Platform    `json:"platform"`
	EventType   string      `json:"event_type"`
	Timestamp   time.Time   `json:"timestamp"`
	RawData     interface{} `json:"raw_data"`
	Signature   string      `json:"signature,omitempty"`
}

// ChatbotRequest represents request to chatbot service
type ChatbotRequest struct {
	Platform  Platform `json:"platform"`
	UserID    string   `json:"user_id"`
	PageID    string   `json:"page_id"`
	Content   string   `json:"content"`
	CommentID string   `json:"comment_id,omitempty"`
	PostID    string   `json:"post_id,omitempty"`
}

// ChatbotResponse represents response from chatbot service
type ChatbotResponse struct {
	ShouldReply      bool    `json:"should_reply"`
	Response         string  `json:"response,omitempty"`
	HasMedia         bool    `json:"has_media,omitempty"`
	MediaDescription *string `json:"media_description,omitempty"`
	MatchedKeyword   string  `json:"matched_keyword,omitempty"`
	MatchType        string  `json:"match_type,omitempty"`
}

// PlatformConfig represents configuration for each platform
type PlatformConfig struct {
	Enabled      bool              `json:"enabled"`
	Credentials  map[string]string `json:"credentials"`
	WebhookPath  string            `json:"webhook_path"`
	VerifyToken  string            `json:"verify_token,omitempty"`
}

// WebhookConfig represents overall webhook configuration
type WebhookConfig struct {
	Port      string                    `json:"port"`
	Host      string                    `json:"host"`
	Platforms map[Platform]PlatformConfig `json:"platforms"`
	ChatbotURL string                   `json:"chatbot_url"`
}
