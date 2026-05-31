package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger(mode string) (*zap.Logger, error) {
	var config zap.Config
	if mode == "release" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var err error
	Logger, err = config.Build()
	if err != nil {
		return nil, err
	}

	return Logger, nil
}

func SyncLogger() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}
