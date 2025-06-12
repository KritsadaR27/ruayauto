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

// ChatBotCoreClient handles communication with the internal ChatBotCore service
type ChatBotCoreClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewChatBotCoreClient creates a new ChatBotCore client
func NewChatBotCoreClient(baseURL string) *ChatBotCoreClient {
	return &ChatBotCoreClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

// ProcessMessage sends a message to ChatBotCore for processing
func (c *ChatBotCoreClient) ProcessMessage(ctx context.Context, message *models.InternalMessage) (*models.InternalResponse, error) {
	url := fmt.Sprintf("%s/api/messages/process", c.baseURL)

	jsonData, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to ChatBotCore: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ChatBotCore error: status %d", resp.StatusCode)
	}

	var response models.InternalResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// HealthCheck checks if ChatBotCore is healthy
func (c *ChatBotCoreClient) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/health", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform health check: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ChatBotCore unhealthy: status %d", resp.StatusCode)
	}

	return nil
}
