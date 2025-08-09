package app

import "go.uber.org/fx"

// AppModule 应用层模块
var AppModule = fx.Options(
	fx.Provide(
		NewFiberApp,
	),
)
