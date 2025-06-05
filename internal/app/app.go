package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"nebulaLive/internal/config"
	"nebulaLive/internal/di"
	"nebulaLive/internal/repository"
	"nebulaLive/pkg/logger"
)

func register(lifecycle fx.Lifecycle, app *fiber.App, cfg *config.Config, logger *logger.Logger, client *repository.Client) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.GetZapLogger().Info("Starting server", zap.String("addr", cfg.Server.Addr))
			go func() {
				if err := app.Listen(cfg.Server.Addr); err != nil {
					logger.GetZapLogger().Fatal("Failed to start server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			logger.GetZapLogger().Info("Stopping server")
			if err := client.Close(); err != nil {
				logger.GetZapLogger().Fatal("Failed to close client", zap.Error(err))
			}
			if err := logger.Sync(); err != nil {
				logger.GetZapLogger().Fatal("Failed to sync logger", zap.Error(err))
			}
			return app.Shutdown()
		},
	})
}

func New() *fx.App {
	return fx.New(
		di.InfrastructureModule,
		di.RepositoryModule,
		di.ServiceModule,
		di.HandlerModule,
		di.RouterModule,
		fx.WithLogger(func(logger *logger.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger.GetZapLogger()}
		}),
		fx.Provide(
			NewFiberApp,
		),
		fx.Invoke(register),
	)
}
