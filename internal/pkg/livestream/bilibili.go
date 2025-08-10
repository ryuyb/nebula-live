package livestream

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

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
	UID              int    `json:"uid"`
	RoomID           int    `json:"room_id"`
	ShortID          int    `json:"short_id"`
	Attention        int    `json:"attention"`
	Online           int    `json:"online"`
	IsPortrait       bool   `json:"is_portrait"`
	Description      string `json:"description"`
	LiveStatus       int    `json:"live_status"`
	AreaID           int    `json:"area_id"`
	ParentAreaID     int    `json:"parent_area_id"`
	ParentAreaName   string `json:"parent_area_name"`
	OldAreaID        int    `json:"old_area_id"`
	Background       string `json:"background"`
	Title            string `json:"title"`
	UserCover        string `json:"user_cover"`
	Keyframe         string `json:"keyframe"`
	IsStrictRoom     bool   `json:"is_strict_room"`
	LiveTime         string `json:"live_time"`
	Tags             string `json:"tags"`
	IsAnchor         int    `json:"is_anchor"`
	RoomSilentType   string `json:"room_silent_type"`
	RoomSilentLevel  int    `json:"room_silent_level"`
	RoomSilentSecond int    `json:"room_silent_second"`
	AreaName         string `json:"area_name"`
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

func (b *bilibiliProvider) GetRoomInfo(ctx context.Context, roomID string) (*RoomInfo, error) {
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
		return nil, fmt.Errorf("failed to fetch bilibili room info: %w", err)
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

	roomInfo := &RoomInfo{
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
		roomInfo.Title = roomData.Title
		roomInfo.Description = roomData.Description
		roomInfo.Cover = roomData.UserCover
		roomInfo.Keyframe = roomData.Keyframe
		roomInfo.OwnerID = strconv.Itoa(roomData.UID)
		roomInfo.ViewerCount = int64(roomData.Online)
		roomInfo.Category = roomData.AreaName

		// live_status: 0=not streaming, 1=streaming, 2=rebroadcast
		if roomData.LiveStatus == 1 {
			roomInfo.Status = StreamStatusOnline
		}

		// Get owner information from user API
		ownerInfo, err := b.getOwnerInfo(ctx, roomData.UID)
		if err == nil {
			roomInfo.OwnerName = ownerInfo.Name
			roomInfo.OwnerAvatar = ownerInfo.Avatar
		}
	}

	return roomInfo, nil
}

type bilibiliMasterResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    struct {
		Info struct {
			UID              int    `json:"uid"`
			UName            string `json:"uname"`
			Face             string `json:"face"`
			OfficialVerify   struct {
				Type int    `json:"type"`
				Desc string `json:"desc"`
			} `json:"official_verify"`
			Gender int `json:"gender"`
		} `json:"info"`
		Exp struct {
			MasterLevel struct {
				Level   int   `json:"level"`
				Color   int   `json:"color"`
				Current []int `json:"current"`
				Next    []int `json:"next"`
			} `json:"master_level"`
		} `json:"exp"`
		FollowerNum   int    `json:"follower_num"`
		RoomID        int    `json:"room_id"`
		MedalName     string `json:"medal_name"`
		GloryCount    int    `json:"glory_count"`
		Pendant       string `json:"pendant"`
		LinkGroupNum  int    `json:"link_group_num"`
		RoomNews      struct {
			Content   string `json:"content"`
			CTime     string `json:"ctime"`
			CTimeText string `json:"ctime_text"`
		} `json:"room_news"`
	} `json:"data"`
}

func (b *bilibiliProvider) getOwnerInfo(ctx context.Context, uid int) (*struct {
	Name   string
	Avatar string
}, error) {
	url := "https://api.live.bilibili.com/live_user/v1/Master/info"
	
	var masterResp bilibiliMasterResponse
	resp, err := b.client.R().
		SetContext(ctx).
		SetResult(&masterResp).
		SetQueryParam("uid", strconv.Itoa(uid)).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		Get(url)

	if err != nil || resp.StatusCode() != 200 || masterResp.Code != 0 {
		return nil, fmt.Errorf("failed to fetch master info")
	}

	return &struct {
		Name   string
		Avatar string
	}{
		Name:   masterResp.Data.Info.UName,
		Avatar: masterResp.Data.Info.Face,
	}, nil
}
