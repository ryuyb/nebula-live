package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ZapLogger 创建使用zap的Fiber日志中间件
func ZapLogger(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		
		// 继续处理请求
		err := c.Next()
		
		// 记录请求日志
		duration := time.Since(start)
		status := c.Response().StatusCode()
		
		fields := []zap.Field{
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", status),
			zap.Duration("latency", duration),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
		}
		
		// 添加请求ID如果存在
		if requestID := c.GetRespHeader("X-Request-ID"); requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}
		
		// 根据状态码选择日志级别
		if status >= 500 {
			logger.Error("HTTP request", fields...)
		} else if status >= 400 {
			logger.Warn("HTTP request", fields...)
		} else {
			logger.Info("HTTP request", fields...)
		}
		
		return err
	}
}