package handler

import "go.uber.org/fx"

// HandlerModule 处理器层模块
var HandlerModule = fx.Options(
	fx.Provide(
		NewUserHandler,
		NewAuthHandler,
		NewRoleHandler,
		NewPermissionHandler,
		NewLiveStreamHandler,
		NewUserPushSettingHandler,
		NewUserPushHandler,
	),
)
