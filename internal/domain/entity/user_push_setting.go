package entity

import (
	"encoding/json"
	"time"
)

// UserPushSetting 用户推送设置实体
type UserPushSetting struct {
	ID         uint                   `json:"id"`
	UserID     uint                   `json:"user_id"`
	Provider   string                 `json:"provider"`        // 推送服务提供商（如：bark）
	Enabled    bool                   `json:"enabled"`         // 是否启用
	DeviceID   string                 `json:"device_id"`       // 设备ID
	DeviceName string                 `json:"device_name"`     // 设备名称
	Settings   map[string]interface{} `json:"settings"`        // 提供商特定设置
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// BarkSettings Bark推送服务的特定设置
type BarkSettings struct {
	BaseURL    string `json:"base_url,omitempty"`    // 自定义Bark服务器地址
	Sound      string `json:"sound,omitempty"`       // 默认铃声
	Icon       string `json:"icon,omitempty"`        // 默认图标
	Group      string `json:"group,omitempty"`       // 默认分组
	Level      string `json:"level,omitempty"`       // 默认通知级别
	AutoCopy   bool   `json:"auto_copy,omitempty"`   // 是否自动复制
	Call       bool   `json:"call,omitempty"`        // 是否响铃30秒
}

// GetBarkSettings 获取Bark设置
func (ups *UserPushSetting) GetBarkSettings() (*BarkSettings, error) {
	if ups.Provider != "bark" {
		return nil, nil
	}
	
	if ups.Settings == nil {
		return &BarkSettings{}, nil
	}
	
	settingsBytes, err := json.Marshal(ups.Settings)
	if err != nil {
		return nil, err
	}
	
	var barkSettings BarkSettings
	err = json.Unmarshal(settingsBytes, &barkSettings)
	if err != nil {
		return nil, err
	}
	
	return &barkSettings, nil
}

// SetBarkSettings 设置Bark设置
func (ups *UserPushSetting) SetBarkSettings(settings *BarkSettings) error {
	if ups.Provider != "bark" {
		return nil
	}
	
	settingsMap := make(map[string]interface{})
	settingsBytes, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	
	err = json.Unmarshal(settingsBytes, &settingsMap)
	if err != nil {
		return err
	}
	
	ups.Settings = settingsMap
	ups.UpdatedAt = time.Now()
	return nil
}

// IsValid 检查推送设置是否有效
func (ups *UserPushSetting) IsValid() bool {
	if ups.UserID == 0 || ups.Provider == "" || ups.DeviceID == "" {
		return false
	}
	return true
}

// Enable 启用推送设置
func (ups *UserPushSetting) Enable() {
	ups.Enabled = true
	ups.UpdatedAt = time.Now()
}

// Disable 禁用推送设置
func (ups *UserPushSetting) Disable() {
	ups.Enabled = false
	ups.UpdatedAt = time.Now()
}

// UpdateDeviceName 更新设备名称
func (ups *UserPushSetting) UpdateDeviceName(name string) {
	ups.DeviceName = name
	ups.UpdatedAt = time.Now()
}