package di

import (
	"go.uber.org/fx"
	"nebulaLive/internal/api/router"
)

func asRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(router.Router)),
		fx.ResultTags(`group:"routers"`),
	)
}

// RouterModule 提供路由相关的依赖项
var RouterModule = fx.Options(
	fx.Provide(
		asRoute(router.NewUserRouter),
		fx.Annotate(
			router.NewRouterRegistry,
			fx.ParamTags(`group:"routers"`),
		),
	),
)
