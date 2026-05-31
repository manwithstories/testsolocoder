package utils

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/natefinsh/lumberjack.v2"
)

type LoggerConfig struct {
	Path       string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

func InitLogger(cfg LoggerConfig) (io.Writer, error) {
	if err := os.MkdirAll(cfg.Path, 0755); err != nil {
		return nil, err
	}

	logFile := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Path, "app.log"),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   true,
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	return multiWriter, nil
}
