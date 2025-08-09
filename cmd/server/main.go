package main

import (
	"context"

	"nebula-live/ent"
	"nebula-live/internal/app"
	"nebula-live/internal/domain/service"
	"nebula-live/internal/infrastructure"
	"nebula-live/internal/infrastructure/persistence"
	"nebula-live/internal/infrastructure/web/handler"
	"nebula-live/internal/infrastructure/web/middleware"
	"nebula-live/internal/infrastructure/web/router"
	"nebula-live/pkg/logger"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fxApp := fx.New(
		// 基础设施层模块
		infrastructure.InfrastructureModule,
		
		// 仓储层模块
		persistence.PersistenceModule,
		
		// 服务层模块
		service.ServiceModule,
		
		// 中间件模块
		middleware.MiddlewareModule,
		
		// 处理器层模块
		handler.HandlerModule,
		
		// 路由模块
		router.RouterModule,
		
		// 应用层模块
		app.AppModule,
		fx.Invoke(func(lc fx.Lifecycle, server *app.Server, client *ent.Client, zapLogger *zap.Logger) {
			// 初始化全局logger
			logger.Initialize(zapLogger)
			
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					// 运行数据库迁移
					if err := persistence.RunMigrations(ctx, client, zapLogger); err != nil {
						zapLogger.Error("Failed to run migrations", zap.Error(err))
						return err
					}
					
					logger.Info("Starting nebula-live server")
					go func() {
						if err := server.Start(); err != nil {
							logger.Error("Server start error", zap.Error(err))
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.Info("Stopping nebula-live server")
					if err := server.Stop(); err != nil {
						logger.Error("Error stopping server", zap.Error(err))
					}
					
					// 关闭数据库连接
					if err := persistence.CloseEntClient(client, zapLogger); err != nil {
						logger.Error("Error closing database connection", zap.Error(err))
						return err
					}
					
					return nil
				},
			})
		}),
	)

	fxApp.Run()
}