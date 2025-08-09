package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"nebula-live/internal/infrastructure/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	// 只有在需要输出到文件时才创建日志目录
	if cfg.Log.EnableFile {
		logDir := filepath.Dir(cfg.Log.Output)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
	}

	// 配置日志级别
	var level zapcore.Level
	switch cfg.Log.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	// 配置编码器
	var encoderConfig zapcore.EncoderConfig
	if cfg.Log.Format == "json" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	
	// 配置彩色级别编码
	if cfg.Log.EnableColor {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	// 文件写入器（使用 lumberjack 实现日志轮转）
	fileWriter := &lumberjack.Logger{
		Filename:   cfg.Log.Output,
		MaxSize:    cfg.Log.MaxSize,
		MaxAge:     cfg.Log.MaxAge,
		MaxBackups: cfg.Log.MaxBackups,
		Compress:   cfg.Log.Compress,
	}

	// 控制台写入器
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 根据配置决定输出目标
	var cores []zapcore.Core
	
	// 控制台输出
	if cfg.Log.EnableConsole {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(consoleEncoder, consoleWriter, level))
	}
	
	// 文件输出
	if cfg.Log.EnableFile {
		// 文件编码器始终使用普通级别编码器
		fileEncoderConfig := encoderConfig
		fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(fileWriter), level))
	}
	
	// 创建core
	var core zapcore.Core
	if len(cores) == 0 {
		// 如果没有启用任何输出，默认输出到控制台
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		core = zapcore.NewCore(consoleEncoder, consoleWriter, level)
	} else if len(cores) == 1 {
		core = cores[0]
	} else {
		core = zapcore.NewTee(cores...)
	}

	// 创建 logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return logger, nil
}