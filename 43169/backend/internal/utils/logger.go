package utils

import (
	"io"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	file   *lumberjack.Logger
	stdout io.Writer
	writer io.Writer
}

func NewLogger(filename string) *Logger {
	l := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
	}
	return &Logger{
		file:   l,
		stdout: os.Stdout,
		writer: io.MultiWriter(l, os.Stdout),
	}
}

func (l *Logger) Write(p []byte) (n int, err error) {
	return l.writer.Write(p)
}

func (l *Logger) Close() error {
	return l.file.Close()
}

func LogFormat() string {
	return `{"time":"` + time.RFC3339 + `","level":"{{.Level}}","msg":{{.Message}}}` + "\n"
}
