package models

import (
	"time"
)

type TestConfig struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	BaseURL     string            `yaml:"base_url"`
	Variables   map[string]string `yaml:"variables"`
	Headers     map[string]string `yaml:"headers"`
	Timeout     int               `yaml:"timeout"`
	Retries     int               `yaml:"retries"`
	Concurrency ConcurrencyConfig `yaml:"concurrency"`
	Scenarios   []Scenario        `yaml:"scenarios"`
	Inherits    string            `yaml:"inherits,omitempty"`
}

type ConcurrencyConfig struct {
	Mode          string `yaml:"mode"`
	Workers       int    `yaml:"workers"`
	TotalRequests int    `yaml:"total_requests"`
	Duration      int    `yaml:"duration"`
	RampUp        int    `yaml:"ramp_up"`
	Steps         []Step `yaml:"steps,omitempty"`
}

type Step struct {
	Workers  int `yaml:"workers"`
	Duration int `yaml:"duration"`
}

type Scenario struct {
	Name        string            `yaml:"name"`
	Weight      int               `yaml:"weight"`
	Requests    []Request         `yaml:"requests"`
	Variables   map[string]string `yaml:"variables"`
	Headers     map[string]string `yaml:"headers"`
}

type Request struct {
	Name        string            `yaml:"name"`
	Method      string            `yaml:"method"`
	Path        string            `yaml:"path"`
	Headers     map[string]string `yaml:"headers"`
	Body        string            `yaml:"body"`
	Variables   map[string]string `yaml:"variables"`
	Extract     []ExtractRule     `yaml:"extract"`
	Assertions  []AssertionRule   `yaml:"assertions"`
	Timeout     int               `yaml:"timeout"`
	Retries     int               `yaml:"retries"`
	ThinkTime   int               `yaml:"think_time"`
}

type ExtractRule struct {
	Name string `yaml:"name"`
	From string `yaml:"from"`
	Path string `yaml:"path"`
}

type AssertionRule struct {
	Type      string      `yaml:"type"`
	Operator  string      `yaml:"operator"`
	Value     interface{} `yaml:"value"`
	Field     string      `yaml:"field,omitempty"`
	Condition string      `yaml:"condition,omitempty"`
}

type RequestResult struct {
	ScenarioName  string
	RequestName   string
	Method        string
	URL           string
	StartTime     time.Time
	EndTime       time.Time
	Duration      time.Duration
	StatusCode    int
	ResponseSize  int64
	Error         error
	AssertionErrs []error
	Success       bool
	Headers       map[string]string
	Body          string
}

type TestStats struct {
	TotalRequests    int64
	SuccessRequests  int64
	FailedRequests   int64
	TotalDuration    time.Duration
	StartTime        time.Time
	EndTime          time.Time
	ResponseTimes    []time.Duration
	StatusCodes      map[int]int64
	Errors           map[string]int64
	AssertionErrors  map[string]int64
	CurrentQPS       float64
	CurrentWorkers   int
}

type Report struct {
	TestName      string
	StartTime     time.Time
	EndTime       time.Time
	Duration      time.Duration
	TotalRequests int64
	SuccessCount  int64
	FailedCount   int64
	ErrorRate     float64
	QPS           float64
	AvgRT         time.Duration
	MinRT         time.Duration
	MaxRT         time.Duration
	P50           time.Duration
	P90           time.Duration
	P95           time.Duration
	P99           time.Duration
	StatusCodes   map[int]int64
	TopErrors     map[string]int64
	ScenarioStats map[string]*ScenarioStat
}

type ScenarioStat struct {
	Name          string
	TotalRequests int64
	SuccessCount  int64
	FailedCount   int64
	ErrorRate     float64
	AvgRT         time.Duration
	MinRT         time.Duration
	MaxRT         time.Duration
}

type LogEntry struct {
	Timestamp     time.Time
	Level         string
	ScenarioName  string
	RequestName   string
	Method        string
	URL           string
	Duration      time.Duration
	StatusCode    int
	Error         string
	AssertionErrs []string
	WorkerID      int
}
