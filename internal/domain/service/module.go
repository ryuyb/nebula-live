package service

import "go.uber.org/fx"

// ServiceModule 服务层模块
var ServiceModule = fx.Options(
	fx.Provide(
		NewUserService,
		NewRBACService,
		NewLiveStreamService,
	),
)
