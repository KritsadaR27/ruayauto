package facebook

import (
	"time"

	"webhook/internal/model"
)

// MapToUnifiedMessage converts Facebook webhook data to UnifiedMessage
func MapToUnifiedMessage(entry FacebookEntry, change FacebookChange) *model.UnifiedMessage {
	value := change.Value

	// Create base unified message
	unified := &model.UnifiedMessage{
		Platform:    model.PlatformFacebook,
		UserID:      value.From.ID,
		PageID:      entry.ID,
		Content:     value.Message,
		MessageType: getMessageType(change.Field, value.Verb),
		Timestamp:   time.Unix(value.CreatedTime, 0),
		Metadata: map[string]interface{}{
			"verb":      value.Verb,
			"field":     change.Field,
			"user_name": value.From.Name,
		},
	}

	// Set platform-specific IDs based on the type of event
	if change.Field == "comments" {
		unified.CommentID = value.CommentID
		unified.PostID = value.PostID
		unified.ID = value.CommentID

		// Handle comment replies
		if value.ParentID != "" {
			unified.Metadata["parent_id"] = value.ParentID
			unified.MessageType = model.MessageTypeReply
		}
	} else if change.Field == "posts" {
		unified.PostID = value.PostID
		unified.ID = value.PostID
	}

	return unified
}

// MapMessagingToUnifiedMessage converts Facebook messaging data to UnifiedMessage
func MapMessagingToUnifiedMessage(entry FacebookEntry, messaging FacebookMessaging) *model.UnifiedMessage {
	return &model.UnifiedMessage{
		ID:          messaging.Message.MID,
		Platform:    model.PlatformFacebook,
		UserID:      messaging.Sender.ID,
		PageID:      entry.ID,
		Content:     messaging.Message.Text,
		MessageType: model.MessageTypeText,
		Timestamp:   time.Unix(messaging.Timestamp/1000, 0), // Facebook uses milliseconds
		Metadata: map[string]interface{}{
			"recipient_id": messaging.Recipient.ID,
			"message_type": "direct_message",
		},
	}
}

// getMessageType determines the message type based on Facebook webhook data
func getMessageType(field, verb string) model.MessageType {
	switch field {
	case "comments":
		if verb == "add" {
			return model.MessageTypeComment
		}
		return model.MessageTypeText
	case "posts":
		return model.MessageTypeText
	default:
		return model.MessageTypeText
	}
}

// MapFromUnifiedResponse maps chatbot response back to Facebook-specific format
func MapFromUnifiedResponse(response *model.ChatbotResponse, originalMsg *model.UnifiedMessage) map[string]interface{} {
	result := map[string]interface{}{
		"platform":     "facebook",
		"should_reply": response.ShouldReply,
	}

	if response.ShouldReply {
		result["response"] = response.Response
		result["has_media"] = response.HasMedia
		if response.MediaDescription != nil {
			result["media_description"] = *response.MediaDescription
		}
		result["matched_keyword"] = response.MatchedKeyword

		// Set target based on message type
		if originalMsg.CommentID != "" {
			result["target_type"] = "comment"
			result["target_id"] = originalMsg.CommentID
		} else {
			result["target_type"] = "message"
			result["target_id"] = originalMsg.UserID
		}
	}

	return result
}
