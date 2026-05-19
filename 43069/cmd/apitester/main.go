package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "apitester",
	Short: "API性能压测与诊断工具",
	Long: `apitester 是一个功能强大的API性能压测与诊断工具。
支持并发测试、链式接口调用、断言校验、结果分析等功能。`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
