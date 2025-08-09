package router

import (
	"go.uber.org/fx"
)

// asRoute 将路由器标记为Router组的成员
func asRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Router)),
		fx.ResultTags(`group:"routers"`),
	)
}

// RouterModule 路由模块，提供所有路由相关的依赖
var RouterModule = fx.Options(
	// 提供各种路由器
	fx.Provide(asRoute(NewUserRouter)),
	fx.Provide(asRoute(NewAuthRouter)),
	fx.Provide(asRoute(NewRoleRouter)),
	fx.Provide(asRoute(NewPermissionRouter)),
	fx.Provide(asRoute(NewLiveStreamRouter)),
	fx.Provide(asRoute(NewUserPushSettingRouter)),
	fx.Provide(asRoute(NewUserPushRouter)),

	// 提供路由注册器
	fx.Provide(NewRouterRegistry),
)
