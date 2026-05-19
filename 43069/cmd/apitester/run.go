package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/apitester/apitester/internal/config"
	"github.com/apitester/apitester/internal/executor"
	"github.com/apitester/apitester/internal/logger"
	"github.com/apitester/apitester/pkg/models"
	"github.com/spf13/cobra"
)

var (
	runConfigFile string
	runLogLevel   string
	runLogFile    string
	runWorkers    int
	runDuration   int
	runRequests   int
	runReportFile string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "运行API压测",
	Long:  `根据配置文件运行API性能压测`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if runConfigFile == "" {
			return fmt.Errorf("请指定配置文件路径")
		}

		logger.Init(runLogLevel, runLogFile, true)
		defer logger.Close()

		cfg, err := config.LoadConfig(runConfigFile)
		if err != nil {
			return fmt.Errorf("加载配置文件失败: %w", err)
		}

		if runWorkers > 0 {
			cfg.Concurrency.Workers = runWorkers
		}
		if runDuration > 0 {
			cfg.Concurrency.Duration = runDuration
		}
		if runRequests > 0 {
			cfg.Concurrency.TotalRequests = runRequests
		}

		exec := executor.NewExecutor(cfg)
		report, err := exec.Run()
		if err != nil {
			return fmt.Errorf("执行测试失败: %w", err)
		}

		if runReportFile != "" {
			if err := saveReport(report, runReportFile); err != nil {
				logger.Warn("保存报告失败: %v", err)
			} else {
				logger.Info("报告已保存到: %s", runReportFile)
			}
		}

		return nil
	},
}

func saveReport(report *models.Report, path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.Create(absPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(generateTextReport(report))
	return err
}

func generateTextReport(report *models.Report) string {
	return fmt.Sprintf(`API性能测试报告
========================================
测试名称: %s
开始时间: %s
结束时间: %s
持续时间: %s

请求统计:
  总请求数: %d
  成功请求: %d
  失败请求: %d
  错误率: %.2f%%
  QPS: %.2f

响应时间(ms):
  平均: %d
  最小: %d
  最大: %d
  P50: %d
  P90: %d
  P95: %d
  P99: %d
`,
		report.TestName,
		report.StartTime.Format("2006-01-02 15:04:05"),
		report.EndTime.Format("2006-01-02 15:04:05"),
		report.Duration.Round(time.Second),
		report.TotalRequests,
		report.SuccessCount,
		report.FailedCount,
		report.ErrorRate,
		report.QPS,
		report.AvgRT.Milliseconds(),
		report.MinRT.Milliseconds(),
		report.MaxRT.Milliseconds(),
		report.P50.Milliseconds(),
		report.P90.Milliseconds(),
		report.P95.Milliseconds(),
		report.P99.Milliseconds(),
	)
}

func init() {
	runCmd.Flags().StringVarP(&runConfigFile, "config", "c", "", "配置文件路径 (YAML)")
	runCmd.Flags().StringVarP(&runLogLevel, "log-level", "l", "info", "日志级别 (debug, info, warn, error)")
	runCmd.Flags().StringVar(&runLogFile, "log-file", "", "日志文件路径")
	runCmd.Flags().IntVarP(&runWorkers, "workers", "w", 0, "并发数 (覆盖配置文件)")
	runCmd.Flags().IntVarP(&runDuration, "duration", "d", 0, "测试持续时间(秒) (覆盖配置文件)")
	runCmd.Flags().IntVarP(&runRequests, "requests", "n", 0, "总请求数 (覆盖配置文件)")
	runCmd.Flags().StringVarP(&runReportFile, "report", "r", "", "输出报告文件路径")

	rootCmd.AddCommand(runCmd)
}
