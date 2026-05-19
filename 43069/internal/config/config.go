package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apitester/apitester/pkg/models"
	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*models.TestConfig, error) {
	return loadConfig(path, false)
}

func loadConfig(path string, isInherited bool) (*models.TestConfig, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid config path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config models.TestConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if config.Inherits != "" {
		inheritedPath := filepath.Join(filepath.Dir(absPath), config.Inherits)
		parentConfig, err := loadConfig(inheritedPath, true)
		if err != nil {
			return nil, fmt.Errorf("failed to load inherited config: %w", err)
		}
		mergeConfigs(parentConfig, &config)
		config = *parentConfig
	}

	applyEnvOverrides(&config)
	setDefaults(&config)

	if isInherited {
		if err := validateConfigLight(&config); err != nil {
			return nil, err
		}
	} else {
		if err := validateConfig(&config); err != nil {
			return nil, err
		}
	}

	return &config, nil
}

func mergeConfigs(parent, child *models.TestConfig) {
	if child.Name != "" {
		parent.Name = child.Name
	}
	if child.Description != "" {
		parent.Description = child.Description
	}
	if child.BaseURL != "" {
		parent.BaseURL = child.BaseURL
	}
	if child.Timeout > 0 {
		parent.Timeout = child.Timeout
	}
	if child.Retries > 0 {
		parent.Retries = child.Retries
	}

	if child.Variables != nil {
		if parent.Variables == nil {
			parent.Variables = make(map[string]string)
		}
		for k, v := range child.Variables {
			parent.Variables[k] = v
		}
	}

	if child.Headers != nil {
		if parent.Headers == nil {
			parent.Headers = make(map[string]string)
		}
		for k, v := range child.Headers {
			parent.Headers[k] = v
		}
	}

	if child.Concurrency.Mode != "" {
		parent.Concurrency.Mode = child.Concurrency.Mode
	}
	if child.Concurrency.Workers > 0 {
		parent.Concurrency.Workers = child.Concurrency.Workers
	}
	if child.Concurrency.TotalRequests > 0 {
		parent.Concurrency.TotalRequests = child.Concurrency.TotalRequests
	}
	if child.Concurrency.Duration > 0 {
		parent.Concurrency.Duration = child.Concurrency.Duration
	}
	if child.Concurrency.RampUp > 0 {
		parent.Concurrency.RampUp = child.Concurrency.RampUp
	}
	if len(child.Concurrency.Steps) > 0 {
		parent.Concurrency.Steps = child.Concurrency.Steps
	}

	if len(child.Scenarios) > 0 {
		parent.Scenarios = child.Scenarios
	}
}

func applyEnvOverrides(config *models.TestConfig) {
	for key, val := range config.Variables {
		if strings.HasPrefix(val, "${") && strings.HasSuffix(val, "}") {
			envKey := strings.TrimSuffix(strings.TrimPrefix(val, "${"), "}")
			if envVal := os.Getenv(envKey); envVal != "" {
				config.Variables[key] = envVal
			}
		}
	}

	if baseURL := os.Getenv("APITESTER_BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}
}

func setDefaults(config *models.TestConfig) {
	if config.Timeout == 0 {
		config.Timeout = 30
	}
	if config.Retries == 0 {
		config.Retries = 0
	}
	if config.Concurrency.Mode == "" {
		config.Concurrency.Mode = "workers"
	}
	if config.Concurrency.Workers == 0 {
		config.Concurrency.Workers = 10
	}
	if config.Concurrency.Duration == 0 {
		config.Concurrency.Duration = 60
	}

	for i := range config.Scenarios {
		if config.Scenarios[i].Weight == 0 {
			config.Scenarios[i].Weight = 1
		}
		for j := range config.Scenarios[i].Requests {
			req := &config.Scenarios[i].Requests[j]
			if req.Method == "" {
				req.Method = "GET"
			}
			if req.Timeout == 0 {
				req.Timeout = config.Timeout
			}
			if req.Retries == 0 {
				req.Retries = config.Retries
			}
		}
	}
}

func validateConfig(config *models.TestConfig) error {
	if err := validateConfigLight(config); err != nil {
		return err
	}
	if len(config.Scenarios) == 0 {
		return fmt.Errorf("at least one scenario is required")
	}
	for _, scenario := range config.Scenarios {
		if scenario.Name == "" {
			return fmt.Errorf("all scenarios must have a name")
		}
		if len(scenario.Requests) == 0 {
			return fmt.Errorf("scenario '%s' has no requests", scenario.Name)
		}
		for _, req := range scenario.Requests {
			if req.Path == "" {
				return fmt.Errorf("request in scenario '%s' has no path", scenario.Name)
			}
		}
	}
	return nil
}

func validateConfigLight(config *models.TestConfig) error {
	if config.BaseURL == "" {
		return fmt.Errorf("base_url is required")
	}
	if config.Concurrency.Mode != "" && config.Concurrency.Mode != "workers" && config.Concurrency.Mode != "steps" && config.Concurrency.Mode != "duration" {
		return fmt.Errorf("invalid concurrency mode: %s (must be workers, steps, or duration)", config.Concurrency.Mode)
	}
	return nil
}
