#!/bin/bash

# 酒店管理系统前端启动脚本

echo "========================================"
echo "  酒店管理系统 - 前端服务启动脚本"
echo "========================================"

# 检查Node环境
if ! command -v node &> /dev/null; then
    echo "❌ 错误: 未找到Node.js环境，请先安装Node.js 18+"
    exit 1
fi

echo "✅ Node.js环境检测通过"

# 进入前端目录
cd "$(dirname "$0")/frontend" || exit 1

# 检查node_modules
if [ ! -d "node_modules" ]; then
    echo "📦 安装依赖中..."
    npm install
fi

# 启动开发服务器
echo "🚀 启动前端开发服务器..."
echo "📡 服务地址: http://localhost:5173"
echo "📝 按 Ctrl+C 停止服务"
echo ""

npm run dev
