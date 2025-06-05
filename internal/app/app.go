package app

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"nebulaLive/internal/config"
	"nebulaLive/internal/di"
	"nebulaLive/internal/entity/ent"
	"nebulaLive/pkg/logger"
)

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
		fx.Invoke(func(app *fiber.App, cfg *config.Config, logger *logger.Logger) {
			logger.GetZapLogger().Info("Starting server", zap.String("addr", cfg.Server.Addr))
			if err := app.Listen(cfg.Server.Addr); err != nil {
				logger.GetZapLogger().Fatal("Failed to start server", zap.Error(err))
			}
		}),
		fx.Invoke(func(lifecycle fx.Lifecycle, logger *logger.Logger, client *ent.Client) {
			lifecycle.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					fmt.Println("Stopping server")
					logger.GetZapLogger().Info("Closing resources")
					errs := make([]error, 0)
					if err := client.Close(); err != nil {
						errs = append(errs, err)
						logger.GetZapLogger().Fatal("Failed to close client", zap.Error(err))
					}
					if err := logger.Sync(); err != nil {
						errs = append(errs, err)
						logger.GetZapLogger().Fatal("Failed to sync logger", zap.Error(err))
					}
					if len(errs) != 0 {
						return fmt.Errorf("failed to close resources: %v", errs)
					}
					return nil
				},
			})
		}),
	)
}
