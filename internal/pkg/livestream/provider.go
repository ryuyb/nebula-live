package livestream

import "context"

// Provider interface for live streaming platforms
type Provider interface {
	GetStreamStatus(ctx context.Context, roomID string) (*StreamInfo, error)
	GetRoomInfo(ctx context.Context, roomID string) (*RoomInfo, error)
	GetPlatformName() string
}
