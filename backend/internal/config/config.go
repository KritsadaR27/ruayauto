package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration values for the application
type Config struct {
	// Application Configuration
	Port    string `mapstructure:"PORT"`
	GinMode string `mapstructure:"GIN_MODE"`

	// Database Configuration
	DatabaseURL  string `mapstructure:"DATABASE_URL"`
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       string `mapstructure:"DB_PORT"`
	DBUser       string `mapstructure:"DB_USER"`
	DBPassword   string `mapstructure:"DB_PASSWORD"`
	DBName       string `mapstructure:"DB_NAME"`
	DBSSLMode    string `mapstructure:"DB_SSLMODE"`

	// Facebook Configuration
	FacebookAppSecret     string `mapstructure:"FACEBOOK_APP_SECRET"`
	FacebookVerifyToken   string `mapstructure:"FACEBOOK_VERIFY_TOKEN"`
	FacebookPageToken     string `mapstructure:"FACEBOOK_PAGE_ACCESS_TOKEN"`
	FBPageToken           string `mapstructure:"FB_PAGE_TOKEN"`           // Alternative naming
	FBVerifyToken         string `mapstructure:"FB_VERIFY_TOKEN"`         // Alternative naming

	// Messenger Configuration
	MessengerVerifyToken      string `mapstructure:"MESSENGER_VERIFY_TOKEN"`
	MessengerPageAccessToken  string `mapstructure:"MESSENGER_PAGE_ACCESS_TOKEN"`

	// Redis Configuration (for caching)
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	// Security Configuration
	APIKey    string `mapstructure:"API_KEY"`
	JWTSecret string `mapstructure:"JWT_SECRET"`

	// Rate Limiting
	RateLimitRequests int `mapstructure:"RATE_LIMIT_REQUESTS"`
	RateLimitWindow   int `mapstructure:"RATE_LIMIT_WINDOW"`

	// Logging
	LogLevel  string `mapstructure:"LOG_LEVEL"`
	LogFormat string `mapstructure:"LOG_FORMAT"`

	// Performance
	CacheTTL       int `mapstructure:"CACHE_TTL"`
	MaxConnections int `mapstructure:"MAX_CONNECTIONS"`

	// Features
	EnableAnalytics    bool `mapstructure:"ENABLE_ANALYTICS"`
	EnableSmartReply   bool `mapstructure:"ENABLE_SMART_REPLY"`
	EnableFuzzyMatch   bool `mapstructure:"ENABLE_FUZZY_MATCHING"`

	// External Services
	WebhookTimeout int `mapstructure:"WEBHOOK_TIMEOUT"`
	APITimeout     int `mapstructure:"API_TIMEOUT"`
}

// Load initializes and loads configuration from environment variables and .env file
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Configure Viper
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")           // Look in current directory
	v.AddConfigPath("../")         // Look in parent directory (for when running from backend/)
	v.AddConfigPath("../../")      // Look two levels up

	// Enable automatic environment variable binding
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read .env file if it exists (optional)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is OK, we'll use environment variables
		fmt.Printf("DEBUG: .env file not found, using environment variables only\n")
	} else {
		fmt.Printf("DEBUG: Successfully loaded .env file from: %s\n", v.ConfigFileUsed())
	}

	// Debug: Check what Viper can see
	fmt.Printf("DEBUG: Viper values - FACEBOOK_VERIFY_TOKEN: '%s'\n", v.GetString("FACEBOOK_VERIFY_TOKEN"))
	fmt.Printf("DEBUG: Viper values - FACEBOOK_PAGE_ACCESS_TOKEN: '%s'\n", v.GetString("FACEBOOK_PAGE_ACCESS_TOKEN"))
	fmt.Printf("DEBUG: Viper values - FB_VERIFY_TOKEN: '%s'\n", v.GetString("FB_VERIFY_TOKEN"))
	fmt.Printf("DEBUG: Viper values - FB_PAGE_TOKEN: '%s'\n", v.GetString("FB_PAGE_TOKEN"))

	// Debug logging - remove this in production
	fmt.Printf("DEBUG: Starting configuration load\n")

	// Unmarshal configuration into struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Debug logging - remove this in production  
	fmt.Printf("DEBUG: Raw config values - FacebookVerifyToken: '%s', FacebookPageToken: '%s'\n", 
		config.FacebookVerifyToken, config.FacebookPageToken)
	fmt.Printf("DEBUG: Raw config values - FBVerifyToken: '%s', FBPageToken: '%s'\n", 
		config.FBVerifyToken, config.FBPageToken)

	// Post-process configuration
	if err := postProcess(&config); err != nil {
		return nil, fmt.Errorf("error post-processing config: %w", err)
	}

	// Debug logging - remove this in production
	fmt.Printf("DEBUG: Loaded Facebook tokens - Verify: '%s', Page: '%s'\n", 
		config.GetFacebookVerifyToken(), config.GetFacebookPageToken())

	return &config, nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	// Application defaults
	v.SetDefault("PORT", "3006")
	v.SetDefault("GIN_MODE", "release")

	// Database defaults
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_USER", "ruay")
	v.SetDefault("DB_PASSWORD", "ruay1234")
	v.SetDefault("DB_NAME", "ruayAutoMsg")
	v.SetDefault("DB_SSLMODE", "disable")

	// Redis defaults
	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", "6379")
	v.SetDefault("REDIS_DB", 0)

	// Rate limiting defaults
	v.SetDefault("RATE_LIMIT_REQUESTS", 100)
	v.SetDefault("RATE_LIMIT_WINDOW", 60)

	// Logging defaults
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_FORMAT", "json")

	// Performance defaults
	v.SetDefault("CACHE_TTL", 300)
	v.SetDefault("MAX_CONNECTIONS", 100)

	// Feature defaults
	v.SetDefault("ENABLE_ANALYTICS", true)
	v.SetDefault("ENABLE_SMART_REPLY", true)
	v.SetDefault("ENABLE_FUZZY_MATCHING", true)

	// External service defaults
	v.SetDefault("WEBHOOK_TIMEOUT", 30)
	v.SetDefault("API_TIMEOUT", 10)
}

// postProcess handles any post-processing of configuration values
func postProcess(config *Config) error {
	// Build DatabaseURL if not provided but individual components are available
	if config.DatabaseURL == "" && config.DBHost != "" {
		config.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			config.DBUser,
			config.DBPassword,
			config.DBHost,
			config.DBPort,
			config.DBName,
			config.DBSSLMode,
		)
	}

	// Handle alternative Facebook token naming
	if config.FBPageToken != "" && config.FacebookPageToken == "" {
		config.FacebookPageToken = config.FBPageToken
	}
	if config.FBVerifyToken != "" && config.FacebookVerifyToken == "" {
		config.FacebookVerifyToken = config.FBVerifyToken
	}

	return nil
}

// GetDatabaseURL returns the database connection URL
func (c *Config) GetDatabaseURL() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode,
	)
}

// GetRedisURL returns the Redis connection URL
func (c *Config) GetRedisURL() string {
	if c.RedisPassword != "" {
		return fmt.Sprintf("redis://:%s@%s:%s/%d", c.RedisPassword, c.RedisHost, c.RedisPort, c.RedisDB)
	}
	return fmt.Sprintf("redis://%s:%s/%d", c.RedisHost, c.RedisPort, c.RedisDB)
}

// GetFacebookPageToken returns the Facebook page access token
func (c *Config) GetFacebookPageToken() string {
	if c.FacebookPageToken != "" {
		return c.FacebookPageToken
	}
	return c.FBPageToken
}

// GetFacebookVerifyToken returns the Facebook verify token
func (c *Config) GetFacebookVerifyToken() string {
	if c.FacebookVerifyToken != "" {
		return c.FacebookVerifyToken
	}
	return c.FBVerifyToken
}

// isPlaceholderValue checks if a value is a placeholder
func isPlaceholderValue(value string) bool {
	placeholders := []string{
		"your_", "test_", "placeholder", "change_me", "replace_me", "xxx", "yyy",
	}
	lowerValue := strings.ToLower(value)
	for _, placeholder := range placeholders {
		if strings.Contains(lowerValue, placeholder) {
			return true
		}
	}
	return false
}

// Validate checks if required configuration values are present and valid
func (c *Config) Validate() error {
	var missing []string

	// Check required Facebook configuration
	verifyToken := c.GetFacebookVerifyToken()
	if verifyToken == "" || isPlaceholderValue(verifyToken) {
		missing = append(missing, "FACEBOOK_VERIFY_TOKEN or FB_VERIFY_TOKEN")
	}
	
	pageToken := c.GetFacebookPageToken()
	if pageToken == "" || isPlaceholderValue(pageToken) {
		missing = append(missing, "FACEBOOK_PAGE_ACCESS_TOKEN or FB_PAGE_TOKEN")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required configuration: %s", strings.Join(missing, ", "))
	}

	return nil
}
