package dto

import (
	"errors"
	"time"
)

// CreateUserPushSettingRequest 创建用户推送设置请求
type CreateUserPushSettingRequest struct {
	Provider   string                 `json:"provider" validate:"required,oneof=bark"`
	DeviceID   string                 `json:"device_id" validate:"required,min=1,max=255"`
	DeviceName string                 `json:"device_name" validate:"required,min=1,max=100"`
	Settings   map[string]interface{} `json:"settings,omitempty"`
}

// Validate 验证创建用户推送设置请求
func (r *CreateUserPushSettingRequest) Validate() error {
	if r.Provider == "" {
		return errors.New("provider is required")
	}
	
	if r.Provider != "bark" {
		return errors.New("provider must be one of: bark")
	}
	
	if r.DeviceID == "" {
		return errors.New("device_id is required")
	}
	
	if len(r.DeviceID) > 255 {
		return errors.New("device_id must not exceed 255 characters")
	}
	
	if r.DeviceName == "" {
		return errors.New("device_name is required")
	}
	
	if len(r.DeviceName) > 100 {
		return errors.New("device_name must not exceed 100 characters")
	}
	
	return nil
}

// UpdateUserPushSettingRequest 更新用户推送设置请求
type UpdateUserPushSettingRequest struct {
	Enabled    *bool                  `json:"enabled,omitempty"`
	DeviceName *string                `json:"device_name,omitempty"`
	Settings   map[string]interface{} `json:"settings,omitempty"`
}

// Validate 验证更新用户推送设置请求
func (r *UpdateUserPushSettingRequest) Validate() error {
	if r.DeviceName != nil && *r.DeviceName == "" {
		return errors.New("device_name cannot be empty")
	}
	
	if r.DeviceName != nil && len(*r.DeviceName) > 100 {
		return errors.New("device_name must not exceed 100 characters")
	}
	
	return nil
}

// ValidateDeviceRequest 验证设备请求
type ValidateDeviceRequest struct {
	Provider string `json:"provider" validate:"required,oneof=bark"`
	DeviceID string `json:"device_id" validate:"required,min=1,max=255"`
}

// Validate 验证设备请求
func (r *ValidateDeviceRequest) Validate() error {
	if r.Provider == "" {
		return errors.New("provider is required")
	}
	
	if r.Provider != "bark" {
		return errors.New("provider must be one of: bark")
	}
	
	if r.DeviceID == "" {
		return errors.New("device_id is required")
	}
	
	if len(r.DeviceID) > 255 {
		return errors.New("device_id must not exceed 255 characters")
	}
	
	return nil
}

// UserPushSettingResponse 用户推送设置响应
type UserPushSettingResponse struct {
	ID         uint                   `json:"id"`
	UserID     uint                   `json:"user_id"`
	Provider   string                 `json:"provider"`
	Enabled    bool                   `json:"enabled"`
	DeviceID   string                 `json:"device_id"`
	DeviceName string                 `json:"device_name"`
	Settings   map[string]interface{} `json:"settings,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// UserPushRequest 用户推送请求
type UserPushRequest struct {
	Title    string `json:"title" validate:"required,min=1,max=200"`
	Body     string `json:"body" validate:"required,min=1,max=1000"`
	URL      string `json:"url,omitempty"`
	Sound    string `json:"sound,omitempty"`
	Icon     string `json:"icon,omitempty"`
	Group    string `json:"group,omitempty"`
	Level    string `json:"level,omitempty"`
	AutoCopy bool   `json:"auto_copy,omitempty"`
	Call     bool   `json:"call,omitempty"`
}

// Validate 验证用户推送请求
func (r *UserPushRequest) Validate() error {
	if r.Title == "" {
		return errors.New("title is required")
	}
	
	if len(r.Title) > 200 {
		return errors.New("title must not exceed 200 characters")
	}
	
	if r.Body == "" {
		return errors.New("body is required")
	}
	
	if len(r.Body) > 1000 {
		return errors.New("body must not exceed 1000 characters")
	}
	
	return nil
}

// PushResponse 推送响应
type PushResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id,omitempty"`
	Provider  string `json:"provider"`
	Error     string `json:"error,omitempty"`
}

// UserPushResult 用户推送结果
type UserPushResult struct {
	UserID       uint           `json:"user_id"`
	Provider     string         `json:"provider,omitempty"`
	TotalDevices int            `json:"total_devices"`
	SuccessCount int            `json:"success_count"`
	FailedCount  int            `json:"failed_count"`
	Responses    []PushResponse `json:"responses"`
	Message      string         `json:"message,omitempty"`
}

// ListResponse 通用列表响应
type ListResponse[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}