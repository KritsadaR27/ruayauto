package facebook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client represents Facebook Graph API client
type Client struct {
	pageToken  string
	appSecret  string
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Facebook API client
func NewClient(pageToken, appSecret string) *Client {
	return &Client{
		pageToken:  pageToken,
		appSecret:  appSecret,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    "https://graph.facebook.com/v18.0",
	}
}

// SendComment sends a reply to a Facebook comment
func (c *Client) SendComment(commentID, message string) (*FacebookAPIResponse, error) {
	url := fmt.Sprintf("%s/%s/comments", c.baseURL, commentID)

	payload := map[string]string{
		"message":     message,
		"access_token": c.pageToken,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var apiResp FacebookAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &apiResp, fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	return &apiResp, nil
}

// SendMessage sends a direct message to a Facebook user
func (c *Client) SendMessage(recipientID, message string) (*FacebookAPIResponse, error) {
	url := fmt.Sprintf("%s/me/messages", c.baseURL)

	payload := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientID,
		},
		"message": map[string]string{
			"text": message,
		},
		"access_token": c.pageToken,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var apiResp FacebookAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &apiResp, fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	return &apiResp, nil
}

// SendCommentWithImage sends a reply to a Facebook comment with an image
func (c *Client) SendCommentWithImage(commentID, message, imageURL string) (*FacebookAPIResponse, error) {
	url := fmt.Sprintf("%s/%s/comments", c.baseURL, commentID)

	payload := map[string]interface{}{
		"message":      message,
		"attachment_url": imageURL,
		"access_token": c.pageToken,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var apiResp FacebookAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &apiResp, fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	return &apiResp, nil
}

// SendMessageWithImage sends a direct message to a Facebook user with an image
func (c *Client) SendMessageWithImage(recipientID, message, imageURL string) (*FacebookAPIResponse, error) {
	url := fmt.Sprintf("%s/me/messages", c.baseURL)

	payload := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientID,
		},
		"message": map[string]interface{}{
			"text": message,
			"attachment": map[string]interface{}{
				"type": "image",
				"payload": map[string]string{
					"url": imageURL,
				},
			},
		},
		"access_token": c.pageToken,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var apiResp FacebookAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &apiResp, fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	return &apiResp, nil
}
