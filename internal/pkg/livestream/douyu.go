package livestream

import (
	"context"
	"fmt"
	"strconv"

	"resty.dev/v3"
)

// Douyu provider implementation
type douyuProvider struct {
	client *resty.Client
}

type douyuResponse struct {
	Room struct {
		ShowStatus   int    `json:"show_status"`
		RoomName     string `json:"room_name"`
		OwnerUID     int    `json:"owner_uid"`
		Nickname     string `json:"nickname"`
		RoomSrc      string `json:"room_src"`
		Avatar       struct {
			Big    string `json:"big"`
			Middle string `json:"middle"`
			Small  string `json:"small"`
		} `json:"avatar"`
		CateName     string `json:"cate_name"`
		ShowDetails  string `json:"show_details"`
		ShowTime     int64  `json:"show_time"`
		RoomPic      string `json:"room_pic"`
		CoverSrc     string `json:"coverSrc"`
		RoomBizAll   struct {
			Hot string `json:"hot"`
		} `json:"room_biz_all"`
	} `json:"room"`
}

func NewDouyuProvider(client *resty.Client) Provider {
	return &douyuProvider{
		client: client,
	}
}

func (d *douyuProvider) GetPlatformName() string {
	return "douyu"
}

func (d *douyuProvider) GetStreamStatus(ctx context.Context, roomID string) (*StreamInfo, error) {
	if roomID == "" {
		return nil, ErrInvalidRoomID
	}

	url := fmt.Sprintf("https://www.douyu.com/betard/%s", roomID)

	var douyuResp douyuResponse
	resp, err := d.client.R().
		SetContext(ctx).
		SetResult(&douyuResp).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch douyu stream status: %w", err)
	}

	if resp.StatusCode() == 404 {
		return nil, ErrRoomNotFound
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("douyu API returned status code: %d", resp.StatusCode())
	}

	streamInfo := &StreamInfo{
		Platform: d.GetPlatformName(),
		RoomID:   roomID,
		Status:   StreamStatusOffline,
	}

	if douyuResp.Room.ShowStatus == 1 {
		streamInfo.Status = StreamStatusOnline
	}

	return streamInfo, nil
}

func (d *douyuProvider) GetRoomInfo(ctx context.Context, roomID string) (*RoomInfo, error) {
	if roomID == "" {
		return nil, ErrInvalidRoomID
	}

	url := fmt.Sprintf("https://www.douyu.com/betard/%s", roomID)

	var douyuResp douyuResponse
	resp, err := d.client.R().
		SetContext(ctx).
		SetResult(&douyuResp).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch douyu room info: %w", err)
	}

	if resp.StatusCode() == 404 {
		return nil, ErrRoomNotFound
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("douyu API returned status code: %d", resp.StatusCode())
	}

	// Parse viewer count from string
	var viewerCount int64
	if douyuResp.Room.RoomBizAll.Hot != "" {
		if count, err := strconv.ParseInt(douyuResp.Room.RoomBizAll.Hot, 10, 64); err == nil {
			viewerCount = count
		}
	}

	roomInfo := &RoomInfo{
		Platform:      d.GetPlatformName(),
		RoomID:        roomID,
		Status:        StreamStatusOffline,
		Title:         douyuResp.Room.RoomName,
		Description:   douyuResp.Room.ShowDetails,
		Cover:         douyuResp.Room.CoverSrc,
		Keyframe:      douyuResp.Room.RoomPic,
		OwnerID:       strconv.Itoa(douyuResp.Room.OwnerUID),
		OwnerName:     douyuResp.Room.Nickname,
		OwnerAvatar:   douyuResp.Room.Avatar.Big,
		LiveStartTime: douyuResp.Room.ShowTime,
		ViewerCount:   viewerCount,
		Category:      douyuResp.Room.CateName,
	}

	if douyuResp.Room.ShowStatus == 1 {
		roomInfo.Status = StreamStatusOnline
	}

	return roomInfo, nil
}
