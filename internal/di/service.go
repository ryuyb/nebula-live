package di

import (
	"go.uber.org/fx"
	"nebulaLive/internal/service"
)

// ServiceModule 提供服务相关的依赖项
var ServiceModule = fx.Options(
	fx.Provide(
		service.NewUserService,
	),
)
