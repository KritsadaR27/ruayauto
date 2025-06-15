package line

// LineWebhookPayload represents LINE webhook payload
type LineWebhookPayload struct {
	Destination string      `json:"destination"`
	Events      []LineEvent `json:"events"`
}

// LineEvent represents a LINE webhook event
type LineEvent struct {
	Type           string       `json:"type"`
	Mode           string       `json:"mode"`
	Timestamp      int64        `json:"timestamp"`
	Source         LineSource   `json:"source"`
	WebhookEventID string       `json:"webhookEventId"`
	DeliveryMode   string       `json:"deliveryContext"`
	Message        *LineMessage `json:"message,omitempty"`
	ReplyToken     string       `json:"replyToken,omitempty"`
}

// LineSource represents the source of a LINE event
type LineSource struct {
	Type    string `json:"type"` // user, group, room
	UserID  string `json:"userId,omitempty"`
	GroupID string `json:"groupId,omitempty"`
	RoomID  string `json:"roomId,omitempty"`
}

// LineMessage represents a LINE message
type LineMessage struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Text      string `json:"text,omitempty"`
	PackageID string `json:"packageId,omitempty"`
	StickerID string `json:"stickerId,omitempty"`
}

// LineReplyMessage represents message to reply via LINE
type LineReplyMessage struct {
	ReplyToken string            `json:"replyToken"`
	Messages   []LineTextMessage `json:"messages"`
}

// LineTextMessage represents a text message for LINE
type LineTextMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// LinePushMessage represents message to push via LINE
type LinePushMessage struct {
	To       string            `json:"to"`
	Messages []LineTextMessage `json:"messages"`
}

// LineAPIResponse represents response from LINE API
type LineAPIResponse struct {
	SentMessages []struct {
		ID            string `json:"id"`
		QuotaConsumed int    `json:"quotaConsumed"`
	} `json:"sentMessages,omitempty"`
}

// LineProfile represents LINE user profile
type LineProfile struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl,omitempty"`
	StatusMessage string `json:"statusMessage,omitempty"`
}
