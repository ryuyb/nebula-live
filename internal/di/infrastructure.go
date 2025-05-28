package di

import (
	"nebulaLive/internal/config"
	"nebulaLive/pkg/logger"

	"go.uber.org/fx"
)

// InfrastructureModule 提供基础设施相关的依赖项
var InfrastructureModule = fx.Options(
	fx.Provide(
		config.GetConfig,
		logger.NewLoggerWithDefaults,
	),
)
