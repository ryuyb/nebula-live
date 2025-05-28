package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

// Logger 是一个封装了 zap.Logger 的结构体
type Logger struct {
	zapLogger *zap.Logger
}

// LoggerConfig 定义了日志配置选项
type LoggerConfig struct {
	Level           string // 日志级别
	OutputToConsole bool   // 是否输出到控制台
	OutputToFile    bool   // 是否输出到文件
	OutputFilePath  string // 输出文件路径
	UseColorLevel   bool   // 是否使用彩色级别编码
	MaxSizeMB       int    // 日志文件最大大小（MB）
	MaxBackups      int    // 保留的最大备份文件数
	MaxAgeDays      int    // 保留的最大天数
	Compress        bool   // 是否压缩备份文件
	CallerSkip      int    // 跳过调用者的层级数
	AddStacktrace   bool   // 是否添加堆栈跟踪
}

// NewLogger 创建一个新的 Logger 实例，使用自定义配置结构体
func NewLogger(cfg LoggerConfig) (*Logger, error) {
	// 使用 zap 的生产环境配置作为基础
	config := zap.NewProductionConfig()
	// 设置日志级别
	var level zap.AtomicLevel
	switch cfg.Level {
	case "debug":
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	config.Level = level
	// 设置默认的时间编码器为 ISO8601 格式
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 设置级别编码器
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 设置输出路径和 lumberjack 控制
	var syncers []zapcore.WriteSyncer
	if cfg.OutputToConsole {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}
	if cfg.OutputToFile && cfg.OutputFilePath != "" {
		lumberjackLogger := &lumberjack.Logger{
			Filename:   cfg.OutputFilePath,
			MaxSize:    cfg.MaxSizeMB,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   cfg.Compress,
		}
		syncers = append(syncers, zapcore.AddSync(lumberjackLogger))
	}
	// 如果都没有设置输出目标，则默认输出到控制台
	if len(syncers) == 0 {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}
	// 构建 zap logger，使用自定义的 writeSyncer
	var cores []zapcore.Core
	if cfg.OutputToConsole {
		// 控制台输出使用彩色级别编码（如果启用）
		consoleEncoderConfig := config.EncoderConfig
		if cfg.UseColorLevel {
			consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), config.Level))
	}
	if cfg.OutputToFile && cfg.OutputFilePath != "" {
		// 文件输出不使用彩色级别编码
		fileEncoder := zapcore.NewJSONEncoder(config.EncoderConfig)
		lumberjackLogger := &lumberjack.Logger{
			Filename:   cfg.OutputFilePath,
			MaxSize:    cfg.MaxSizeMB,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   cfg.Compress,
		}
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberjackLogger), config.Level))
	}
	// 如果都没有设置输出目标，则默认输出到控制台
	if len(cores) == 0 {
		consoleEncoderConfig := config.EncoderConfig
		if cfg.UseColorLevel {
			consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), config.Level))
	}
	// 组合多个核心
	var zapLogger *zap.Logger
	if len(cores) > 1 {
		zapLogger = zap.New(zapcore.NewTee(cores...), zap.WithCaller(true), zap.AddCallerSkip(cfg.CallerSkip))
	} else {
		zapLogger = zap.New(cores[0], zap.WithCaller(true), zap.AddCallerSkip(cfg.CallerSkip))
	}
	if cfg.AddStacktrace {
		zapLogger = zapLogger.WithOptions(zap.AddStacktrace(zapcore.ErrorLevel))
	}
	return &Logger{zapLogger: zapLogger}, nil
}

// NewLoggerWithDefaults 创建一个具有默认配置的 Logger 实例
func NewLoggerWithDefaults() (*Logger, error) {
	return NewLogger(LoggerConfig{
		Level:           "info",
		OutputToConsole: true, // 默认输出到控制台
		OutputToFile:    false,
		OutputFilePath:  "",
		UseColorLevel:   true,
		MaxSizeMB:       10,    // 默认最大大小 10MB
		MaxBackups:      5,     // 默认保留 5 个备份文件
		MaxAgeDays:      30,    // 默认保留 30 天
		Compress:        true,  // 默认启用压缩
		CallerSkip:      0,     // 默认不跳过调用者
		AddStacktrace:   false, // 默认不添加堆栈跟踪
	})
}

// GetZapLogger 返回内部的 zap.Logger 实例
func (l *Logger) GetZapLogger() *zap.Logger {
	return l.zapLogger
}

// Sync 调用内部 zap.Logger 的 Sync 方法，确保日志被写入
func (l *Logger) Sync() error {
	return l.zapLogger.Sync()
}
