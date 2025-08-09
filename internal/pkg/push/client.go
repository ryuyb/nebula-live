package push

import (
	"context"
	"fmt"
	"time"

	"resty.dev/v3"
)

// Client provides a unified interface for push notification providers
type Client struct {
	providers  map[string]Provider
	httpClient *resty.Client
}

// ClientConfig holds the configuration for all push providers
type ClientConfig struct {
	Bark BarkConfig `mapstructure:"bark"`
}

// NewClient creates a new push notification client
func NewClient(config ClientConfig) *Client {
	httpClient := resty.New()
	httpClient.SetTimeout(30 * time.Second)
	httpClient.SetRetryCount(3)
	httpClient.SetRetryWaitTime(1 * time.Second)

	client := &Client{
		providers:  make(map[string]Provider),
		httpClient: httpClient,
	}

	// Register providers
	client.RegisterProvider(NewBarkProvider(httpClient, config.Bark))

	return client
}

// RegisterProvider registers a new push notification provider
func (c *Client) RegisterProvider(provider Provider) {
	c.providers[provider.GetProviderName()] = provider
}

// SendMessage sends a push notification via the specified provider
func (c *Client) SendMessage(ctx context.Context, providerName string, message *PushMessage) (*PushResponse, error) {
	provider, exists := c.providers[providerName]
	if !exists {
		return nil, ErrProviderNotFound
	}

	if !provider.IsEnabled() {
		return nil, ErrProviderNotEnabled
	}

	return provider.SendMessage(ctx, message)
}

// SendToAll sends a push notification to all enabled providers
func (c *Client) SendToAll(ctx context.Context, message *PushMessage) ([]*PushResponse, error) {
	var responses []*PushResponse
	var lastError error

	for _, provider := range c.providers {
		if !provider.IsEnabled() {
			continue
		}

		// Create a copy of the message for each provider
		msgCopy := *message
		resp, err := provider.SendMessage(ctx, &msgCopy)
		if err != nil {
			lastError = err
			// Create error response if provider returned an error
			if resp == nil {
				resp = &PushResponse{
					Success:  false,
					Error:    err.Error(),
					Provider: provider.GetProviderName(),
				}
			}
		}

		if resp != nil {
			responses = append(responses, resp)
		}
	}

	// Return error only if no providers succeeded
	if len(responses) == 0 && lastError != nil {
		return nil, fmt.Errorf("all providers failed: %w", lastError)
	}

	return responses, nil
}

// GetSupportedProviders returns a list of supported providers
func (c *Client) GetSupportedProviders() []string {
	providers := make([]string, 0, len(c.providers))
	for name := range c.providers {
		providers = append(providers, name)
	}
	return providers
}

// GetEnabledProviders returns a list of enabled providers
func (c *Client) GetEnabledProviders() []string {
	var providers []string
	for name, provider := range c.providers {
		if provider.IsEnabled() {
			providers = append(providers, name)
		}
	}
	return providers
}

// IsProviderEnabled checks if a specific provider is enabled
func (c *Client) IsProviderEnabled(providerName string) bool {
	provider, exists := c.providers[providerName]
	if !exists {
		return false
	}
	return provider.IsEnabled()
}
