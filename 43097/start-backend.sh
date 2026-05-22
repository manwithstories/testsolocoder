#!/bin/bash

# 酒店管理系统启动脚本

echo "========================================"
echo "  酒店管理系统 - 后端服务启动脚本"
echo "========================================"

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到Go环境，请先安装Go 1.21+"
    exit 1
fi

echo "✅ Go环境检测通过"

# 进入后端目录
cd "$(dirname "$0")/backend" || exit 1

# 下载依赖
echo "📦 下载依赖中..."
go mod download

# 检查是否存在.env文件
if [ ! -f .env ]; then
    echo "⚠️  警告: .env文件不存在，将使用默认配置"
fi

# 启动服务
echo "🚀 启动后端服务..."
echo "📡 服务地址: http://localhost:8080"
echo "📝 按 Ctrl+C 停止服务"
echo ""

go run ./cmd/server/
