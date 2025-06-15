package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"facebook-connect/config"
)

// FacebookAPIService handles Facebook Graph API interactions
type FacebookAPIService struct {
	config *config.Config
	client *http.Client
}

// NewFacebookAPIService creates a new Facebook API service
func NewFacebookAPIService(cfg *config.Config) *FacebookAPIService {
	return &FacebookAPIService{
		config: cfg,
		client: &http.Client{},
	}
}

// SendMessage sends a message to a Facebook user
func (s *FacebookAPIService) SendMessage(recipientID, message string) error {
	url := fmt.Sprintf("%s/me/messages?access_token=%s", s.config.GraphAPIURL, s.config.PageAccessToken)

	payload := map[string]interface{}{
		"recipient": map[string]string{"id": recipientID},
		"message":   map[string]string{"text": message},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message payload: %w", err)
	}

	resp, err := s.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send Facebook message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Facebook API returned status: %d", resp.StatusCode)
	}

	log.Printf("✅ Message sent to user %s: %s", recipientID, message)
	return nil
}
