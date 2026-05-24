package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
	"business-registration-platform/config"
)

func InitLogger() (*os.File, error) {
	logDir := filepath.Dir(config.AppConfig.Log.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(config.AppConfig.Log.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func LogOperation(module, action string, userID, targetID *uint, targetType, content, result string) {
	logEntry := map[string]interface{}{
		"module":     module,
		"action":     action,
		"userId":     userID,
		"targetType": targetType,
		"targetId":   targetID,
		"content":    content,
		"result":     result,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	data, _ := yaml.Marshal(logEntry)
	fmt.Printf("[OPERATION] %s\n", string(data))
}
