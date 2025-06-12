package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"ruaymanagement/backend/external/FacebookConnect/models"
)

// FacebookAPIClient handles communication with Facebook Graph API
type FacebookAPIClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewFacebookAPIClient creates a new Facebook API client
func NewFacebookAPIClient() *FacebookAPIClient {
	return &FacebookAPIClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://graph.facebook.com/v18.0",
	}
}

// SendCommentReply sends a reply to a Facebook comment
func (c *FacebookAPIClient) SendCommentReply(ctx context.Context, req models.FacebookReplyRequest) error {
	url := fmt.Sprintf("%s/%s/comments", c.baseURL, req.CommentID)

	payload := map[string]string{
		"message":      req.Message,
		"access_token": req.PageToken,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Facebook API error: status %d", resp.StatusCode)
	}

	return nil
}

// VerifyWebhook verifies Facebook webhook verification requests
func (c *FacebookAPIClient) VerifyWebhook(mode, token, challenge, verifyToken string) (string, error) {
	if mode != "subscribe" {
		return "", fmt.Errorf("invalid mode: %s", mode)
	}

	if token != verifyToken {
		return "", fmt.Errorf("invalid verification token")
	}

	return challenge, nil
}
