package livestream

import (
	"context"
	"encoding/json"
	"fmt"

	"resty.dev/v3"
)

// Bilibili provider implementation
type bilibiliProvider struct {
	client *resty.Client
}

type bilibiliResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    any    `json:"data"` // Can be struct or empty array
}

type bilibiliRoomData struct {
	RoomID        int    `json:"room_id"`
	ShortID       int    `json:"short_id"`
	UID           int    `json:"uid"`
	Title         string `json:"title"`
	Cover         string `json:"cover"`
	UserCover     string `json:"user_cover"`
	Keyframe      string `json:"keyframe"`
	LiveStatus    int    `json:"live_status"`
	LiveStartTime int64  `json:"live_start_time"`
	IsHidden      bool   `json:"is_hidden"`
	IsLocked      bool   `json:"is_locked"`
	IsPortrait    bool   `json:"is_portrait"`
	LiveTime      string `json:"live_time"`
	Tags          string `json:"tags"`
	Description   string `json:"description"`
}

func NewBilibiliProvider(client *resty.Client) Provider {
	return &bilibiliProvider{
		client: client,
	}
}

func (b *bilibiliProvider) GetPlatformName() string {
	return "bilibili"
}

func (b *bilibiliProvider) GetStreamStatus(ctx context.Context, roomID string) (*StreamInfo, error) {
	if roomID == "" {
		return nil, ErrInvalidRoomID
	}

	url := "https://api.live.bilibili.com/room/v1/Room/get_info"

	var bilibiliResp bilibiliResponse
	resp, err := b.client.R().
		SetContext(ctx).
		SetResult(&bilibiliResp).
		SetQueryParam("room_id", roomID).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch bilibili stream status: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("bilibili API returned status code: %d", resp.StatusCode())
	}

	// Check API response code
	if bilibiliResp.Code != 0 {
		// Code 1 means room not found
		if bilibiliResp.Code == 1 {
			return nil, ErrRoomNotFound
		}
		return nil, fmt.Errorf("bilibili API error: %s (code: %d)", bilibiliResp.Message, bilibiliResp.Code)
	}

	streamInfo := &StreamInfo{
		Platform: b.GetPlatformName(),
		RoomID:   roomID,
		Status:   StreamStatusOffline,
	}

	// Parse the data field - it can be a struct or empty array
	dataBytes, err := json.Marshal(bilibiliResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bilibili data: %w", err)
	}

	// Try to unmarshal as room data struct
	var roomData bilibiliRoomData
	if err := json.Unmarshal(dataBytes, &roomData); err == nil {
		// Successfully parsed as room data
		// live_status: 0=not streaming, 1=streaming, 2=rebroadcast
		if roomData.LiveStatus == 1 {
			streamInfo.Status = StreamStatusOnline
		}
	}
	// If unmarshal fails, it's likely an empty array, so room is offline (default)

	return streamInfo, nil
}
