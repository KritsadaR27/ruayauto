package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the FacebookConnect service
type Config struct {
	Port           string
	Host           string
	VerifyToken    string
	PageToken      string
	AppSecret      string
	ChatBotCoreURL string
	LogLevel       string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Port:           getEnv("PORT", "8091"),
		Host:           getEnv("HOST", "0.0.0.0"),
		VerifyToken:    getEnv("FB_VERIFY_TOKEN", ""),
		PageToken:      getEnv("FB_PAGE_TOKEN", ""),
		AppSecret:      getEnv("FB_APP_SECRET", ""),
		ChatBotCoreURL: getEnv("CHATBOT_CORE_URL", "http://localhost:8090"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}

	// Validate required fields
	if cfg.VerifyToken == "" {
		return nil, fmt.Errorf("FB_VERIFY_TOKEN is required")
	}
	if cfg.PageToken == "" {
		return nil, fmt.Errorf("FB_PAGE_TOKEN is required")
	}

	return cfg, nil
}

// GetServerAddr returns server address
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// GetChatBotURL returns the ChatBot service URL with endpoint
func (c *Config) GetChatBotURL(endpoint string) string {
	return fmt.Sprintf("%s/api/%s", c.ChatBotCoreURL, endpoint)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
