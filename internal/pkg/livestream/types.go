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

// RoomInfo contains detailed information about a live room
type RoomInfo struct {
	Platform      string       `json:"platform"`
	RoomID        string       `json:"room_id"`
	Status        StreamStatus `json:"status"`
	Title         string       `json:"title,omitempty"`
	Description   string       `json:"description,omitempty"`
	Cover         string       `json:"cover,omitempty"`
	Keyframe      string       `json:"keyframe,omitempty"`
	OwnerID       string       `json:"owner_id,omitempty"`
	OwnerName     string       `json:"owner_name,omitempty"`
	OwnerAvatar   string       `json:"owner_avatar,omitempty"`
	LiveStartTime int64        `json:"live_start_time,omitempty"`
	ViewerCount   int64        `json:"viewer_count,omitempty"`
	Category      string       `json:"category,omitempty"`
}

// Common errors
var (
	ErrRoomNotFound     = errors.New("live room not found")
	ErrPlatformNotFound = errors.New("platform not supported")
	ErrInvalidRoomID    = errors.New("invalid room ID")
)
