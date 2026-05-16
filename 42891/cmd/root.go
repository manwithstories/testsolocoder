package cmd

import (
	"fmt"
	"os"
	"timetrack/cmd/backupcmd"
	"timetrack/cmd/exportcmd"
	"timetrack/cmd/project"
	"timetrack/cmd/report"
	"timetrack/cmd/timecmd"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "timetrack",
	Short: "时间追踪工具 - 专为自由职业者设计",
	Long: `timetrack 是一个命令行时间追踪工具，帮助自由职业者管理项目、记录工作时间、
生成统计报表并导出数据。支持项目管理、标签分类、跨天检测、备份恢复等功能。`,
	Version: "1.0.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(project.ProjectCmd)
	rootCmd.AddCommand(timecmd.TimeCmd)
	rootCmd.AddCommand(report.ReportCmd)
	rootCmd.AddCommand(exportcmd.ExportCmd)
	rootCmd.AddCommand(backupcmd.BackupCmd)
	rootCmd.AddCommand(backupcmd.RestoreCmd)
}
