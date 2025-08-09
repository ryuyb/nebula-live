package push

import (
	"errors"
)

// PushLevel represents the notification level
type PushLevel string

const (
	PushLevelCritical      PushLevel = "critical"
	PushLevelActive        PushLevel = "active"
	PushLevelTimeSensitive PushLevel = "timeSensitive"
	PushLevelPassive       PushLevel = "passive"
)

// PushMessage represents a push notification message
type PushMessage struct {
	Title    string            `json:"title,omitempty"`
	Subtitle string            `json:"subtitle,omitempty"`
	Body     string            `json:"body"`
	DeviceID string            `json:"device_id"`
	Badge    int               `json:"badge,omitempty"`
	Sound    string            `json:"sound,omitempty"`
	Icon     string            `json:"icon,omitempty"`
	Group    string            `json:"group,omitempty"`
	URL      string            `json:"url,omitempty"`
	Level    PushLevel         `json:"level,omitempty"`
	Call     bool              `json:"call,omitempty"`
	AutoCopy bool              `json:"auto_copy,omitempty"`
	Copy     string            `json:"copy,omitempty"`
	Extra    map[string]string `json:"extra,omitempty"`
}

// PushResponse represents the response from a push provider
type PushResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id,omitempty"`
	Error     string `json:"error,omitempty"`
	Provider  string `json:"provider"`
}

// Common errors for push notifications
var (
	ErrInvalidDeviceID    = errors.New("invalid device ID")
	ErrEmptyMessage       = errors.New("message body cannot be empty")
	ErrProviderNotFound   = errors.New("push provider not found")
	ErrProviderNotEnabled = errors.New("push provider not enabled")
	ErrSendFailed         = errors.New("failed to send push notification")
)
