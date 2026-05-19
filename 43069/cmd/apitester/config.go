package main

import (
	"fmt"

	"github.com/apitester/apitester/internal/config"
	"github.com/spf13/cobra"
)

var (
	validateConfigFile string
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "验证配置文件",
	Long:  `验证YAML配置文件的语法和内容是否正确`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if validateConfigFile == "" {
			return fmt.Errorf("请指定配置文件路径")
		}

		cfg, err := config.LoadConfig(validateConfigFile)
		if err != nil {
			return fmt.Errorf("配置文件验证失败: %w", err)
		}

		fmt.Printf("配置文件验证通过!\n\n")
		fmt.Printf("测试名称: %s\n", cfg.Name)
		fmt.Printf("描述: %s\n", cfg.Description)
		fmt.Printf("Base URL: %s\n", cfg.BaseURL)
		fmt.Printf("场景数量: %d\n", len(cfg.Scenarios))
		fmt.Printf("并发模式: %s\n", cfg.Concurrency.Mode)
		fmt.Printf("并发数: %d\n", cfg.Concurrency.Workers)

		for i, s := range cfg.Scenarios {
			fmt.Printf("\n场景 %d: %s (权重: %d)\n", i+1, s.Name, s.Weight)
			fmt.Printf("  请求数: %d\n", len(s.Requests))
			for j, r := range s.Requests {
				fmt.Printf("    请求 %d: %s %s\n", j+1, r.Method, r.Path)
			}
		}

		return nil
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "生成示例配置文件",
	Long:  `生成一个示例YAML配置文件模板`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`示例配置文件内容:

name: "API性能测试"
description: "测试示例API的性能"
base_url: "https://api.example.com"
timeout: 30
retries: 2

variables:
  env: "production"
  token: "${API_TOKEN}"

headers:
  Content-Type: "application/json"
  User-Agent: "apitester/1.0"

concurrency:
  mode: "duration"
  workers: 50
  duration: 60
  ramp_up: 10

scenarios:
  - name: "用户登录流程"
    weight: 3
    requests:
      - name: "登录"
        method: "POST"
        path: "/api/auth/login"
        body: '{"username": "user", "password": "pass"}'
        assertions:
          - type: "status_code"
            operator: "eq"
            value: 200
          - type: "response_time"
            operator: "lte"
            value: 500
        extract:
          - name: "auth_token"
            from: "body"
            path: "data.token"

      - name: "获取用户信息"
        method: "GET"
        path: "/api/user/profile"
        headers:
          Authorization: "Bearer {{auth_token}}"
        assertions:
          - type: "status_code"
            operator: "eq"
            value: 200
          - type: "body_field"
            field: "code"
            operator: "eq"
            value: 0

  - name: "数据查询"
    weight: 1
    requests:
      - name: "查询列表"
        method: "GET"
        path: "/api/items?page=1&size=20"
        assertions:
          - type: "status_code"
            operator: "eq"
            value: 200
          - type: "response_time"
            operator: "lte"
            value: 1000
`)
	},
}

func init() {
	validateCmd.Flags().StringVarP(&validateConfigFile, "config", "c", "", "配置文件路径")
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(initCmd)
}
