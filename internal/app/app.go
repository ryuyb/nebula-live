package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
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
		fx.Provide(
			NewFiberApp,
		),
		fx.Invoke(func(app *fiber.App, cfg *config.Config, logger *logger.Logger) {
			// 启动服务器，监听在 3000 端口
			logger.GetZapLogger().Info("Starting server", zap.String("addr", cfg.Server.Addr))
			if err := app.Listen(cfg.Server.Addr); err != nil {
				logger.GetZapLogger().Fatal("Failed to start server", zap.Error(err))
			}
		}),
		fx.Invoke(func(lifecycle fx.Lifecycle, logger *logger.Logger, client *ent.Client) {
			lifecycle.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					if err := logger.Sync(); err != nil {
						return err
					}
					return client.Close()
				},
			})
		}),
	)
}
