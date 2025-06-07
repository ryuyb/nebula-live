package live

import "resty.dev/v3"

var (
	Client = resty.New()
)

type Live interface {
}

type RoomInfo struct {
	RoomId     string // 直播间id
	ShortId    string //直播间短号
	Attention  int64
	Online     int64
	IsPortrait bool // 是否竖屏
}
