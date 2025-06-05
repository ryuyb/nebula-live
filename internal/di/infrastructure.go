package di

import (
	"nebulaLive/internal/config"
	"nebulaLive/pkg/logger"

	"go.uber.org/fx"
)

func NewLogger(cfg *config.Config) (*logger.Logger, error) {
	return logger.NewLogger(logger.LoggerConfig{
		Level:           cfg.Logging.Level,
		OutputToConsole: cfg.Logging.OutputToConsole,
		OutputToFile:    cfg.Logging.OutputToFile,
		OutputFilePath:  cfg.Logging.OutputFilePath,
		UseColorLevel:   cfg.Logging.UseColorLevel,
		MaxSizeMB:       cfg.Logging.MaxSizeMB,
		MaxBackups:      cfg.Logging.MaxBackups,
		MaxAgeDays:      cfg.Logging.MaxAgeDays,
		Compress:        cfg.Logging.Compress,
		CallerSkip:      cfg.Logging.CallerSkip,
		AddStacktrace:   cfg.Logging.AddStacktrace,
	})
}

// InfrastructureModule 提供基础设施相关的依赖项
var InfrastructureModule = fx.Options(
	fx.Provide(
		config.GetConfig,
		NewLogger,
	),
)
