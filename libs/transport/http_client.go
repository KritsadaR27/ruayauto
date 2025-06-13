package transport

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
)

// HTTPClient wraps resty client with common functionality
type HTTPClient struct {
	client      *resty.Client
	rateLimiter *rate.Limiter
	baseURL     string
	timeout     time.Duration
}

// HTTPClientConfig represents HTTP client configuration
type HTTPClientConfig struct {
	BaseURL        string
	Timeout        time.Duration
	RateLimitRPS   int // requests per second
	RateLimitBurst int // burst size
	RetryAttempts  int
	RetryDelay     time.Duration
	Headers        map[string]string
}

// NewHTTPClient creates a new HTTP client with rate limiting
func NewHTTPClient(config HTTPClientConfig) *HTTPClient {
	client := resty.New()

	// Configure timeout
	client.SetTimeout(config.Timeout)

	// Configure retries
	client.SetRetryCount(config.RetryAttempts)
	client.SetRetryWaitTime(config.RetryDelay)

	// Set default headers
	for key, value := range config.Headers {
		client.SetHeader(key, value)
	}

	// Create rate limiter
	limiter := rate.NewLimiter(rate.Limit(config.RateLimitRPS), config.RateLimitBurst)

	return &HTTPClient{
		client:      client,
		rateLimiter: limiter,
		baseURL:     config.BaseURL,
		timeout:     config.Timeout,
	}
}

// GET performs a GET request with rate limiting
func (h *HTTPClient) GET(ctx context.Context, endpoint string) (*resty.Response, error) {
	// Wait for rate limiter
	if err := h.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	resp, err := h.client.R().
		SetContext(ctx).
		Get(h.baseURL + endpoint)

	return resp, err
}

// POST performs a POST request with rate limiting
func (h *HTTPClient) POST(ctx context.Context, endpoint string, body interface{}) (*resty.Response, error) {
	// Wait for rate limiter
	if err := h.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	resp, err := h.client.R().
		SetContext(ctx).
		SetBody(body).
		Post(h.baseURL + endpoint)

	return resp, err
}

// SetAuthToken sets authorization token for requests
func (h *HTTPClient) SetAuthToken(token string) {
	h.client.SetAuthToken(token)
}

// SetHeader sets a custom header for all requests
func (h *HTTPClient) SetHeader(key, value string) {
	h.client.SetHeader(key, value)
}
