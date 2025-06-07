package bilibili

import (
	"fmt"
	"nebulaLive/pkg/live"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2/log"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

const (
	Name   = "bilibili"
	CnName = "哔哩哔哩"

	roomInitUrl = "https://api.live.bilibili.com/room/v1/Room/room_init"
	roomInfoUrl = "https://api.live.bilibili.com/room/v1/Room/get_info"
	streamUrl   = "https://api.live.bilibili.com/xlive/web-room/v2/index/getRoomPlayInfo"
)

type Bilibili struct{}

func (l *Bilibili) parseRealRoomId(roomId string) (string, error) {
	if roomId == "" {
		return "", live.ErrRoomNotExist
	}
	resp, err := live.Client.R().
		SetQueryParam("id", roomId).
		Get(roomInitUrl)
	if err != nil {
		log.Errorw("get real room from room init fail", zap.Error(err))
		return "", fmt.Errorf("failed to get real room id: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return "", live.ErrRoomNotExist
	}
	if gjson.Get(resp.String(), "code").Int() != 0 {
		return "", live.ErrRoomNotExist
	}
	return gjson.Get(resp.String(), "data.room_id").String(), nil
}

func (l *Bilibili) GetInfo(roomId string) (*RoomInfo, error) {
	realRoomId, err := l.parseRealRoomId(roomId)
	if err != nil {
		return nil, err
	}

	var response ResponseWrapper[RoomInfo]
	resp, err := live.Client.R().
		SetQueryParam("room_id", realRoomId).
		SetQueryParam("from", "room").
		SetResult(&response).
		Get(roomInfoUrl)
	if err != nil {
		log.Errorw("get room info fail", zap.Error(err))
		return nil, fmt.Errorf("failed to get room info: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, live.ErrRoomNotExist
	}
	if response.Code != 0 {
		return nil, live.ErrRoomNotExist
	}

	return &response.Data, nil
}

// GetStreams 获取直播流信息
func (l *Bilibili) GetStreams(roomId string, quality Quality, cookies []*http.Cookie) ([]StreamInfo, error) {
	// 解析真实房间号
	realRoomId, err := l.parseRealRoomId(roomId)
	if err != nil {
		return nil, err
	}

	if cookies == nil {
		cookies = []*http.Cookie{}
	}

	var playInfoResp ResponseWrapper[RoomPlayInfoResponse]
	resp, err := live.Client.R().
		SetQueryParams(map[string]string{
			"room_id":  realRoomId,
			"qn":       strconv.Itoa(int(quality)),
			"protocol": "0,1",   // 0: http_stream, 1: http_hls
			"format":   "0,1,2", // 0: flv, 1: ts, 2: fmp4
			"codec":    "0,1",   // 0: avc, 1: hevc
		}).
		SetCookies(cookies).
		SetResult(&playInfoResp).
		Get(streamUrl)

	if err != nil {
		log.Errorw("get room stream info fail", zap.Error(err))
		return nil, fmt.Errorf("failed to get room stream info: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, live.ErrRoomNotExist
	}
	if playInfoResp.Code != 0 {
		return nil, live.ErrRoomNotExist
	}

	streamInfos := make([]StreamInfo, 0)
	streamResp := playInfoResp.Data.PlayurlInfo.Playurl.Stream
	for _, stream := range streamResp {
		for _, format := range stream.Format {
			for _, codec := range format.Codec {
				for _, urlInfo := range codec.URLInfo {
					info := StreamInfo{
						URL:           urlInfo.Host + codec.BaseURL + urlInfo.Extra,
						ProtocolName:  stream.ProtocolName,
						Format:        format.FormatName,
						Codec:         codec.CodecName,
						Quality:       codec.CurrentQn,
						AcceptQuality: codec.AcceptQn,
					}
					streamInfos = append(streamInfos, info)
				}
			}
		}
	}

	return streamInfos, nil
}
