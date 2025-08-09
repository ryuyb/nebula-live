package logger

import (
	"go.uber.org/zap"
)

var (
	// Logger 全局logger实例
	Logger *zap.Logger
)

// Initialize 初始化全局logger
func Initialize(logger *zap.Logger) {
	Logger = logger
}

// Info 信息级别日志
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Error 错误级别日志
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Warn 警告级别日志
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Debug 调试级别日志
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}
