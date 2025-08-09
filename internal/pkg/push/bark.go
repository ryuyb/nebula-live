package push

import (
	"context"
	"fmt"

	"resty.dev/v3"
)

// Bark provider implementation
type barkProvider struct {
	client  *resty.Client
	baseURL string
	enabled bool
}

// BarkConfig holds the configuration for Bark provider
type BarkConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Enabled bool   `mapstructure:"enabled"`
}

// barkRequest represents the Bark API request payload
type barkRequest struct {
	Body      string `json:"body"`
	DeviceKey string `json:"device_key"`
	Title     string `json:"title,omitempty"`
	Subtitle  string `json:"subtitle,omitempty"`
	Badge     int    `json:"badge,omitempty"`
	Sound     string `json:"sound,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Group     string `json:"group,omitempty"`
	URL       string `json:"url,omitempty"`
	Level     string `json:"level,omitempty"`
	Call      string `json:"call,omitempty"`
	AutoCopy  string `json:"autoCopy,omitempty"`
	Copy      string `json:"copy,omitempty"`
}

// barkResponse represents the Bark API response
type barkResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// NewBarkProvider creates a new Bark provider
func NewBarkProvider(client *resty.Client, config BarkConfig) Provider {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.day.app" // Default Bark server
	}

	return &barkProvider{
		client:  client,
		baseURL: baseURL,
		enabled: config.Enabled,
	}
}

// GetProviderName returns the provider name
func (b *barkProvider) GetProviderName() string {
	return "bark"
}

// IsEnabled returns whether the provider is enabled
func (b *barkProvider) IsEnabled() bool {
	return b.enabled
}

// ValidateMessage validates the message for Bark provider
func (b *barkProvider) ValidateMessage(message *PushMessage) error {
	if message.DeviceID == "" {
		return ErrInvalidDeviceID
	}
	if message.Body == "" {
		return ErrEmptyMessage
	}
	return nil
}

// SendMessage sends a push notification via Bark
func (b *barkProvider) SendMessage(ctx context.Context, message *PushMessage) (*PushResponse, error) {
	if !b.enabled {
		return nil, ErrProviderNotEnabled
	}

	if err := b.ValidateMessage(message); err != nil {
		return nil, err
	}

	// Prepare Bark request payload
	barkReq := barkRequest{
		Body:      message.Body,
		DeviceKey: message.DeviceID,
		Title:     message.Title,
		Subtitle:  message.Subtitle,
		Badge:     message.Badge,
		Sound:     message.Sound,
		Icon:      message.Icon,
		Group:     message.Group,
		URL:       message.URL,
	}

	// Convert level to string
	if message.Level != "" {
		barkReq.Level = string(message.Level)
	}

	// Convert boolean flags to string for Bark API
	if message.Call {
		barkReq.Call = "1"
	}
	if message.AutoCopy {
		barkReq.AutoCopy = "1"
		barkReq.Copy = message.Copy
	}

	// Send request to Bark API
	var barkResp barkResponse
	resp, err := b.client.R().
		SetContext(ctx).
		SetResult(&barkResp).
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetBody(barkReq).
		Post(b.baseURL + "/push")

	if err != nil {
		return &PushResponse{
			Success:  false,
			Error:    fmt.Sprintf("failed to send bark notification: %v", err),
			Provider: b.GetProviderName(),
		}, nil
	}

	if resp.StatusCode() != 200 {
		return &PushResponse{
			Success:  false,
			Error:    fmt.Sprintf("bark API returned status code: %d", resp.StatusCode()),
			Provider: b.GetProviderName(),
		}, nil
	}

	// Check Bark response code
	if barkResp.Code != 200 {
		return &PushResponse{
			Success:  false,
			Error:    fmt.Sprintf("bark API error: %s (code: %d)", barkResp.Message, barkResp.Code),
			Provider: b.GetProviderName(),
		}, nil
	}

	return &PushResponse{
		Success:   true,
		MessageID: fmt.Sprintf("%d", barkResp.Timestamp),
		Provider:  b.GetProviderName(),
	}, nil
}
