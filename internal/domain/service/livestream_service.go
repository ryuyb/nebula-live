package service

import (
	"context"

	"nebula-live/internal/pkg/livestream"
)

// LiveStreamService manages multiple live streaming platforms
type LiveStreamService interface {
	GetStreamStatus(ctx context.Context, platformName string, roomID string) (*livestream.StreamInfo, error)
	GetRoomInfo(ctx context.Context, platformName string, roomID string) (*livestream.RoomInfo, error)
	GetSupportedPlatforms() []string
}

type liveStreamService struct {
	client *livestream.Client
}

func NewLiveStreamService() LiveStreamService {
	return &liveStreamService{
		client: livestream.NewClient(),
	}
}

func (s *liveStreamService) GetStreamStatus(ctx context.Context, platformName string, roomID string) (*livestream.StreamInfo, error) {
	return s.client.GetStreamStatus(ctx, platformName, roomID)
}

func (s *liveStreamService) GetRoomInfo(ctx context.Context, platformName string, roomID string) (*livestream.RoomInfo, error) {
	return s.client.GetRoomInfo(ctx, platformName, roomID)
}

func (s *liveStreamService) GetSupportedPlatforms() []string {
	return s.client.GetSupportedPlatforms()
}
