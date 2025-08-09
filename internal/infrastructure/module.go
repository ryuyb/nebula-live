package infrastructure

import (
	"nebula-live/internal/infrastructure/config"
	"nebula-live/internal/infrastructure/logger"
	"nebula-live/internal/infrastructure/persistence"

	"go.uber.org/fx"
)

// InfrastructureModule 基础设施层模块
var InfrastructureModule = fx.Options(
	fx.Provide(
		config.NewConfig,
		logger.NewLogger,
		persistence.NewEntClient,
	),
)
