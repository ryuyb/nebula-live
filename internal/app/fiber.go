package app

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"nebulaLive/internal/api/middleware"
	"nebulaLive/internal/api/router"
	"nebulaLive/pkg/logger"
)

func NewFiberApp(registry *router.RouterRegistry, logger *logger.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "Nebula",
		ErrorHandler: middleware.ErrorHandler,
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
	})

	log.SetLogger(fiberzap.NewLogger(fiberzap.LoggerConfig{
		SetLogger: logger.GetZapLogger(),
	}))

	middleware.Common(app, logger.GetZapLogger())

	// 设置 API 路由
	router.SetupRouter(app, registry)

	middleware.NotFound(app)

	return app
}
