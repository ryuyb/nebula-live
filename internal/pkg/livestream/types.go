package livestream

import "errors"

// StreamStatus represents the status of a live stream
type StreamStatus string

const (
	StreamStatusOnline  StreamStatus = "online"
	StreamStatusOffline StreamStatus = "offline"
)

// StreamInfo contains information about a live stream
type StreamInfo struct {
	Platform string       `json:"platform"`
	RoomID   string       `json:"room_id"`
	Status   StreamStatus `json:"status"`
}

// Common errors
var (
	ErrRoomNotFound     = errors.New("live room not found")
	ErrPlatformNotFound = errors.New("platform not supported")
	ErrInvalidRoomID    = errors.New("invalid room ID")
)
