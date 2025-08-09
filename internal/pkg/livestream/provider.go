package livestream

import "context"

// Provider interface for live streaming platforms
type Provider interface {
	GetStreamStatus(ctx context.Context, roomID string) (*StreamInfo, error)
	GetPlatformName() string
}