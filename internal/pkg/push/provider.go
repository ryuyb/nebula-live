package push

import (
	"context"
)

// Provider defines the interface for push notification providers
type Provider interface {
	// SendMessage sends a push notification message
	SendMessage(ctx context.Context, message *PushMessage) (*PushResponse, error)

	// GetProviderName returns the name of the provider
	GetProviderName() string

	// IsEnabled returns whether the provider is enabled
	IsEnabled() bool

	// ValidateMessage validates if the message is compatible with this provider
	ValidateMessage(message *PushMessage) error
}
