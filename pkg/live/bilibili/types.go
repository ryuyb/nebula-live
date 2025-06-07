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

// RoomInfoResponse API响应结构
type RoomInfoResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    RoomInfo `json:"data"`
}

// StreamInfo 播放流信息
type StreamInfo struct {
	URL       string `json:"url"`        // 播放地址
	Format    string `json:"format"`     // 格式
	Codec     string `json:"codec"`      // 编码
	Quality   int    `json:"quality"`    // 清晰度
	Bandwidth int    `json:"bandwidth"`  // 带宽
	FrameRate string `json:"frame_rate"` // 帧率
}

// StreamResponse 播放流API响应
type StreamResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		RoomID           int64    `json:"room_id"`            // 房间号
		ShortID          int64    `json:"short_id"`           // 短号
		UID              int64    `json:"uid"`                // 主播uid
		NeedP2P          int      `json:"need_p2p"`           // 是否需要P2P
		IsHidden         bool     `json:"is_hidden"`          // 是否隐藏
		IsLocked         bool     `json:"is_locked"`          // 是否锁定
		LockTill         int      `json:"lock_till"`          // 锁定时间
		HiddenTill       int      `json:"hidden_till"`        // 隐藏时间
		Encrypted        bool     `json:"encrypted"`          // 是否加密
		PwdVerified      bool     `json:"pwd_verified"`       // 密码验证
		LiveStatus       int      `json:"live_status"`        // 直播状态
		LiveTime         int      `json:"live_time"`          // 直播时长
		LiveType         int      `json:"live_type"`          // 直播类型
		AreaID           int      `json:"area_id"`            // 分区ID
		ParentAreaID     int      `json:"parent_area_id"`     // 父分区ID
		ParentAreaName   string   `json:"parent_area_name"`   // 父分区名称
		OldAreaID        int      `json:"old_area_id"`        // 旧分区ID
		Background       string   `json:"background"`         // 背景图
		Title            string   `json:"title"`              // 标题
		UserCover        string   `json:"user_cover"`         // 用户封面
		Keyframe         string   `json:"keyframe"`           // 关键帧
		IsStrictRoom     bool     `json:"is_strict_room"`     // 是否严格房间
		LiveTimeStr      string   `json:"live_time_str"`      // 直播时长字符串
		Tags             string   `json:"tags"`               // 标签
		IsAnchor         int      `json:"is_anchor"`          // 是否主播
		RoomSilentType   string   `json:"room_silent_type"`   // 房间禁言类型
		RoomSilentLevel  int      `json:"room_silent_level"`  // 房间禁言等级
		RoomSilentSecond int      `json:"room_silent_second"` // 房间禁言时间
		AreaName         string   `json:"area_name"`          // 分区名称
		Pendants         string   `json:"pendants"`           // 头像框
		AreaPendants     string   `json:"area_pendants"`      // 分区头像框
		HotWords         []string `json:"hot_words"`          // 热词
		HotWordsStatus   int      `json:"hot_words_status"`   // 热词状态
		Verify           string   `json:"verify"`             // 认证信息
		NewPendants      struct {
			Frame struct {
				Name       string `json:"name"`         // 名称
				Value      string `json:"value"`        // 值
				Position   int    `json:"position"`     // 位置
				Desc       string `json:"desc"`         // 描述
				Area       int    `json:"area"`         // 分区
				AreaOld    int    `json:"area_old"`     // 旧分区
				BgColor    string `json:"bg_color"`     // 背景色
				BgPic      string `json:"bg_pic"`       // 背景图
				UseOldArea bool   `json:"use_old_area"` // 是否使用旧分区
			} `json:"frame"` // 头像框
			Badge struct {
				Name     string `json:"name"`     // 名称
				Position int    `json:"position"` // 位置
				Value    string `json:"value"`    // 值
				Desc     string `json:"desc"`     // 描述
			} `json:"badge"` // 大V
		} `json:"new_pendants"` // 新头像框
		UpSession            string `json:"up_session"`              // 上行会话
		PkStatus             int    `json:"pk_status"`               // PK状态
		PkID                 int    `json:"pk_id"`                   // PK ID
		BattleID             int    `json:"battle_id"`               // 战斗ID
		AllowChangeAreaTime  int    `json:"allow_change_area_time"`  // 允许更换分区时间
		AllowUploadCoverTime int    `json:"allow_upload_cover_time"` // 允许上传封面时间
		StudioInfo           struct {
			Status     int   `json:"status"`      // 状态
			MasterList []int `json:"master_list"` // 主播列表
		} `json:"studio_info"` // 工作室信息
		PlayURLInfo struct {
			ConfJSON string `json:"conf_json"` // 配置JSON
			PlayURL  struct {
				StreamInfo []struct {
					ProtocolName string `json:"protocol_name"` // 协议名称
					Format       []struct {
						FormatName string `json:"format_name"` // 格式名称
						Codec      []struct {
							CodecName string `json:"codec_name"` // 编码名称
							CurrentQn int    `json:"current_qn"` // 当前清晰度
							AcceptQn  []int  `json:"accept_qn"`  // 可接受的清晰度
							BaseURL   string `json:"base_url"`   // 基础URL
							URLInfo   []struct {
								Host      string `json:"host"`       // 主机
								Extra     string `json:"extra"`      // 额外信息
								StreamTTL int    `json:"stream_ttl"` // 流TTL
							} `json:"url_info"` // URL信息
							HdrQn     interface{} `json:"hdr_qn"`     // HDR清晰度
							DolbyType int         `json:"dolby_type"` // 杜比类型
							AttrName  string      `json:"attr_name"`  // 属性名称
						} `json:"codec"` // 编码
					} `json:"format"` // 格式
				} `json:"stream_info"` // 流信息
			} `json:"play_url"` // 播放URL
			P2PData struct {
				P2P      bool   `json:"p2p"`       // 是否P2P
				P2PType  int    `json:"p2p_type"`  // P2P类型
				MP2P     bool   `json:"m_p2p"`     // 是否移动P2P
				MServers string `json:"m_servers"` // 移动服务器
			} `json:"p2p_data"` // P2P数据
			DolbyQn interface{} `json:"dolby_qn"` // 杜比清晰度
		} `json:"play_url_info"` // 播放URL信息
	} `json:"data"`
}
