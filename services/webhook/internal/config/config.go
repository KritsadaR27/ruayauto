package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Facebook FacebookConfig `json:"facebook"`
	Line     LineConfig     `json:"line"`
	Chatbot  ChatbotConfig  `json:"chatbot"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

// FacebookConfig holds Facebook-specific configuration
type FacebookConfig struct {
	Enabled     bool   `json:"enabled"`
	VerifyToken string `json:"verify_token"`
	PageToken   string `json:"page_token"`
	AppSecret   string `json:"app_secret"`
}

// LineConfig holds LINE-specific configuration
type LineConfig struct {
	Enabled       bool   `json:"enabled"`
	ChannelSecret string `json:"channel_secret"`
	ChannelToken  string `json:"channel_token"`
}

// ChatbotConfig holds chatbot service configuration
type ChatbotConfig struct {
	URL string `json:"url"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("WEBHOOK_PORT", "8091"),
			Host: getEnv("WEBHOOK_HOST", "0.0.0.0"),
		},
		Facebook: FacebookConfig{
			Enabled:     getEnvBool("FACEBOOK_ENABLED", true),
			VerifyToken: getEnv("FACEBOOK_VERIFY_TOKEN", ""),
			PageToken:   getEnv("FACEBOOK_PAGE_TOKEN", ""),
			AppSecret:   getEnv("FACEBOOK_APP_SECRET", ""),
		},
		Line: LineConfig{
			Enabled:       getEnvBool("LINE_ENABLED", true),
			ChannelSecret: getEnv("LINE_CHANNEL_SECRET", ""),
			ChannelToken:  getEnv("LINE_CHANNEL_TOKEN", ""),
		},
		Chatbot: ChatbotConfig{
			URL: getEnv("CHATBOT_URL", "http://localhost:8090"),
		},
	}

	// Validate required configurations
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// validate validates the configuration
func (c *Config) validate() error {
	if c.Facebook.Enabled {
		if c.Facebook.VerifyToken == "" {
			return fmt.Errorf("FACEBOOK_VERIFY_TOKEN is required when Facebook is enabled")
		}
		if c.Facebook.PageToken == "" {
			return fmt.Errorf("FACEBOOK_PAGE_TOKEN is required when Facebook is enabled")
		}
		if c.Facebook.AppSecret == "" {
			return fmt.Errorf("FACEBOOK_APP_SECRET is required when Facebook is enabled")
		}
	}

	if c.Line.Enabled {
		if c.Line.ChannelSecret == "" {
			return fmt.Errorf("LINE_CHANNEL_SECRET is required when LINE is enabled")
		}
		if c.Line.ChannelToken == "" {
			return fmt.Errorf("LINE_CHANNEL_TOKEN is required when LINE is enabled")
		}
	}

	if c.Chatbot.URL == "" {
		return fmt.Errorf("CHATBOT_URL is required")
	}

	return nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool gets a boolean environment variable with a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
