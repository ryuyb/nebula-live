package bilibili

import (
	"encoding/json"
	"fmt"
	"io"
	"nebulaLive/pkg/live"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2/log"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"resty.dev/v3"
)

const (
	Name   = "bilibili"
	CnName = "哔哩哔哩"

	roomInitUrl = "https://api.live.bilibili.com/room/v1/Room/room_init"
	roomInfoUrl = "https://api.live.bilibili.com/room/v1/Room/get_info"
	streamUrl   = "https://api.live.bilibili.com/room/v2/index/getRoomPlayInfo"
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

	var response RoomInfoResponse
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

// GetStream 获取直播流信息
func (b *Bilibili) GetStream(roomId string, quality int) (*StreamInfo, error) {
	// 解析真实房间号
	realRoomId, err := b.parseRealRoomId(roomId)
	if err != nil {
		return nil, fmt.Errorf("获取真实房间号失败: %v", err)
	}

	// 构建请求参数
	params := map[string]string{
		"room_id":  realRoomId,
		"qn":       strconv.Itoa(quality),
		"protocol": "0,1",   // 0: http_stream, 1: http_hls
		"format":   "0,1,2", // 0: flv, 1: ts, 2: fmp4
		"codec":    "0,1",   // 0: avc, 1: hevc
	}

	// 发送请求
	client := resty.New()
	req := client.R()
	for k, v := range params {
		req.SetQueryParam(k, v)
	}
	resp, err := req.Get(streamUrl)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var streamResp StreamResponse
	if err := json.Unmarshal(body, &streamResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查响应状态
	if streamResp.Code != 0 {
		return nil, fmt.Errorf("API返回错误: %s", streamResp.Message)
	}

	// 检查直播状态
	if streamResp.Data.LiveStatus != 1 {
		return nil, fmt.Errorf("直播间未开播")
	}

	// 获取播放流信息
	if len(streamResp.Data.PlayURLInfo.PlayURL.StreamInfo) == 0 {
		return nil, fmt.Errorf("未找到可用的播放流")
	}

	// 遍历流信息找到合适的流
	for _, stream := range streamResp.Data.PlayURLInfo.PlayURL.StreamInfo {
		if stream.ProtocolName == "http_stream" {
			for _, format := range stream.Format {
				if format.FormatName == "flv" {
					for _, codec := range format.Codec {
						if codec.CurrentQn == quality {
							// 构建完整的播放地址
							url := codec.BaseURL
							if len(codec.URLInfo) > 0 {
								url = codec.URLInfo[0].Host + codec.BaseURL + codec.URLInfo[0].Extra
							}

							return &StreamInfo{
								URL:       url,
								Format:    format.FormatName,
								Codec:     codec.CodecName,
								Quality:   codec.CurrentQn,
								Bandwidth: 0,  // 新API中没有带宽信息
								FrameRate: "", // 新API中没有帧率信息
							}, nil
						}
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("未找到指定清晰度的播放流")
}
