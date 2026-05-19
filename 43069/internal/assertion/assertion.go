package assertion

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/apitester/apitester/pkg/models"
)

type AssertionResult struct {
	Passed bool
	Error  error
}

func RunAssertions(req *models.Request, resp *ResponseData) []error {
	var errors []error

	for _, assertion := range req.Assertions {
		result := checkAssertion(assertion, resp)
		if !result.Passed {
			errors = append(errors, result.Error)
		}
	}

	return errors
}

func checkAssertion(rule models.AssertionRule, resp *ResponseData) AssertionResult {
	switch rule.Type {
	case "status_code":
		return checkStatusCode(rule, resp.StatusCode)
	case "response_time":
		return checkResponseTime(rule, resp.Duration)
	case "body_field":
		return checkBodyField(rule, resp.Body)
	case "header":
		return checkHeader(rule, resp.Headers)
	case "body_contains":
		return checkBodyContains(rule, resp.Body)
	default:
		return AssertionResult{
			Passed: false,
			Error:  fmt.Errorf("unknown assertion type: %s", rule.Type),
		}
	}
}

func checkStatusCode(rule models.AssertionRule, actual int) AssertionResult {
	expected, err := toInt(rule.Value)
	if err != nil {
		return AssertionResult{Passed: false, Error: fmt.Errorf("invalid status code value: %v", rule.Value)}
	}

	passed, err := compareInt(rule.Operator, actual, expected)
	if err != nil {
		return AssertionResult{Passed: false, Error: err}
	}

	if !passed {
		return AssertionResult{
			Passed: false,
			Error:  fmt.Errorf("status code assertion failed: %d %s %d", actual, rule.Operator, expected),
		}
	}
	return AssertionResult{Passed: true}
}

func checkResponseTime(rule models.AssertionRule, actual time.Duration) AssertionResult {
	expectedMs, err := toInt(rule.Value)
	if err != nil {
		return AssertionResult{Passed: false, Error: fmt.Errorf("invalid response time value: %v", rule.Value)}
	}

	actualMs := int(actual.Milliseconds())
	passed, err := compareInt(rule.Operator, actualMs, expectedMs)
	if err != nil {
		return AssertionResult{Passed: false, Error: err}
	}

	if !passed {
		return AssertionResult{
			Passed: false,
			Error:  fmt.Errorf("response time assertion failed: %dms %s %dms", actualMs, rule.Operator, expectedMs),
		}
	}
	return AssertionResult{Passed: true}
}

func checkBodyField(rule models.AssertionRule, body string) AssertionResult {
	if rule.Field == "" {
		return AssertionResult{Passed: false, Error: fmt.Errorf("field is required for body_field assertion")}
	}

	actual, err := GetJSONValue(body, rule.Field)
	if err != nil {
		return AssertionResult{Passed: false, Error: fmt.Errorf("failed to extract field %s: %v", rule.Field, err)}
	}

	passed, err := compareValues(rule.Operator, actual, rule.Value)
	if err != nil {
		return AssertionResult{Passed: false, Error: err}
	}

	if !passed {
		return AssertionResult{
			Passed: false,
			Error:  fmt.Errorf("body field assertion failed: %s=%v %s %v", rule.Field, actual, rule.Operator, rule.Value),
		}
	}
	return AssertionResult{Passed: true}
}

func checkHeader(rule models.AssertionRule, headers map[string]string) AssertionResult {
	if rule.Field == "" {
		return AssertionResult{Passed: false, Error: fmt.Errorf("field is required for header assertion")}
	}

	actual, exists := headers[rule.Field]
	if !exists {
		return AssertionResult{
			Passed: false,
			Error:  fmt.Errorf("header %s not found", rule.Field),
		}
	}

	passed, err := compareStrings(rule.Operator, actual, fmt.Sprintf("%v", rule.Value))
	if err != nil {
		return AssertionResult{Passed: false, Error: err}
	}

	if !passed {
		return AssertionResult{
			Passed: false,
			Error:  fmt.Errorf("header assertion failed: %s=%s %s %v", rule.Field, actual, rule.Operator, rule.Value),
		}
	}
	return AssertionResult{Passed: true}
}

func checkBodyContains(rule models.AssertionRule, body string) AssertionResult {
	expected := fmt.Sprintf("%v", rule.Value)
	if !strings.Contains(body, expected) {
		return AssertionResult{
			Passed: false,
			Error:  fmt.Errorf("body does not contain: %s", expected),
		}
	}
	return AssertionResult{Passed: true}
}

func compareInt(op string, actual, expected int) (bool, error) {
	switch op {
	case "eq", "==":
		return actual == expected, nil
	case "ne", "!=":
		return actual != expected, nil
	case "gt", ">":
		return actual > expected, nil
	case "lt", "<":
		return actual < expected, nil
	case "gte", ">=":
		return actual >= expected, nil
	case "lte", "<=":
		return actual <= expected, nil
	default:
		return false, fmt.Errorf("unsupported operator: %s", op)
	}
}

func compareStrings(op string, actual, expected string) (bool, error) {
	switch op {
	case "eq", "==":
		return actual == expected, nil
	case "ne", "!=":
		return actual != expected, nil
	case "contains":
		return strings.Contains(actual, expected), nil
	case "not_contains":
		return !strings.Contains(actual, expected), nil
	case "starts_with":
		return strings.HasPrefix(actual, expected), nil
	case "ends_with":
		return strings.HasSuffix(actual, expected), nil
	default:
		return false, fmt.Errorf("unsupported string operator: %s", op)
	}
}

func compareValues(op string, actual, expected interface{}) (bool, error) {
	actualStr := fmt.Sprintf("%v", actual)
	expectedStr := fmt.Sprintf("%v", expected)

	if actualFloat, err := strconv.ParseFloat(actualStr, 64); err == nil {
		if expectedFloat, err := strconv.ParseFloat(expectedStr, 64); err == nil {
			switch op {
			case "eq", "==":
				return actualFloat == expectedFloat, nil
			case "ne", "!=":
				return actualFloat != expectedFloat, nil
			case "gt", ">":
				return actualFloat > expectedFloat, nil
			case "lt", "<":
				return actualFloat < expectedFloat, nil
			case "gte", ">=":
				return actualFloat >= expectedFloat, nil
			case "lte", "<=":
				return actualFloat <= expectedFloat, nil
			}
		}
	}

	return compareStrings(op, actualStr, expectedStr)
}

func toInt(v interface{}) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case float64:
		return int(val), nil
	case string:
		return strconv.Atoi(val)
	default:
		return 0, fmt.Errorf("cannot convert %v to int", reflect.TypeOf(v))
	}
}

func GetJSONValue(body, path string) (interface{}, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return nil, err
	}

	parts := strings.Split(path, ".")
	current := data

	for _, part := range parts {
		switch m := current.(type) {
		case map[string]interface{}:
			var ok bool
			current, ok = m[part]
			if !ok {
				return nil, fmt.Errorf("field not found: %s", part)
			}
		default:
			return nil, fmt.Errorf("cannot access field %s on non-object", part)
		}
	}

	return current, nil
}

type ResponseData struct {
	StatusCode int
	Duration   time.Duration
	Headers    map[string]string
	Body       string
}
