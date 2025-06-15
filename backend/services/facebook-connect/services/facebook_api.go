package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"facebook-connect/models"
)

// FacebookAPI handles communication with Facebook Graph API
type FacebookAPI struct {
	pageToken string
	baseURL   string
	client    *http.Client
}

// NewFacebookAPI creates a new Facebook API client
func NewFacebookAPI(pageToken string) *FacebookAPI {
	return &FacebookAPI{
		pageToken: pageToken,
		baseURL:   "https://graph.facebook.com/v18.0",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ReplyToComment posts a reply to a Facebook comment
func (f *FacebookAPI) ReplyToComment(ctx context.Context, commentID, message string) (*models.FacebookAPIResponse, error) {
	log.Printf("📤 Replying to comment %s: %s", commentID, message)

	// Prepare request data
	data := url.Values{}
	data.Set("message", message)
	data.Set("access_token", f.pageToken)

	// Create request URL
	requestURL := fmt.Sprintf("%s/%s/comments", f.baseURL, commentID)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var apiResponse models.FacebookAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ Facebook API error: %+v", apiResponse.Error)
		return &apiResponse, fmt.Errorf("Facebook API error: %s", apiResponse.Error.Message)
	}

	log.Printf("✅ Successfully replied to comment %s", commentID)
	return &apiResponse, nil
}
