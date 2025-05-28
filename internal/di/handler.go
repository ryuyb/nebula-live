package di

import (
	"go.uber.org/fx"
	"nebulaLive/internal/api/handler"
)

// HandlerModule 提供处理程序相关的依赖项
var HandlerModule = fx.Options(
	fx.Provide(
		handler.NewUserHandler,
	),
)
