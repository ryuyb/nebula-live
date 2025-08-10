package app

import (
	"fmt"

	"nebula-live/internal/infrastructure/config"
	"nebula-live/internal/infrastructure/web/middleware"
	"nebula-live/internal/infrastructure/web/router"
	"nebula-live/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"go.uber.org/zap"

	_ "nebula-live/docs" // swagger docs
)

type Server struct {
	app    *fiber.App
	config *config.Config
	logger *zap.Logger
}

func NewFiberApp(cfg *config.Config, log *zap.Logger, routerRegistry *router.RouterRegistry) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal server error"

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}

			log.Error("Request failed",
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.Int("status", code),
				zap.Error(err),
			)

			return c.Status(code).JSON(errors.NewAPIError(code, "Request failed", message))
		},
	})

	// 全局中间件
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(middleware.ZapLogger(log))

	// CORS 配置
	app.Use(cors.New(cors.Config{
		AllowOrigins:     joinStrings(cfg.CORS.AllowedOrigins),
		AllowMethods:     joinStrings(cfg.CORS.AllowedMethods),
		AllowHeaders:     joinStrings(cfg.CORS.AllowedHeaders),
		ExposeHeaders:    joinStrings(cfg.CORS.ExposedHeaders),
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           cfg.CORS.MaxAge,
	}))

	// 健康检查
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": cfg.App.Name,
			"version": cfg.App.Version,
		})
	})

	// Swagger API 文档
	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html", fiber.StatusMovedPermanently)
	})
	app.Get("/swagger/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html", fiber.StatusMovedPermanently)
	})
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// 设置路由
	routerRegistry.RegisterAllRoutes(app)

	return &Server{
		app:    app,
		config: cfg,
		logger: log,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	s.logger.Info("Server starting", zap.String("address", addr))
	return s.app.Listen(addr)
}

func (s *Server) Stop() error {
	s.logger.Info("Server stopping")
	return s.app.Shutdown()
}

func joinStrings(slice []string) string {
	if len(slice) == 0 {
		return ""
	}
	result := slice[0]
	for i := 1; i < len(slice); i++ {
		result += "," + slice[i]
	}
	return result
}
