package middleware

import "go.uber.org/fx"

// MiddlewareModule 中间件模块
var MiddlewareModule = fx.Options(
	fx.Provide(
		NewAuthMiddleware,
	),
)