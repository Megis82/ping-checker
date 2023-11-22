package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey: "ts",
		// LevelKey:       "level",
		// NameKey:        "logger",
		// CallerKey:      "caller",
		FunctionKey: zapcore.OmitKey,
		MessageKey:  "msg",
		// StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		// EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func NewLogger(LogFileName string) (*zap.Logger, error) {
	cfg := NewProductionEncoderConfig()
	// cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths:      []string{LogFileName, "stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    cfg,
	}
	logger, err := config.Build()
	return logger, err
}
