package middleware

import (
	"github.com/samber/lo"
	"go.uber.org/zap"
	"nebulaLive/internal/entity"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// ErrorHandler 是一个自定义的错误处理器，使用通用的Response结构
func ErrorHandler(c *fiber.Ctx, err error) error {
	// 默认错误代码和消息
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// 检查是否为Fiber错误
	if fiberErr, ok := lo.ErrorsAs[*fiber.Error](err); ok {
		code = fiberErr.Code
		message = fiberErr.Message
	}

	// 当错误代码为500时，记录错误日志
	if code == fiber.StatusInternalServerError {
		log.Errorw("Internal Server Error: ", zap.Error(err))
	}

	// 返回格式化的错误响应
	return c.Status(code).JSON(entity.ErrorResponse(code, message))
}
