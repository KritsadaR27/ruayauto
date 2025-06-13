package line

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client represents LINE Messaging API client
type Client struct {
	channelToken string
	httpClient   *http.Client
	baseURL      string
}

// NewClient creates a new LINE API client
func NewClient(channelToken string) *Client {
	return &Client{
		channelToken: channelToken,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		baseURL:      "https://api.line.me/v2/bot",
	}
}

// ReplyMessage sends a reply message to LINE
func (c *Client) ReplyMessage(replyToken, message string) (*LineAPIResponse, error) {
	url := fmt.Sprintf("%s/message/reply", c.baseURL)
	
	payload := LineReplyMessage{
		ReplyToken: replyToken,
		Messages: []LineTextMessage{
			{
				Type: "text",
				Text: message,
			},
		},
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.channelToken)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	var apiResp LineAPIResponse
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &apiResp, nil
}

// PushMessage sends a push message to LINE user
func (c *Client) PushMessage(userID, message string) (*LineAPIResponse, error) {
	url := fmt.Sprintf("%s/message/push", c.baseURL)
	
	payload := LinePushMessage{
		To: userID,
		Messages: []LineTextMessage{
			{
				Type: "text",
				Text: message,
			},
		},
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.channelToken)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	var apiResp LineAPIResponse
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &apiResp, nil
}

// GetProfile gets LINE user profile
func (c *Client) GetProfile(userID string) (*LineProfile, error) {
	url := fmt.Sprintf("%s/profile/%s", c.baseURL, userID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+c.channelToken)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}
	
	var profile LineProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &profile, nil
}
