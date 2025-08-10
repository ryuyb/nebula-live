package livestream

import (
	"context"
	"time"

	"resty.dev/v3"
)

// Client provides a unified interface for live streaming platforms
type Client struct {
	providers  map[string]Provider
	httpClient *resty.Client
}

// NewClient creates a new livestream client
func NewClient() *Client {
	httpClient := resty.New()
	httpClient.SetTimeout(10 * time.Second)
	httpClient.SetRetryCount(3)
	httpClient.SetRetryWaitTime(1 * time.Second)

	client := &Client{
		providers:  make(map[string]Provider),
		httpClient: httpClient,
	}

	// Register default providers
	client.RegisterProvider(NewDouyuProvider(httpClient))
	client.RegisterProvider(NewBilibiliProvider(httpClient))

	return client
}

// RegisterProvider registers a new provider
func (c *Client) RegisterProvider(provider Provider) {
	c.providers[provider.GetPlatformName()] = provider
}

// GetStreamStatus gets the status of a live stream
func (c *Client) GetStreamStatus(ctx context.Context, platform, roomID string) (*StreamInfo, error) {
	provider, exists := c.providers[platform]
	if !exists {
		return nil, ErrPlatformNotFound
	}

	return provider.GetStreamStatus(ctx, roomID)
}

// GetRoomInfo gets detailed information about a live room
func (c *Client) GetRoomInfo(ctx context.Context, platform, roomID string) (*RoomInfo, error) {
	provider, exists := c.providers[platform]
	if !exists {
		return nil, ErrPlatformNotFound
	}

	return provider.GetRoomInfo(ctx, roomID)
}

// GetSupportedPlatforms returns a list of supported platforms
func (c *Client) GetSupportedPlatforms() []string {
	platforms := make([]string, 0, len(c.providers))
	for name := range c.providers {
		platforms = append(platforms, name)
	}
	return platforms
}
