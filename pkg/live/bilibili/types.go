package bilibili

// RoomInfo 直播间信息
type RoomInfo struct {
	UID              int64      `json:"uid"`                // 主播mid
	RoomID           int64      `json:"room_id"`            // 直播间长号
	ShortID          int64      `json:"short_id"`           // 直播间短号
	Attention        int64      `json:"attention"`          // 关注数量
	Online           int64      `json:"online"`             // 观看人数
	IsPortrait       bool       `json:"is_portrait"`        // 是否竖屏
	Description      string     `json:"description"`        // 描述
	LiveStatus       int        `json:"live_status"`        // 直播状态 0：未开播 1：直播中 2：轮播中
	AreaID           int        `json:"area_id"`            // 分区id
	ParentAreaID     int        `json:"parent_area_id"`     // 父分区id
	ParentAreaName   string     `json:"parent_area_name"`   // 父分区名称
	OldAreaID        int        `json:"old_area_id"`        // 旧版分区id
	Background       string     `json:"background"`         // 背景图片链接
	Title            string     `json:"title"`              // 标题
	UserCover        string     `json:"user_cover"`         // 封面
	Keyframe         string     `json:"keyframe"`           // 关键帧
	IsStrictRoom     bool       `json:"is_strict_room"`     // 未知
	LiveTime         string     `json:"live_time"`          // 直播开始时间
	Tags             string     `json:"tags"`               // 标签
	IsAnchor         int        `json:"is_anchor"`          // 未知
	RoomSilentType   string     `json:"room_silent_type"`   // 禁言状态
	RoomSilentLevel  int        `json:"room_silent_level"`  // 禁言等级
	RoomSilentSecond int        `json:"room_silent_second"` // 禁言时间
	AreaName         string     `json:"area_name"`          // 分区名称
	HotWords         []string   `json:"hot_words"`          // 热词
	HotWordsStatus   int        `json:"hot_words_status"`   // 热词状态
	NewPendants      Pendants   `json:"new_pendants"`       // 头像框\大v
	StudioInfo       StudioInfo `json:"studio_info"`        // 工作室信息
}

// Pendants 头像框和大V信息
type Pendants struct {
	Frame       PendantItem `json:"frame"`        // 头像框
	MobileFrame PendantItem `json:"mobile_frame"` // 手机版头像框
	Badge       PendantItem `json:"badge"`        // 大v
	MobileBadge PendantItem `json:"mobile_badge"` // 手机版大v
}

// PendantItem 头像框或大V的具体信息
type PendantItem struct {
	Name       string `json:"name"`         // 名称
	Value      string `json:"value"`        // 值
	Position   int    `json:"position"`     // 位置
	Desc       string `json:"desc"`         // 描述
	Area       int    `json:"area"`         // 分区
	AreaOld    int    `json:"area_old"`     // 旧分区
	BgColor    string `json:"bg_color"`     // 背景色
	BgPic      string `json:"bg_pic"`       // 背景图
	UseOldArea bool   `json:"use_old_area"` // 是否旧分区号
}

// StudioInfo 工作室信息
type StudioInfo struct {
	Status     int   `json:"status"`      // 状态
	MasterList []int `json:"master_list"` // 主播列表
}

type ResponseWrapper[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// StreamInfo 播放流信息
type StreamInfo struct {
	URL           string `json:"url"`            // 播放地址
	ProtocolName  string `json:"protocol_name"`  // 协议名
	Format        string `json:"format"`         // 格式
	Codec         string `json:"codec"`          // 编码
	Quality       int    `json:"quality"`        // 清晰度
	AcceptQuality []int  `json:"accept_quality"` // 可选择的清晰度
}

// RoomPlayInfoResponse 播放流API响应
type RoomPlayInfoResponse struct {
	RoomID          int   `json:"room_id"`      // 直播间id
	ShortID         int   `json:"short_id"`     // 直播间短id
	UID             int   `json:"uid"`          // 主播uid
	IsHidden        bool  `json:"is_hidden"`    // 直播间是否被隐藏
	IsLocked        bool  `json:"is_locked"`    // 直播间是否被锁定
	IsPortrait      bool  `json:"is_portrait"`  // 是否竖屏
	LiveStatus      int   `json:"live_status"`  // 直播状态 0：未开播 1：直播中 2：轮播中
	HiddenTill      int   `json:"hidden_till"`  // 隐藏结束时间
	LockTill        int   `json:"lock_till"`    // 封禁结束时间 秒级时间戳
	Encrypted       bool  `json:"encrypted"`    // 直播间为加密直播间
	PwdVerified     bool  `json:"pwd_verified"` // 是否通过密码验证 当encrypted为true时才有意义
	LiveTime        int   `json:"live_time"`    // 本次开播时间 秒级时间戳
	RoomShield      int   `json:"room_shield"`
	AllSpecialTypes []any `json:"all_special_types"`
	PlayurlInfo     struct {
		ConfJSON string `json:"conf_json"`
		Playurl  struct {
			Cid     int `json:"cid"` // 直播间id
			GQnDesc []struct {
				Qn       int    `json:"qn"`   // 清晰度代码
				Desc     string `json:"desc"` // 清晰度描述
				HdrDesc  string `json:"hdr_desc"`
				AttrDesc any    `json:"attr_desc"`
			} `json:"g_qn_desc"` // 清晰度列表
			Stream []struct {
				ProtocolName string `json:"protocol_name"` // 协议名
				Format       []struct {
					FormatName string `json:"format_name"` // 格式名
					Codec      []struct {
						CodecName string `json:"codec_name"` // 编码名
						CurrentQn int    `json:"current_qn"` // 当前清晰度编码
						AcceptQn  []int  `json:"accept_qn"`  // 可用清晰度编码列表
						BaseURL   string `json:"base_url"`   // 播放源路径
						URLInfo   []struct {
							Host      string `json:"host"`  // 域名
							Extra     string `json:"extra"` // URL参数
							StreamTTL int    `json:"stream_ttl"`
						} `json:"url_info"` // 域名信息列表
						HdrQn     any    `json:"hdr_qn"`
						DolbyType int    `json:"dolby_type"`
						AttrName  string `json:"attr_name"`
					} `json:"codec"` // 编码列表
				} `json:"format"` // 格式列表
			} `json:"stream"` // 直播流信息
			P2PData struct {
				P2P      bool `json:"p2p"`
				P2PType  int  `json:"p2p_type"`
				MP2P     bool `json:"m_p2p"`
				MServers any  `json:"m_servers"`
			} `json:"p2p_data"`
			DolbyQn any `json:"dolby_qn"`
		} `json:"playurl"`
	} `json:"playurl_info"` // 直播流信息
}

// Quality 可接受的清晰度
type Quality int

const (
	Dolby    Quality = 30000 // 杜比
	FourK    Quality = 20000 // 4K
	Original Quality = 10000 // 原画
	BluRay   Quality = 400   // 蓝光
	UltraHd  Quality = 250   //超清
	Hd       Quality = 150   // 高清
	Smooth   Quality = 80    // 流畅
)
