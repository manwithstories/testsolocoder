package scenario

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/apitester/apitester/internal/assertion"
	"github.com/apitester/apitester/internal/httpclient"
	"github.com/apitester/apitester/internal/logger"
	"github.com/apitester/apitester/pkg/models"
)

type ScenarioExecutor struct {
	config   *models.TestConfig
	scenario *models.Scenario
	client   *httpclient.Client
	vars     map[string]string
}

func NewScenarioExecutor(config *models.TestConfig, scenario *models.Scenario) *ScenarioExecutor {
	vars := make(map[string]string)
	for k, v := range config.Variables {
		vars[k] = v
	}
	for k, v := range scenario.Variables {
		vars[k] = v
	}

	return &ScenarioExecutor{
		config:   config,
		scenario: scenario,
		client:   httpclient.NewClient(time.Duration(config.Timeout) * time.Second),
		vars:     vars,
	}
}

func (se *ScenarioExecutor) Execute(ctx context.Context, workerID int) []*models.RequestResult {
	var results []*models.RequestResult

	for i, req := range se.scenario.Requests {
		select {
		case <-ctx.Done():
			return results
		default:
		}

		result := se.executeRequest(ctx, &req, workerID)
		results = append(results, result)

		if !result.Success {
			logger.Debug("Scenario '%s' stopped at request %d due to failure", se.scenario.Name, i+1)
			break
		}

		if req.ThinkTime > 0 {
			select {
			case <-ctx.Done():
				return results
			case <-time.After(time.Duration(req.ThinkTime) * time.Millisecond):
			}
		}
	}

	return results
}

func (se *ScenarioExecutor) executeRequest(ctx context.Context, req *models.Request, workerID int) *models.RequestResult {
	startTime := time.Now()

	method := req.Method
	url := se.config.BaseURL + se.renderTemplate(req.Path, se.vars)

	headers := make(map[string]string)
	for k, v := range se.config.Headers {
		headers[k] = se.renderTemplate(v, se.vars)
	}
	for k, v := range se.scenario.Headers {
		headers[k] = se.renderTemplate(v, se.vars)
	}
	for k, v := range req.Headers {
		headers[k] = se.renderTemplate(v, se.vars)
	}

	for k, v := range req.Variables {
		se.vars[k] = se.renderTemplate(v, se.vars)
	}

	body := se.renderTemplate(req.Body, se.vars)

	logger.Debug("[Worker %d] Executing request: %s %s", workerID, method, url)

	resp := se.client.Do(ctx, method, url, headers, body, req.Retries)

	result := &models.RequestResult{
		ScenarioName: se.scenario.Name,
		RequestName:  req.Name,
		Method:       method,
		URL:          url,
		StartTime:    startTime,
		EndTime:      time.Now(),
		Duration:     resp.Duration,
		StatusCode:   resp.StatusCode,
		Headers:      resp.Headers,
		Body:         resp.Body,
		Error:        resp.Error,
		Success:      resp.Error == nil,
	}

	if resp.Error == nil {
		assertionErrs := assertion.RunAssertions(req, &assertion.ResponseData{
			StatusCode: resp.StatusCode,
			Duration:   resp.Duration,
			Headers:    resp.Headers,
			Body:       resp.Body,
		})
		result.AssertionErrs = assertionErrs
		if len(assertionErrs) > 0 {
			result.Success = false
			logger.Debug("Assertion failures for %s: %v", req.Name, assertionErrs)
		}

		for _, extract := range req.Extract {
			if value, err := se.extractValue(extract, resp); err == nil {
				se.vars[extract.Name] = value
				logger.Debug("Extracted variable %s = %s", extract.Name, value)
			} else {
				logger.Warn("Failed to extract %s: %v", extract.Name, err)
			}
		}
	} else {
		logger.Debug("Request failed: %v", resp.Error)
	}

	entry := &models.LogEntry{
		Timestamp:     time.Now(),
		Level:         "INFO",
		ScenarioName:  se.scenario.Name,
		RequestName:   req.Name,
		Method:        method,
		URL:           url,
		Duration:      resp.Duration,
		StatusCode:    resp.StatusCode,
		WorkerID:      workerID,
	}
	if resp.Error != nil {
		entry.Error = resp.Error.Error()
		entry.Level = "ERROR"
	}
	for _, ae := range result.AssertionErrs {
		entry.AssertionErrs = append(entry.AssertionErrs, ae.Error())
	}
	logger.LogRequest(entry)

	return result
}

func (se *ScenarioExecutor) renderTemplate(tpl string, vars map[string]string) string {
	if !strings.Contains(tpl, "{{") {
		return tpl
	}

	result := tpl
	for k, v := range vars {
		placeholder := "{{" + k + "}}"
		result = strings.ReplaceAll(result, placeholder, v)
	}

	if strings.Contains(result, "{{") {
		logger.Debug("Some variables could not be replaced in template: %s", result)
	}

	return result
}

func (se *ScenarioExecutor) extractValue(rule models.ExtractRule, resp *httpclient.Response) (string, error) {
	switch rule.From {
	case "body":
		val, err := assertion.GetJSONValue(resp.Body, rule.Path)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%v", val), nil
	case "header":
		if val, ok := resp.Headers[rule.Path]; ok {
			return val, nil
		}
		return "", fmt.Errorf("header not found: %s", rule.Path)
	case "cookie":
		if cookieHeader, ok := resp.Headers["Set-Cookie"]; ok {
			for _, cookie := range strings.Split(cookieHeader, ";") {
				parts := strings.SplitN(strings.TrimSpace(cookie), "=", 2)
				if len(parts) == 2 && parts[0] == rule.Path {
					return parts[1], nil
				}
			}
		}
		return "", fmt.Errorf("cookie not found: %s", rule.Path)
	default:
		return "", fmt.Errorf("unknown extract source: %s", rule.From)
	}
}

func GetScenariosByWeight(config *models.TestConfig) []*models.Scenario {
	var weightedScenarios []*models.Scenario
	for i := range config.Scenarios {
		scenario := &config.Scenarios[i]
		weight := scenario.Weight
		if weight <= 0 {
			weight = 1
		}
		for w := 0; w < weight; w++ {
			weightedScenarios = append(weightedScenarios, scenario)
		}
	}
	return weightedScenarios
}
