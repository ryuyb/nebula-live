package livestream

import (
	"context"
	"fmt"

	"resty.dev/v3"
)

// Douyu provider implementation
type douyuProvider struct {
	client *resty.Client
}

type douyuResponse struct {
	Room struct {
		ShowStatus int `json:"show_status"`
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
