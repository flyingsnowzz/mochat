package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log 全局 zap.Logger 实例，用于结构化日志输出
var Log *zap.Logger

// Sugar 全局 zap.SugaredLogger 实例，用于格式化日志输出
var Sugar *zap.SugaredLogger

// Init 初始化日志系统
// level: 日志级别 (debug, info, warn, error, dpanic, panic, fatal)
// output: 日志输出路径 (stdout, stderr, 或文件路径)
func Init(level string, output string) error {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{output},
		ErrorOutputPaths: []string{output},
	}

	var err error
	Log, err = config.Build()
	if err != nil {
		return err
	}
	Sugar = Log.Sugar()
	return nil
}
