package executor

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/apitester/apitester/internal/logger"
	"github.com/apitester/apitester/internal/scenario"
	"github.com/apitester/apitester/internal/stats"
	"github.com/apitester/apitester/pkg/models"
)

type Executor struct {
	config      *models.TestConfig
	collector   *stats.Collector
	weightedScenarios []*models.Scenario
	scenarioIdx uint64
}

func NewExecutor(config *models.TestConfig) *Executor {
	return &Executor{
		config:            config,
		collector:         stats.NewCollector(),
		weightedScenarios: scenario.GetScenariosByWeight(config),
	}
}

func (e *Executor) Run() (*models.Report, error) {
	logger.Info("Starting test: %s", e.config.Name)
	logger.Info("Base URL: %s", e.config.BaseURL)
	logger.Info("Concurrency mode: %s", e.config.Concurrency.Mode)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Info("\nReceived interrupt signal, stopping test...")
		cancel()
	}()

	go e.monitorProgress(ctx)

	switch e.config.Concurrency.Mode {
	case "workers":
		e.runWorkersMode(ctx)
	case "duration":
		e.runDurationMode(ctx)
	case "steps":
		e.runStepsMode(ctx)
	default:
		return nil, fmt.Errorf("unknown concurrency mode: %s", e.config.Concurrency.Mode)
	}

	logger.Info("Test completed, generating report...")
	report := e.collector.GenerateReport(e.config.Name)
	stats.PrintReport(report)

	return report, nil
}

func (e *Executor) runWorkersMode(ctx context.Context) {
	workers := e.config.Concurrency.Workers
	totalRequests := e.config.Concurrency.TotalRequests
	rampUp := e.config.Concurrency.RampUp

	logger.Info("Running with %d workers, %d total requests, ramp_up: %ds", workers, totalRequests, rampUp)

	var wg sync.WaitGroup
	var requestCount int64
	var activeWorkers int64

	for i := 0; i < workers; i++ {
		wg.Add(1)
		workerID := i

		go func() {
			defer wg.Done()

			if rampUp > 0 && workers > 1 {
				delay := time.Duration(float64(rampUp) / float64(workers-1) * float64(workerID))
				select {
				case <-ctx.Done():
					return
				case <-time.After(delay * time.Second):
				}
			}

			atomic.AddInt64(&activeWorkers, 1)
			e.collector.SetCurrentWorkers(int(atomic.LoadInt64(&activeWorkers)))

			for {
				select {
				case <-ctx.Done():
					return
				default:
					if totalRequests > 0 && atomic.LoadInt64(&requestCount) >= int64(totalRequests) {
						return
					}
					atomic.AddInt64(&requestCount, 1)
					e.executeNextScenario(ctx, workerID)
				}
			}
		}()
	}

	wg.Wait()
}

func (e *Executor) runDurationMode(ctx context.Context) {
	workers := e.config.Concurrency.Workers
	duration := time.Duration(e.config.Concurrency.Duration) * time.Second
	rampUp := e.config.Concurrency.RampUp

	logger.Info("Running with %d workers for %s, ramp_up: %ds", workers, duration, rampUp)

	testCtx, testCancel := context.WithTimeout(ctx, duration)
	defer testCancel()

	var wg sync.WaitGroup
	var activeWorkers int64

	for i := 0; i < workers; i++ {
		wg.Add(1)
		workerID := i

		go func() {
			defer wg.Done()

			if rampUp > 0 && workers > 1 {
				delay := time.Duration(float64(rampUp) / float64(workers-1) * float64(workerID))
				select {
				case <-testCtx.Done():
					return
				case <-time.After(delay * time.Second):
				}
			}

			atomic.AddInt64(&activeWorkers, 1)
			e.collector.SetCurrentWorkers(int(atomic.LoadInt64(&activeWorkers)))

			for {
				select {
				case <-testCtx.Done():
					return
				default:
					e.executeNextScenario(testCtx, workerID)
				}
			}
		}()
	}

	wg.Wait()
}

func (e *Executor) runStepsMode(ctx context.Context) {
	steps := e.config.Concurrency.Steps
	if len(steps) == 0 {
		logger.Error("No steps defined for steps mode")
		return
	}

	logger.Info("Running in steps mode with %d steps", len(steps))

	for stepIdx, step := range steps {
		select {
		case <-ctx.Done():
			return
		default:
		}

		logger.Info("Step %d/%d: %d workers for %ds",
			stepIdx+1, len(steps), step.Workers, step.Duration)

		stepCtx, stepCancel := context.WithTimeout(ctx, time.Duration(step.Duration)*time.Second)

		stepWg := &sync.WaitGroup{}
		e.collector.SetCurrentWorkers(step.Workers)

		for i := 0; i < step.Workers; i++ {
			stepWg.Add(1)
			go func(workerID int) {
				defer stepWg.Done()

				for {
					select {
					case <-stepCtx.Done():
						return
					default:
						e.executeNextScenario(stepCtx, workerID)
					}
				}
			}(i)
		}

		stepWg.Wait()
		stepCancel()

		snapshot := e.collector.GetSnapshot()
		logger.Info("Step %d completed: %d total requests, QPS: %.1f",
			stepIdx+1, snapshot.TotalRequests, snapshot.CurrentQPS)
	}
}

func (e *Executor) executeNextScenario(ctx context.Context, workerID int) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic in worker %d: %v", workerID, r)
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			logger.Error("Stack trace: %s", buf[:n])
		}
	}()

	scenarioList := e.weightedScenarios
	if len(scenarioList) == 0 {
		logger.Error("No scenarios available")
		return
	}

	idx := atomic.AddUint64(&e.scenarioIdx, 1) % uint64(len(scenarioList))
	s := scenarioList[idx]

	executor := scenario.NewScenarioExecutor(e.config, s)
	results := executor.Execute(ctx, workerID)

	for _, result := range results {
		e.collector.Record(result)
	}
}

func (e *Executor) monitorProgress(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.collector.UpdateQPS()
			e.collector.PrintProgress()
		}
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
