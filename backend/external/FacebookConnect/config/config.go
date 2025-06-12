package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds all configuration for FacebookConnect service
type Config struct {
	// Server configuration
	Port     string `mapstructure:"PORT"`
	GinMode  string `mapstructure:"GIN_MODE"`
	
	// Facebook configuration
	VerifyToken string `mapstructure:"FB_VERIFY_TOKEN"`
	PageToken   string `mapstructure:"FB_PAGE_TOKEN"`
	
	// ChatBotCore service configuration
	ChatBotCoreURL string `mapstructure:"CHATBOT_CORE_URL"`
	
	// Logging configuration
	LogLevel string `mapstructure:"LOG_LEVEL"`
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	viper.SetDefault("PORT", "8085")
	viper.SetDefault("GIN_MODE", "release")
	viper.SetDefault("CHATBOT_CORE_URL", "http://chatbot-core:8086")
	viper.SetDefault("LOG_LEVEL", "info")
	
	// Bind environment variables
	viper.AutomaticEnv()
	
	// Try to read from config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	
	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found is OK, we'll use env vars and defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}
	
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	// Validate required configuration
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return &cfg, nil
}

// validate checks if all required configuration is present
func (c *Config) validate() error {
	if c.VerifyToken == "" {
		return fmt.Errorf("FB_VERIFY_TOKEN is required")
	}
	
	if c.PageToken == "" {
		return fmt.Errorf("FB_PAGE_TOKEN is required")
	}
	
	if c.ChatBotCoreURL == "" {
		return fmt.Errorf("CHATBOT_CORE_URL is required")
	}
	
	return nil
}

// GetPort returns the port with colon prefix for server binding
func (c *Config) GetPort() string {
	if c.Port[0] != ':' {
		return ":" + c.Port
	}
	return c.Port
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.GinMode == "debug" || os.Getenv("APP_ENV") == "development"
}
