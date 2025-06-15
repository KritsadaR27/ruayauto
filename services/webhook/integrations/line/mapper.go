package line

import (
	"time"

	"webhook/internal/model"
)

// MapToUnifiedMessage converts LINE webhook event to UnifiedMessage
func MapToUnifiedMessage(event LineEvent) *model.UnifiedMessage {
	unified := &model.UnifiedMessage{
		ID:          event.Message.ID,
		Platform:    model.PlatformLine,
		UserID:      event.Source.UserID,
		Content:     event.Message.Text,
		MessageType: getMessageType(event.Message.Type),
		Timestamp:   time.Unix(event.Timestamp/1000, 0), // LINE uses milliseconds
		ReplyToken:  event.ReplyToken,
		Metadata: map[string]interface{}{
			"event_type":       event.Type,
			"webhook_event_id": event.WebhookEventID,
			"source_type":      event.Source.Type,
		},
	}

	// Set PageID based on source type
	switch event.Source.Type {
	case "user":
		unified.PageID = event.Source.UserID
	case "group":
		unified.PageID = event.Source.GroupID
		unified.Metadata["group_id"] = event.Source.GroupID
	case "room":
		unified.PageID = event.Source.RoomID
		unified.Metadata["room_id"] = event.Source.RoomID
	}

	// Handle sticker messages
	if event.Message.Type == "sticker" {
		unified.Content = "[Sticker]"
		unified.Metadata["package_id"] = event.Message.PackageID
		unified.Metadata["sticker_id"] = event.Message.StickerID
	}

	return unified
}

// getMessageType determines the message type based on LINE message type
func getMessageType(lineType string) model.MessageType {
	switch lineType {
	case "text":
		return model.MessageTypeText
	case "image":
		return model.MessageTypeImage
	case "video":
		return model.MessageTypeVideo
	case "sticker":
		return model.MessageTypeSticker
	default:
		return model.MessageTypeText
	}
}

// MapFromUnifiedResponse maps chatbot response back to LINE-specific format
func MapFromUnifiedResponse(response *model.ChatbotResponse, originalMsg *model.UnifiedMessage) map[string]interface{} {
	result := map[string]interface{}{
		"platform":     "line",
		"should_reply": response.ShouldReply,
	}

	if response.ShouldReply {
		result["response"] = response.Response
		result["matched_keyword"] = response.MatchedKeyword
		result["reply_token"] = originalMsg.ReplyToken
		result["target_user"] = originalMsg.UserID
	}

	return result
}
