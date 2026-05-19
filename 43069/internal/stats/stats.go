package stats

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/apitester/apitester/pkg/models"
)

type Collector struct {
	mu              sync.Mutex
	totalRequests   int64
	successRequests int64
	failedRequests  int64
	startTime       time.Time
	lastUpdateTime  time.Time
	lastRequestCount int64
	currentQPS      float64
	responseTimes   []time.Duration
	statusCodes     map[int]int64
	errors          map[string]int64
	assertionErrors map[string]int64
	scenarioStats   map[string]*scenarioStatInternal
	currentWorkers  int
}

type scenarioStatInternal struct {
	totalRequests  int64
	successCount   int64
	failedCount    int64
	responseTimes  []time.Duration
}

func NewCollector() *Collector {
	return &Collector{
		statusCodes:     make(map[int]int64),
		errors:          make(map[string]int64),
		assertionErrors: make(map[string]int64),
		scenarioStats:   make(map[string]*scenarioStatInternal),
		startTime:       time.Now(),
		lastUpdateTime:  time.Now(),
	}
}

func (c *Collector) Record(result *models.RequestResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	atomic.AddInt64(&c.totalRequests, 1)
	c.responseTimes = append(c.responseTimes, result.Duration)
	c.statusCodes[result.StatusCode]++

	if c.scenarioStats[result.ScenarioName] == nil {
		c.scenarioStats[result.ScenarioName] = &scenarioStatInternal{
			responseTimes: make([]time.Duration, 0),
		}
	}
	ss := c.scenarioStats[result.ScenarioName]
	ss.totalRequests++
	ss.responseTimes = append(ss.responseTimes, result.Duration)

	if result.Success && len(result.AssertionErrs) == 0 {
		atomic.AddInt64(&c.successRequests, 1)
		ss.successCount++
	} else {
		atomic.AddInt64(&c.failedRequests, 1)
		ss.failedCount++
		if result.Error != nil {
			errMsg := truncate(result.Error.Error(), 100)
			c.errors[errMsg]++
		}
		for _, ae := range result.AssertionErrs {
			c.assertionErrors[ae.Error()]++
		}
	}
}

func (c *Collector) UpdateQPS() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(c.lastUpdateTime).Seconds()
	if elapsed > 0 {
		requestsSince := atomic.LoadInt64(&c.totalRequests) - c.lastRequestCount
		c.currentQPS = float64(requestsSince) / elapsed
	}
	c.lastRequestCount = atomic.LoadInt64(&c.totalRequests)
	c.lastUpdateTime = now
}

func (c *Collector) GetSnapshot() models.TestStats {
	c.mu.Lock()
	defer c.mu.Unlock()

	statusCodesCopy := make(map[int]int64)
	for k, v := range c.statusCodes {
		statusCodesCopy[k] = v
	}

	errorsCopy := make(map[string]int64)
	for k, v := range c.errors {
		errorsCopy[k] = v
	}

	assertionErrorsCopy := make(map[string]int64)
	for k, v := range c.assertionErrors {
		assertionErrorsCopy[k] = v
	}

	return models.TestStats{
		TotalRequests:   atomic.LoadInt64(&c.totalRequests),
		SuccessRequests: atomic.LoadInt64(&c.successRequests),
		FailedRequests:  atomic.LoadInt64(&c.failedRequests),
		StartTime:       c.startTime,
		EndTime:         time.Now(),
		ResponseTimes:   append([]time.Duration{}, c.responseTimes...),
		StatusCodes:     statusCodesCopy,
		Errors:          errorsCopy,
		AssertionErrors: assertionErrorsCopy,
		CurrentQPS:      c.currentQPS,
		CurrentWorkers:  c.currentWorkers,
	}
}

func (c *Collector) SetCurrentWorkers(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.currentWorkers = n
}

func (c *Collector) GenerateReport(testName string) *models.Report {
	c.mu.Lock()
	defer c.mu.Unlock()

	endTime := time.Now()
	duration := endTime.Sub(c.startTime)

	total := atomic.LoadInt64(&c.totalRequests)
	success := atomic.LoadInt64(&c.successRequests)
	failed := atomic.LoadInt64(&c.failedRequests)

	errorRate := 0.0
	if total > 0 {
		errorRate = float64(failed) / float64(total) * 100
	}

	qps := 0.0
	if duration.Seconds() > 0 {
		qps = float64(total) / duration.Seconds()
	}

	sortedRT := append([]time.Duration{}, c.responseTimes...)
	sort.Slice(sortedRT, func(i, j int) bool {
		return sortedRT[i] < sortedRT[j]
	})

	scenarioStats := make(map[string]*models.ScenarioStat)
	for name, ss := range c.scenarioStats {
		scenarioStats[name] = &models.ScenarioStat{
			Name:          name,
			TotalRequests: ss.totalRequests,
			SuccessCount:  ss.successCount,
			FailedCount:   ss.failedCount,
			ErrorRate:     float64(ss.failedCount) / float64(ss.totalRequests) * 100,
			AvgRT:         averageDuration(ss.responseTimes),
			MinRT:         minDuration(ss.responseTimes),
			MaxRT:         maxDuration(ss.responseTimes),
		}
	}

	topErrors := make(map[string]int64)
	for k, v := range c.errors {
		if len(topErrors) < 10 {
			topErrors[k] = v
		}
	}

	statusCodesCopy := make(map[int]int64)
	for k, v := range c.statusCodes {
		statusCodesCopy[k] = v
	}

	return &models.Report{
		TestName:      testName,
		StartTime:     c.startTime,
		EndTime:       endTime,
		Duration:      duration,
		TotalRequests: total,
		SuccessCount:  success,
		FailedCount:   failed,
		ErrorRate:     errorRate,
		QPS:           qps,
		AvgRT:         averageDuration(sortedRT),
		MinRT:         minDuration(sortedRT),
		MaxRT:         maxDuration(sortedRT),
		P50:           percentile(sortedRT, 50),
		P90:           percentile(sortedRT, 90),
		P95:           percentile(sortedRT, 95),
		P99:           percentile(sortedRT, 99),
		StatusCodes:   statusCodesCopy,
		TopErrors:     topErrors,
		ScenarioStats: scenarioStats,
	}
}

func (c *Collector) PrintProgress() {
	stats := c.GetSnapshot()
	elapsed := time.Since(c.startTime)

	fmt.Printf("\r\033[K")
	fmt.Printf("Time: %-8s | QPS: %7.1f | Total: %-8d | Success: %-8d | Failed: %-8d | Workers: %d",
		formatDuration(elapsed),
		stats.CurrentQPS,
		stats.TotalRequests,
		stats.SuccessRequests,
		stats.FailedRequests,
		stats.CurrentWorkers,
	)
}

func PrintReport(report *models.Report) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("  TEST REPORT: %s\n", report.TestName)
	fmt.Println(strings.Repeat("=", 80))

	fmt.Printf("\n  Start Time:    %s\n", report.StartTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("  End Time:      %s\n", report.EndTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Duration:      %s\n", formatDuration(report.Duration))

	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("  REQUEST SUMMARY")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("  Total Requests:   %d\n", report.TotalRequests)
	fmt.Printf("  Successful:       %d\n", report.SuccessCount)
	fmt.Printf("  Failed:           %d\n", report.FailedCount)
	fmt.Printf("  Error Rate:       %.2f%%\n", report.ErrorRate)
	fmt.Printf("  QPS:              %.2f\n", report.QPS)

	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("  RESPONSE TIME DISTRIBUTION")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("  Average:  %-15s Min:  %-15s Max:  %-15s\n",
		report.AvgRT.Round(time.Millisecond),
		report.MinRT.Round(time.Millisecond),
		report.MaxRT.Round(time.Millisecond))
	fmt.Printf("  P50:      %-15s P90:  %-15s P95:  %-15s P99:  %-15s\n",
		report.P50.Round(time.Millisecond),
		report.P90.Round(time.Millisecond),
		report.P95.Round(time.Millisecond),
		report.P99.Round(time.Millisecond))

	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("  STATUS CODE DISTRIBUTION")
	fmt.Println(strings.Repeat("-", 80))
	for code, count := range report.StatusCodes {
		fmt.Printf("  %d: %d (%.1f%%)\n", code, count, float64(count)/float64(report.TotalRequests)*100)
	}

	if len(report.TopErrors) > 0 {
		fmt.Println("\n" + strings.Repeat("-", 80))
		fmt.Println("  TOP ERRORS")
		fmt.Println(strings.Repeat("-", 80))
		for err, count := range report.TopErrors {
			fmt.Printf("  %-70s %d\n", truncate(err, 70), count)
		}
	}

	if len(report.ScenarioStats) > 0 {
		fmt.Println("\n" + strings.Repeat("-", 80))
		fmt.Println("  SCENARIO BREAKDOWN")
		fmt.Println(strings.Repeat("-", 80))
		for name, ss := range report.ScenarioStats {
			fmt.Printf("\n  Scenario: %s\n", name)
			fmt.Printf("    Total: %d, Success: %d, Failed: %d, Error Rate: %.2f%%\n",
				ss.TotalRequests, ss.SuccessCount, ss.FailedCount, ss.ErrorRate)
			fmt.Printf("    Avg RT: %s, Min: %s, Max: %s\n",
				ss.AvgRT.Round(time.Millisecond),
				ss.MinRT.Round(time.Millisecond),
				ss.MaxRT.Round(time.Millisecond))
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

func averageDuration(rts []time.Duration) time.Duration {
	if len(rts) == 0 {
		return 0
	}
	var sum time.Duration
	for _, rt := range rts {
		sum += rt
	}
	return sum / time.Duration(len(rts))
}

func minDuration(rts []time.Duration) time.Duration {
	if len(rts) == 0 {
		return 0
	}
	min := rts[0]
	for _, rt := range rts[1:] {
		if rt < min {
			min = rt
		}
	}
	return min
}

func maxDuration(rts []time.Duration) time.Duration {
	if len(rts) == 0 {
		return 0
	}
	max := rts[0]
	for _, rt := range rts[1:] {
		if rt > max {
			max = rt
		}
	}
	return max
}

func percentile(sorted []time.Duration, p int) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	index := (p * len(sorted)) / 100
	if index >= len(sorted) {
		index = len(sorted) - 1
	}
	return sorted[index]
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
