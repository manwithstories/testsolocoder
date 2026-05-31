package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"tea-platform/config"
)

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
	once   sync.Once
)

func Init() {
	once.Do(func() {
		cfg := config.Get()

		level := parseLevel(cfg.Log.Level)
		encoder := parseEncoder(cfg.Log.Format)

		var core zapcore.Core
		switch cfg.Log.Output {
		case "stdout":
			core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
		case "stderr":
			core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), level)
		default:
			writer, err := os.OpenFile(cfg.Log.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
			} else {
				core = zapcore.NewCore(encoder, zapcore.AddSync(writer), level)
			}
		}

		Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		Sugar = Logger.Sugar()
	})
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func parseEncoder(format string) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}
